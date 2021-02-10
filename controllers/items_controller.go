package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cmd-ctrl-q/bookstore_items-api/utils/http_utils"
	"github.com/cmd-ctrl-q/bookstore_utils-go/rest_errors"

	"github.com/cmd-ctrl-q/bookstore_items-api/domain/items"
	"github.com/cmd-ctrl-q/bookstore_items-api/services"

	"github.com/cmd-ctrl-q/bookstore_oauth-go/oauth"
)

var (
	ItemsController itemsControllerInterface = &itemsController{}
)

type itemsControllerInterface interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}

type itemsController struct{}

func (c *itemsController) Create(w http.ResponseWriter, r *http.Request) {
	// try to authenticate the request and check if access token is valid
	if err := oauth.AuthenticateRequest(r); err != nil {
		// cannot validate access_token yet bc api is not up and running.
		// if there is an error when authenticating, do nothing.
		// respErr := rest_errors.RestErr{
		// 	Message: err.Error,
		// 	Status:  err.Status,
		// 	Error:   "authentication_request_error",
		// 	Causes:  nil,
		// }
		// http_utils.RespondError(w, respErr)
		return
	}

	sellerID := oauth.GetCallerID(r)
	// no user in X-Caller-ID, so return unauthroized
	if sellerID == 0 {
		// ie. no access_token so not able to use GetCallerId to validate the user id
		respErr := rest_errors.NewUnauthorizedError("unauthorized user")
		http_utils.RespondError(w, respErr)
		return
	}

	// if access_token is valid OR if there is no access token, get the body of the request.
	// try to get the json body from the request
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.RespondError(w, respErr)
		return
	}
	defer r.Body.Close()

	// try to unmarshal the request body and populate into the items request object
	// this will attempt to populate ALL of the fields in the Items struct
	var itemRequest items.Item
	if err = json.Unmarshal(requestBody, &itemRequest); err != nil {
		respErr := rest_errors.NewBadRequestError("invalid item json body")
		http_utils.RespondError(w, respErr)
		return
	}

	// get the client id from the oauth request (after doing the authenticated request above)
	// ie get user id from GetCallerID and put into the item's Seller field
	// IE. if the callerID in the header is associated to a sellerID,
	// then they can create an item.
	itemRequest.Seller = sellerID

	// finally, call the services to create the item
	result, createErr := services.ItemsService.Create(itemRequest)
	if createErr != nil {
		http_utils.RespondError(w, createErr)
		return
	}

	http_utils.RespondJson(w, http.StatusCreated, result)
}

// Get an item
func (c *itemsController) Get(w http.ResponseWriter, r *http.Request) {

}
