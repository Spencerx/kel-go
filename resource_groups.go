package kel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/derekdowling/go-json-spec-handler"
)

const (
	resourceType string = "resource-groups"
)

// ResourceGroupService represents the link between the ResourceGroup and the
// Client.
type ResourceGroupService struct {
	client *Client
}

// Create sends an HTTP request to create the Kel resource group.
func (srv *ResourceGroupService) Create(resourceGroup *ResourceGroup) *ResourceGroupCreateRequest {
	return &ResourceGroupCreateRequest{
		srv:           srv,
		resourceGroup: resourceGroup,
	}
}

// CreateWithToken sends an HTTP request to create the Kel resource group
// providing the given token.
func (srv *ResourceGroupService) CreateWithToken(resourceGroup *ResourceGroup, token string) *ResourceGroupCreateRequest {
	return &ResourceGroupCreateRequest{
		srv:           srv,
		resourceGroup: resourceGroup,
		token:         token,
	}
}

// ResourceGroupCreateRequest represents a request to create a reource group.
type ResourceGroupCreateRequest struct {
	srv           *ResourceGroupService
	resourceGroup *ResourceGroup
	token         string
}

// Do executes the create resource group request.
func (req *ResourceGroupCreateRequest) Do() error {
	var u url.URL
	u = *req.srv.client.baseURL
	u.Path += "/resource-groups"
	outreq, err := http.NewRequest("POST", (&u).String(), nil)
	if err != nil {
		return err
	}
	obj, objErr := jsh.NewObject(req.resourceGroup.Name, resourceType, req.resourceGroup)
	if objErr != nil {
		return objErr
	}
	if err := obj.Validate(outreq, false); err != nil {
		return fmt.Errorf("error preparing object: %s", err.Error())
	}
	doc := jsh.Build(obj)
	docSerialized, err := json.MarshalIndent(doc, "", " ")
	if err != nil {
		return fmt.Errorf("error serializing document: %s", err.Error())
	}
	if req.token != "" {
		outreq.Header.Set("X-Kel-Token", req.token)
	}
	outreq.Header.Set("Content-Type", jsh.ContentType)
	outreq.Body = jsh.CreateReadCloser(docSerialized)
	outreq.ContentLength = int64(len(docSerialized))
	res, err := req.srv.client.Do(outreq)
	if err != nil {
		return err
	}
	parser := &jsh.Parser{Method: "", Headers: res.Header}
	doc, err = parser.Document(res.Body, jsh.ObjectMode)
	doc.Status = res.StatusCode
	switch res.StatusCode {
	case http.StatusBadRequest:
		return errors.New(doc.Errors[0].Detail)
	default:
		return fmt.Errorf("unknown response from API: %s", res.Status)
	}
}

// List returns all resource groups reachable by the API
func (srv *ResourceGroupService) List() *ResourceGroupListRequest {
	return &ResourceGroupListRequest{srv: srv}
}

// ResourceGroupListRequest represents a request for a list of resource groups.
type ResourceGroupListRequest struct {
	srv *ResourceGroupService
}

// Include ...
func (req *ResourceGroupListRequest) Include() *ResourceGroupListRequest {
	return req
}

// Do ...
func (req *ResourceGroupListRequest) Do() ([]*ResourceGroup, error) {
	var u url.URL
	u = *req.srv.client.baseURL
	u.Path += "/resource-groups"
	outreq, err := http.NewRequest("GET", (&u).String(), nil)
	if err != nil {
		return nil, err
	}
	outreq.Header.Set("Content-Type", jsh.ContentType)
	res, err := req.srv.client.Do(outreq)
	if err != nil {
		return nil, err
	}
	fmt.Println(res.Status)
	return nil, nil
}

// Get returns the resource group with the given name reachable by the API
func (srv *ResourceGroupService) Get(name string) (*ResourceGroup, error) {
	return nil, nil
}

// ResourceGroup represents a Kel resource group.
type ResourceGroup struct {
	Name string `json:"name"`
}

// GetResourceType ...
func (resourceGroup *ResourceGroup) GetResourceType() string {
	return resourceType
}

// GetID ...
func (resourceGroup *ResourceGroup) GetID() string {
	return resourceGroup.Name
}

// Reload will get an updated resource group and point to it as its own.
func (resourceGroup *ResourceGroup) Reload() error {
	return nil
}

// Save will persistent local data with the API.
func (resourceGroup *ResourceGroup) Save() error {
	return nil
}

// Delete will destroy the resource group.
func (resourceGroup *ResourceGroup) Delete() error {
	return nil
}

// ResourceGroupUser represents a Kel resource group user.
type ResourceGroupUser struct {
}
