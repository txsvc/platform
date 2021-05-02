package platform

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"

	"github.com/txsvc/platform/pkg/env"
)

type (
	// Client holds all clients needed to access basic Google Cloud services
	Client struct {
		DatastoreClient *datastore.Client
		StorageClient   *storage.Client
	}
)

var client *Client

func init() {
	if client != nil {
		return // singleton
	}

	// FIXME,DEPRECATED remove this in the future
	projectID := env.GetString("PROJECT_ID", "")
	if projectID == "" {
		return // FIXME should this be a hard exception?
	}

	cl, err := NewClient(context.Background(), projectID, env.GetString("SERVICE_NAME", "default"))

	if err != nil {
		log.Fatal(err)
	}

	client = cl
}

// NewClient creates a new client
func NewClient(ctx context.Context, projectID, serviceName string) (*Client, error) {
	c := Client{}

	// initialize Cloud Datastore
	ds, err := datastore.NewClient(ctx, projectID)
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
}

//
// FIXME,DEPRECATED remove these helpers in the future
//

// Close the platform related clients
func Close() {
	if client != nil {
		client.Close()
	}
}

// DataStore returns a reference to the datastore client
func DataStore() *datastore.Client {
	if client != nil {
		return client.DatastoreClient
	}
	log.Fatal(fmt.Errorf("Google Datastore is not initialized"))
	return nil
}

// Storage returns a reference to the storage client
func Storage() *storage.Client {
	if client != nil {
		return client.StorageClient
	}
	log.Fatal(fmt.Errorf("Google Cloud Storage is not initialized"))
	return nil
}
