package account

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"

	ds "github.com/txsvc/platform/v2/pkg/datastore"
	"github.com/txsvc/platform/v2/pkg/id"
	"github.com/txsvc/platform/v2/pkg/loader"
	"github.com/txsvc/platform/v2/pkg/timestamp"
)

const (
	// AccountActive indicates a confirmed account with a valid login & authorization
	AccountActive = 1
	// AccountLoggedOut indicates a confirmed account that is not logged in
	AccountLoggedOut = 0

	// AccountDeactivated indicates an account that has been deactivated due to e.g. account deletion or UserID swap
	AccountDeactivated = -1
	// AccountBlocked signals an issue with the account that needs intervention
	AccountBlocked = -2
	// AccountUnconfirmed well guess what?
	AccountUnconfirmed = -3

	// DatastoreAccounts collection ACCOUNTS
	datastoreAccounts string = "ACCOUNTS"
)

type (
	// Account represents an account for a user or client (e.g. API, bot)
	Account struct {
		Realm    string `json:"realm"`     // KEY
		UserID   string `json:"user_id"`   // KEY external id for the entity e.g. email for a user
		ClientID string `json:"client_id"` // a unique id within [realm,user_id]
		// status and other metadata
		Status int `json:"status"` // default == AccountUnconfirmed
		// login auditing
		LastLogin  int64  `json:"-"`
		LoginCount int    `json:"-"`
		LoginFrom  string `json:"-"`
		// internal
		Token     string `json:"-"` // a temporary token to confirm the account or to exchanged for the "real" token
		Expires   int64  `json:"-"` // 0 == never
		Confirmed int64  `json:"-"`
		Created   int64  `json:"-"`
		Updated   int64  `json:"-"`
	}
)

var (
	// ErrAccountExists indicates that the account already exists and can't be created
	ErrAccountExists = errors.New("account exists")

	// loader used to cache accounts
	accountLoader = loader.New(AccountLoaderFunc, loader.DefaultTTL)
)

func (acc *Account) Equal(a *Account) bool {
	if a == nil {
		return false
	}
	return acc.Realm == a.Realm && acc.ClientID == a.ClientID && acc.UserID == a.UserID
}

// CreateAccount creates an new account in the given realm. userID has to be unique and a new clientID will be assigned.
func CreateAccount(ctx context.Context, realm, userID string, expires int) (*Account, error) {
	acc, err := FindAccountByUserID(ctx, realm, userID)
	if err != nil {
		return nil, err
	}
	if acc != nil {
		return nil, ErrAccountExists
	}

	now := timestamp.Now()
	token, _ := id.ShortUUID() // temporary token to confirm the new account. this is not an authorization token or such!

	uid := ""
	for {
		// make sure that the new clientID is unique
		uid, _ = id.ShortUUID()
		acc, err := LookupAccount(ctx, realm, uid)
		if err != nil {
			return nil, err
		}
		if acc == nil {
			break
		}
		// try again ...
	}

	account := Account{
		Realm:     realm,
		UserID:    userID,
		ClientID:  uid,
		Status:    AccountUnconfirmed,
		Token:     token,
		Expires:   timestamp.IncT(timestamp.Now(), expires),
		Confirmed: 0,
		Created:   now,
		Updated:   now,
	}

	if err := UpdateAccount(ctx, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// LookupAccount retrieves an account within a given realm
func LookupAccount(ctx context.Context, realm, clientID string) (*Account, error) {

	k := nativeKey(namedKey(realm, clientID))
	account, err := accountLoader.Load(ctx, k.Encode())
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, nil // not found
	}
	return account.(*Account), nil

}

func UpdateAccount(ctx context.Context, account *Account) error {
	k := nativeKey(account.Key())

	// there is a SMALL time window where the cache and the datastore are inconsistent ...

	account.Updated = timestamp.Now()
	if _, err := ds.DataStore().Put(ctx, k, account); err != nil {
		return err
	}

	accountLoader.Remove(ctx, k.Encode())
	return nil
}

// FindAccountUserID retrieves an account bases on the user id
func FindAccountByUserID(ctx context.Context, realm, userID string) (*Account, error) {
	var accounts []*Account
	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(datastoreAccounts).Filter("Realm =", realm).Filter("UserID =", userID), &accounts); err != nil {
		return nil, err
	}
	if accounts == nil {
		return nil, nil
	}
	return accounts[0], nil
}

// FindAccountByToken retrieves an account bases on either the temporary token or the auth token
func FindAccountByToken(ctx context.Context, token string) (*Account, error) {
	var accounts []*Account
	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(datastoreAccounts).Filter("Token =", token), &accounts); err != nil {
		return nil, err
	}
	if accounts == nil {
		return nil, nil
	}
	return accounts[0], nil
}

// ResetAccountChallenge creates a new confirmation token and resets the timer
func ResetAccountChallenge(ctx context.Context, acc *Account, expires int) (*Account, error) {
	token, _ := id.ShortUUID()
	acc.Expires = timestamp.IncT(timestamp.Now(), expires)
	acc.Token = "ac." + token
	acc.Status = AccountUnconfirmed

	if err := UpdateAccount(ctx, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

// ResetTemporaryToken creates a new temporary token and resets the timer
func ResetTemporaryToken(ctx context.Context, acc *Account, expires int) (*Account, error) {
	token, _ := id.ShortUUID()
	acc.Expires = timestamp.IncT(timestamp.Now(), expires)
	acc.Token = "tt." + token
	acc.Status = AccountLoggedOut

	if err := UpdateAccount(ctx, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

//
// keys and dataloader
//

func (acc *Account) Key() string {
	return namedKey(acc.Realm, acc.ClientID)
}

func nativeKey(key string) *datastore.Key {
	return datastore.NameKey(datastoreAccounts, key, nil)
}

func namedKey(part1, part2 string) string {
	return part1 + "." + part2
}

// AccountLoaderFunc implements the LoaderFunc interface for retrieving account resources
func AccountLoaderFunc(ctx context.Context, key string) (interface{}, error) {
	var account Account

	k, err := datastore.DecodeKey(key)
	if err != nil {
		return nil, err
	}

	err = ds.DataStore().Get(ctx, k, &account)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}
