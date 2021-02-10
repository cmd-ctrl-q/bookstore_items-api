package elasticsearch

import (
	"context"
	"fmt"
	"time"

	"github.com/cmd-ctrl-q/bookstore_utils-go/logger"
	"github.com/olivere/elastic"
)

var (
	// during testing, mock this interface to test the callers in it
	Client esClientInterface = &esClient{} // create new client
)

type esClientInterface interface {
	setClient(c *elastic.Client)
	Index(string, interface{}) (*elastic.IndexResponse, error) // for indexing a new document in ES
}

type esClient struct {
	// client points to an es client but we never set it here
	client *elastic.Client // this is the actual library were implementing
}

// Must modify you database structure in production when the application starts.
// make sure you db is already configured, created, and available before using it.
func Init() {
	log := logger.GetLogger()
	// create a new client
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log),
		elastic.SetInfoLog(log),
	)
	if err != nil {
		// panic b/c app cannot run without elastic search
		panic(err)
	}
	// set global var Client to the created client
	Client.setClient(client)
	// now other packages can use the global `Client` variable as it has an ES client and its global

	// Create index (ie es database) if it does not exist.
}

func (c *esClient) setClient(client *elastic.Client) {
	c.client = client
}

// parameter says: generate any document on any index that you want.
func (c *esClient) Index(index string, doc interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	// return c.client.Index().Do(ctx)
	result, err := c.client.Index().
		Index(index).
		BodyJson(doc).
		Do(ctx)

	if err != nil {
		// put logger where error occurs
		logger.Error(
			fmt.Sprintf("error when trying to index document in es: %s", index), err)
		return nil, err
	}
	return result, nil
}
