package kel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/brosner/go-json-spec-handler"
)

// ErrNotFound is the error resulting if the requested object returned a 404
var ErrNotFound = errors.New("object requested was not found")

// CreateRequest represents the behavior of creating a Kel object.
type CreateRequest interface {
	Do() error
}

// ListRequest represents the behavior of listing Kel objects.
type ListRequest interface {
	Include() ListRequest
	Do() error
}

// GetRequest represents the behavior of getting a Kel object.
type GetRequest interface {
	Include() GetRequest
	Do() error
}

// UpdateRequest represents the behavior of updating a Kel object.
type UpdateRequest interface {
	Do() error
}

// DeleteRequest represents the behavior of deleting a Kel object.
type DeleteRequest interface {
	Do() error
}

type createRequest struct {
	client *Client
	path   string
	object Object
	hdr    http.Header
}

func (req *createRequest) Do() error {
	outreq, err := http.NewRequest("POST", req.client.makeURL(req.path).String(), nil)
	if err != nil {
		return err
	}
	obj, objErr := jsh.NewObject(req.object.GetID(), req.object.GetResourceType(), req.object)
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
	if req.hdr != nil {
		outreq.Header = req.hdr
	}
	outreq.Header.Set("Content-Type", jsh.ContentType)
	outreq.Body = jsh.CreateReadCloser(docSerialized)
	outreq.ContentLength = int64(len(docSerialized))
	res, err := req.client.Do(outreq)
	if err != nil {
		return err
	}
	parser := &jsh.Parser{Method: "", Headers: res.Header}
	doc, docErr := parser.Document(res.Body, jsh.ObjectMode)
	if docErr != nil {
		return docErr
	}
	doc.Status = res.StatusCode
	switch res.StatusCode {
	case http.StatusCreated:
		obj = doc.Data[0]
		if objErr := obj.Unmarshal(req.object.GetResourceType(), req.object); objErr != nil {
			return objErr
		}
		return nil
	case http.StatusBadRequest:
		return errors.New(doc.Errors[0].Detail)
	default:
		return fmt.Errorf("unknown response from API: %s", res.Status)
	}
}

type listRequest struct {
	client  *Client
	path    string
	handler func(obj *jsh.Document) error
	hdr     http.Header
}

// Include ...
func (req *listRequest) Include() ListRequest {
	return req
}

// Do ...
func (req *listRequest) Do() error {
	outreq, err := http.NewRequest("GET", req.client.makeURL(req.path).String(), nil)
	if err != nil {
		return err
	}
	if req.hdr != nil {
		outreq.Header = req.hdr
	}
	outreq.Header.Set("Content-Type", jsh.ContentType)
	res, err := req.client.Do(outreq)
	if err != nil {
		return err
	}
	parser := &jsh.Parser{Method: "", Headers: res.Header}
	doc, docErr := parser.Document(res.Body, jsh.ListMode)
	if docErr != nil {
		return docErr
	}
	doc.Status = res.StatusCode
	switch res.StatusCode {
	case http.StatusOK:
		if err := req.handler(doc); err != nil {
			return err
		}
		return nil
	case http.StatusBadRequest:
		return errors.New(doc.Errors[0].Detail)
	default:
		return fmt.Errorf("unknown response from API: %s", res.Status)
	}
}

type getRequest struct {
	client  *Client
	path    string
	handler func(obj *jsh.Document) error
	hdr     http.Header
}

// Include ...
func (req *getRequest) Include() GetRequest {
	return req
}

// Do ...
func (req *getRequest) Do() error {
	outreq, err := http.NewRequest("GET", req.client.makeURL(req.path).String(), nil)
	if err != nil {
		return err
	}
	if req.hdr != nil {
		outreq.Header = req.hdr
	}
	outreq.Header.Set("Content-Type", jsh.ContentType)
	res, err := req.client.Do(outreq)
	if err != nil {
		return err
	}
	parser := &jsh.Parser{Method: "", Headers: res.Header}
	doc, docErr := parser.Document(res.Body, jsh.ObjectMode)
	if docErr != nil {
		return docErr
	}
	doc.Status = res.StatusCode
	switch res.StatusCode {
	case http.StatusOK:
		if err := req.handler(doc); err != nil {
			return err
		}
		return nil
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusBadRequest:
		return errors.New(doc.Errors[0].Detail)
	default:
		return fmt.Errorf("unknown response from API: %s", res.Status)
	}
}

type updateRequest struct {
	client *Client
	path   string
	object Object
	hdr    http.Header
}

func (req *updateRequest) Do() error {
	outreq, err := http.NewRequest("PATCH", req.client.makeURL(req.path).String(), nil)
	if err != nil {
		return err
	}
	obj, objErr := jsh.NewObject(req.object.GetID(), req.object.GetResourceType(), req.object)
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
	if req.hdr != nil {
		outreq.Header = req.hdr
	}
	outreq.Header.Set("Content-Type", jsh.ContentType)
	outreq.Body = jsh.CreateReadCloser(docSerialized)
	outreq.ContentLength = int64(len(docSerialized))
	res, err := req.client.Do(outreq)
	if err != nil {
		return err
	}
	parser := &jsh.Parser{Method: "", Headers: res.Header}
	doc, docErr := parser.Document(res.Body, jsh.ObjectMode)
	if docErr != nil {
		return docErr
	}
	doc.Status = res.StatusCode
	switch res.StatusCode {
	case http.StatusOK:
		obj = doc.Data[0]
		if objErr := obj.Unmarshal(req.object.GetResourceType(), req.object); objErr != nil {
			return objErr
		}
		return nil
	case http.StatusBadRequest:
		return errors.New(doc.Errors[0].Detail)
	default:
		return fmt.Errorf("unknown response from API: %s", res.Status)
	}
}

type deleteRequest struct {
	client *Client
	path   string
	hdr    http.Header
}

func (req *deleteRequest) Do() error {
	outreq, err := http.NewRequest("DELETE", req.client.makeURL(req.path).String(), nil)
	if err != nil {
		return err
	}
	if req.hdr != nil {
		outreq.Header = req.hdr
	}
	outreq.Header.Set("Content-Type", jsh.ContentType)
	res, err := req.client.Do(outreq)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	parser := &jsh.Parser{Method: "", Headers: res.Header}
	doc, docErr := parser.Document(res.Body, jsh.ObjectMode)
	if docErr != nil {
		return docErr
	}
	doc.Status = res.StatusCode
	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted:
		return nil
	case http.StatusBadRequest:
		return errors.New(doc.Errors[0].Detail)
	default:
		return fmt.Errorf("unknown response from API: %s", res.Status)
	}
}
