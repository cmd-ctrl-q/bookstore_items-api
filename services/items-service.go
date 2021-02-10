package services

import (
	"net/http"

	"github.com/cmd-ctrl-q/bookstore_items-api/domain/items"

	"github.com/cmd-ctrl-q/bookstore_utils-go/rest_errors"
)

var (
	ItemsService itemsServiceInterface = &itemsService{}
)

type itemsServiceInterface interface {
	Create(items.Item) (*items.Item, rest_errors.RestErr)
	Get(string) (*items.Item, rest_errors.RestErr)
}

type itemsService struct{}

func (s *itemsService) Create(item items.Item) (*items.Item, rest_errors.RestErr) {

	// save item in mysql, cassandra, elastic search, memcache, etc.
	// ie any persistent layer
	if err := item.Save(); err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *itemsService) Get(string) (*items.Item, rest_errors.RestErr) {
	return nil, rest_errors.NewRestError("implement me!", http.StatusNotFound, "not_implemented", nil)
}
