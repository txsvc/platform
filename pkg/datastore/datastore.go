package datastore

import (
	"context"
	"errors"
	"log"

	ds "cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"

	"github.com/txsvc/platform/v2"
	"github.com/txsvc/platform/v2/pkg/env"
)

type (
	// Client holds all clients needed to access basic Google Cloud services
	Client struct {
		DatastoreClient *ds.Client
		StorageClient   *storage.Client
	}
)

var (
	client *Client

	ErrInvalidConf = errors.New("missing PROJECT_ID")
	ErrNoDatastore = errors.New("google datastore is not initialized")
	ErrNoStorage   = errors.New("google cloud storage is not initialized")
)

func init() {
	cl, err := NewClient(context.Background(), env.GetString("PROJECT_ID", ""), env.GetString("SERVICE_NAME", "default"))
	if err != nil {
		platform.ReportError(err)
	}
	client = cl
}

// NewClient creates a new client
func NewClient(ctx context.Context, projectID, serviceName string) (*Client, error) {
	c := Client{}

	if projectID == "" {
		return nil, ErrInvalidConf
	}
	if serviceName == "" {
		serviceName = "default"
	}

	// initialize Cloud Datastore
	ds, err := ds.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	c.DatastoreClient = ds

	// initialize Cloud Storage
	cs, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	c.StorageClient = cs

	return &c, nil
}

// Close closes all clients to the Google Cloud services
func (c *Client) Close() {
	if c.StorageClient != nil {
		c.StorageClient.Close()
	}
	if c.DatastoreClient != nil {
		c.DatastoreClient.Close()
	}
	c.StorageClient = nil
	c.DatastoreClient = nil
}

//
// DEPRECATED remove these helpers in the future
//

// Close the platform related clients
func Close() {
	if client != nil {
		client.Close()
	}
}

// DataStore returns a reference to the datastore client
func DataStore() *ds.Client {
	if client.DatastoreClient != nil {
		return client.DatastoreClient
	}
	log.Fatal(ErrNoDatastore)
	return nil
}

// Storage returns a reference to the storage client
func Storage() *storage.Client {
	if client.StorageClient != nil {
		return client.StorageClient
	}
	log.Fatal(ErrNoStorage)
	return nil
}
