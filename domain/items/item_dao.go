package items

import (
	"errors"

	"github.com/cmd-ctrl-q/bookstore_items-api/clients/elasticsearch"
	"github.com/cmd-ctrl-q/bookstore_utils-go/rest_errors"
)

const (
	indexItems = "items"
)

func (i *Item) Save() rest_errors.RestErr {
	result, err := elasticsearch.Client.Index(indexItems, i)
	// error when trying to index document in es
	if err != nil {
		// first error message goes to user. they know nothing about elastic search.
		return rest_errors.NewInternalServerError("error when trying to save item", errors.New("database_error"))
	}
	i.ID = result.Id
	return nil
}
