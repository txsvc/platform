package platform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txsvc/platform/v2/pkg/env"
)

func TestPackageInit(t *testing.T) {
	assert.NotNil(t, client)
	assert.NotNil(t, DataStore())
	assert.NotNil(t, Storage())
}

func TestInitClient(t *testing.T) {

	ds, err := NewClient(context.Background(), env.GetString("PROJECT_ID", ""), env.GetString("SERVICE_NAME", "default"))

	assert.NoError(t, err)
	assert.NotNil(t, ds)
	assert.NotNil(t, ds.DatastoreClient)
	assert.NotNil(t, ds.StorageClient)
}

func TestInitClientFail(t *testing.T) {

	ds, err := NewClient(context.Background(), "", "")

	assert.Nil(t, ds)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidConf, err)
}

func TestClientClose(t *testing.T) {
	ds, err := NewClient(context.Background(), env.GetString("PROJECT_ID", ""), env.GetString("SERVICE_NAME", "default"))

	assert.NoError(t, err)
	assert.NotNil(t, ds)

	ds.Close()
	assert.Nil(t, ds.DatastoreClient)
	assert.Nil(t, ds.StorageClient)

	ds.Close() // should be idempotent
	assert.Nil(t, ds.DatastoreClient)
	assert.Nil(t, ds.StorageClient)
}
