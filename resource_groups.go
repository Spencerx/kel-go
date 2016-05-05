package kel

import (
	"fmt"
	"net/http"

	"github.com/brosner/go-json-spec-handler"
)

const (
	resourceGroupResourceType string = "resource-groups"
	resourceGroupAPIPath      string = "/resource-groups"
)

// ResourceGroupService represents the link between the ResourceGroup and the
// Client.
type ResourceGroupService struct {
	client *Client
}

// Create sends an HTTP request to create the Kel resource group.
func (srv *ResourceGroupService) Create(resourceGroup *ResourceGroup) CreateRequest {
	return &createRequest{
		client:   srv.client,
		path:     resourceGroupAPIPath,
		buildDoc: buildResourceGroupDoc(resourceGroup),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(resourceGroupResourceType, resourceGroup); objErr != nil {
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
		path:     resourceGroupAPIPath,
		hdr:      make(http.Header),
		buildDoc: buildResourceGroupDoc(resourceGroup),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(resourceGroupResourceType, resourceGroup); objErr != nil {
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
		path:   resourceGroupAPIPath,
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
	resourceGroup.Name = name // ID for lookup
	return &getRequest{
		client: srv.client,
		path:   resourceGroup.getDetailPath(),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(resourceGroupResourceType, resourceGroup); objErr != nil {
				return objErr
			}
			resourceGroup.srv = srv
			return nil
		},
	}
}

func (resourceGroup *ResourceGroup) getDetailPath() string {
	return fmt.Sprintf(
		"%s/%s",
		resourceGroupAPIPath,
		resourceGroup.GetID(),
	)
}

// GetResourceType ...
func (resourceGroup *ResourceGroup) GetResourceType() string {
	return resourceGroupResourceType
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
		path:     resourceGroup.getDetailPath(),
		buildDoc: buildResourceGroupDoc(resourceGroup),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(resourceGroupResourceType, resourceGroup); objErr != nil {
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
		path:   resourceGroup.getDetailPath(),
	}
	return req.Do()
}

func buildResourceGroupDoc(resourceGroup *ResourceGroup) func(outreq *http.Request) (*jsh.Document, error) {
	return func(outreq *http.Request) (*jsh.Document, error) {
		obj, objErr := jsh.NewObject(resourceGroup.GetID(), resourceGroupResourceType, resourceGroup)
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
