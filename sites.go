package kel

import (
	"fmt"
	"net/http"

	"github.com/brosner/go-json-spec-handler"
)

const (
	sitesResourceType string = "sites"
	sitesAPIPath      string = "/sites"
)

// SiteService represents the link between the Site and the Client.
type SiteService struct {
	client *Client
}

func (srv *SiteService) getPath(resourceGroup *ResourceGroup) string {
	return fmt.Sprintf("%s%s", resourceGroup.getDetailPath(), sitesAPIPath)
}

// Create sends an HTTP request to create the Kel site.
func (srv *SiteService) Create(site *Site) CreateRequest {
	return &createRequest{
		client:   srv.client,
		path:     srv.getPath(site.ResourceGroup),
		buildDoc: buildSiteDoc(site),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(sitesResourceType, site); objErr != nil {
				return objErr
			}
			site.srv = srv
			return nil
		},
	}
}

// List returns all sites reachable by the API
func (srv *SiteService) List(sites *[]*Site) ListRequest {
	return &listRequest{
		client: srv.client,
		path:   sitesAPIPath,
		handler: func(document *jsh.Document) error {
			for i := range document.Data {
				obj := document.Data[i]
				site := &Site{srv: srv}
				if objErr := obj.Unmarshal(site.GetResourceType(), site); objErr != nil {
					return objErr
				}
				*sites = append(*sites, site)
			}
			return nil
		},
	}
}

// Get returns the site with the given name reachable by the API
func (srv *SiteService) Get(name string, site *Site) GetRequest {
	return &getRequest{
		client: srv.client,
		path:   site.getDetailPath(),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(sitesResourceType, site); objErr != nil {
				return objErr
			}
			site.srv = srv
			return nil
		},
	}
}

func (site *Site) getDetailPath() string {
	return fmt.Sprintf(
		"%s%s/%s",
		site.ResourceGroup.getDetailPath(),
		sitesAPIPath,
		site.GetID(),
	)
}

// GetResourceType ...
func (site *Site) GetResourceType() string {
	return sitesResourceType
}

// GetID ...
func (site *Site) GetID() string {
	return site.Name
}

// Reload will get an updated site and point to it as its own.
func (site *Site) Reload() error {
	reloaded := &Site{}
	if err := site.srv.Get(site.GetID(), reloaded).Do(); err != nil {
		return err
	}
	*site = *reloaded
	return nil
}

// Save will persistent local data with the API.
func (site *Site) Save() error {
	req := &updateRequest{
		client:   site.srv.client,
		path:     site.getDetailPath(),
		buildDoc: buildSiteDoc(site),
		handler: func(document *jsh.Document) error {
			obj := document.Data[0]
			if objErr := obj.Unmarshal(sitesResourceType, site); objErr != nil {
				return objErr
			}
			return nil
		},
	}
	return req.Do()
}

// Delete will destroy the site.
func (site *Site) Delete() error {
	req := &deleteRequest{
		client: site.srv.client,
		path:   site.getDetailPath(),
	}
	return req.Do()
}

func buildSiteDoc(site *Site) func(outreq *http.Request) (*jsh.Document, error) {
	return func(outreq *http.Request) (*jsh.Document, error) {
		obj, objErr := jsh.NewObject(site.GetID(), sitesResourceType, site)
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
