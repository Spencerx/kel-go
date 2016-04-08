package kel

import (
	"fmt"
	"net/http"

	"github.com/brosner/go-json-spec-handler"
)

const (
	resourceType string = "resource-groups"
	apiPath      string = "/resource-groups"
)

// ResourceGroupService represents the link between the ResourceGroup and the
// Client.
type ResourceGroupService struct {
	client *Client
}

func (srv *ResourceGroupService) getDetailPath(id string) string {
	return fmt.Sprintf("%s/%s", apiPath, id)
}

// Create sends an HTTP request to create the Kel resource group.
func (srv *ResourceGroupService) Create(resourceGroup *ResourceGroup) CreateRequest {
	return &createRequest{
		client:   srv.client,
		path:     apiPath,
		buildDoc: buildDoc(resourceGroup),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(resourceType, resourceGroup); objErr != nil {
				return objErr
			}
			resourceGroup.srv = srv
			return nil
		},
	}
}

// CreateWithToken sends an HTTP request to create the Kel resource group
// providing the given token.
func (srv *ResourceGroupService) CreateWithToken(resourceGroup *ResourceGroup, token string) CreateRequest {
	req := &createRequest{
		client:   srv.client,
		path:     apiPath,
		hdr:      make(http.Header),
		buildDoc: buildDoc(resourceGroup),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(resourceType, resourceGroup); objErr != nil {
				return objErr
			}
			resourceGroup.srv = srv
			return nil
		},
	}
	req.hdr.Set("X-Kel-Token", token)
	return req
}

// List returns all resource groups reachable by the API
func (srv *ResourceGroupService) List(resourceGroups *[]*ResourceGroup) ListRequest {
	return &listRequest{
		client: srv.client,
		path:   apiPath,
		handler: func(document *jsh.Document) error {
			for i := range document.Data {
				obj := document.Data[i]
				resourceGroup := &ResourceGroup{srv: srv}
				if objErr := obj.Unmarshal(resourceGroup.GetResourceType(), resourceGroup); objErr != nil {
					return objErr
				}
				*resourceGroups = append(*resourceGroups, resourceGroup)
			}
			return nil
		},
	}
}

// Get returns the resource group with the given name reachable by the API
func (srv *ResourceGroupService) Get(name string, resourceGroup *ResourceGroup) GetRequest {
	return &getRequest{
		client: srv.client,
		path:   srv.getDetailPath(name),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(resourceType, resourceGroup); objErr != nil {
				return objErr
			}
			resourceGroup.srv = srv
			return nil
		},
	}
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
	reloaded := &ResourceGroup{}
	if err := resourceGroup.srv.Get(resourceGroup.GetID(), reloaded).Do(); err != nil {
		return err
	}
	*resourceGroup = *reloaded
	return nil
}

// Save will persistent local data with the API.
func (resourceGroup *ResourceGroup) Save() error {
	req := &updateRequest{
		client:   resourceGroup.srv.client,
		path:     resourceGroup.srv.getDetailPath(resourceGroup.GetID()),
		buildDoc: buildDoc(resourceGroup),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(resourceType, resourceGroup); objErr != nil {
				return objErr
			}
			return nil
		},
	}
	return req.Do()
}

// Delete will destroy the resource group.
func (resourceGroup *ResourceGroup) Delete() error {
	req := &deleteRequest{
		client: resourceGroup.srv.client,
		path:   resourceGroup.srv.getDetailPath(resourceGroup.GetID()),
	}
	return req.Do()
}

func buildDoc(resourceGroup *ResourceGroup) func(outreq *http.Request) (*jsh.Document, error) {
	return func(outreq *http.Request) (*jsh.Document, error) {
		obj, objErr := jsh.NewObject(resourceGroup.GetID(), resourceType, resourceGroup)
		if objErr != nil {
			return nil, objErr
		}
		if err := obj.Validate(outreq, false); err != nil {
			return nil, fmt.Errorf("error preparing object: %s", err.Error())
		}
		doc := jsh.Build(obj)
		return doc, nil
	}
}
