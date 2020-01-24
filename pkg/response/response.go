package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/illuminateeducation/rest-service-lib-go/pkg/http/route"
	"net/http"
	"reflect"
	"strings"
)

type Link struct {
	Method string `json:"method"`
	Rel    string `json:"rel"`
	Href   string `json:"href"`
}

type ResourceLinks struct {
	router  *mux.Router
	request *http.Request
	links   []Link
}

func (rl *ResourceLinks) getRouter() *mux.Router {
	return rl.router
}

func (rl *ResourceLinks) addLink(l Link) {
	rl.links = append(rl.links, l)
}

func (rl *ResourceLinks) getRequest() *http.Request {
	return rl.request
}

func (rl *ResourceLinks) GetLinks() []Link {
	return rl.links
}

type Resource struct {
	modelResponse interface{}
	resourceType  string
}

func (r *Resource) getModelResponse() interface{} {
	return r.modelResponse
}

func (r *Resource) GetResourceType() string {
	return r.resourceType
}

type SingleResponse struct {
	ResourceLinks
	Resource
}

func (rl *ResourceLinks) AddLink(method string, rel string, routeName string, routeParams map[string]string) error {
	routePtr := rl.getRouter().Get(routeName)
	if routePtr == nil {
		return errors.New(fmt.Sprintf("Route %s does not exist", routeName))
	}

	flatParams := []string{}
	for key, value := range routeParams {
		flatParams = append(flatParams, key)
		flatParams = append(flatParams, value)
	}

	url, err := routePtr.URL(flatParams...)
	if err != nil {
		return err
	}

	url.Host = rl.getRequest().Host
	url.Scheme = "http"
	if rl.request.URL.Scheme != "" {
		url.Scheme = rl.getRequest().URL.Scheme
	}
	rl.addLink(Link{
		Method: strings.ToUpper(method),
		Rel:    strings.ToLower(rel),
		Href:   url.String(),
	})

	return nil
}

// CreateSingleResponse instantiates a SingleResponse
func CreateSingleResponse(modelResp interface{}, resourceType string, router *mux.Router, req *http.Request) (SingleResponse, error) {
	if resourceType == "" {
		return SingleResponse{}, errors.New("resource Type cannot be empty for Single Responses")
	}

	return SingleResponse{
		ResourceLinks: ResourceLinks{
			router:  router,
			request: req,
		},
		Resource: Resource{
			modelResponse: modelResp,
			resourceType:  resourceType,
		},
	}, nil
}

func (sr SingleResponse) MarshalJSON() ([]byte, error) {
	responseObj := make(map[string]interface{}, 0)
	responseObj[sr.GetResourceType()] = sr.getModelResponse()

	if len(sr.GetLinks()) > 0 {
		responseObj["links"] = sr.GetLinks()
	}

	return json.Marshal(responseObj)
}

type CollectionPaging struct {
	Current int `json:"current"`
	Size    int `json:"size"`
	Pages   int `json:"pages"`
}

type CollectionMetadata struct {
	Count  int                    `json:"count"`
	Paging CollectionPaging       `json:"paging"`
	Sort   map[string]interface{} `json:"sort"`
	Filter map[string]interface{} `json:"filter"`
}

type CollectionResponse struct {
	metadata CollectionMetadata
	items    []SingleResponse
	ResourceLinks
	resourceType string
}

// CreateCollectionResponse instantiates a CollectionResponse
func CreateCollectionResponse(md CollectionMetadata, resourceType string, router *mux.Router, req *http.Request) (CollectionResponse, error) {
	if resourceType == "" {
		return CollectionResponse{}, errors.New("resource Type cannot be empty for Collection Responses")
	}
	return CollectionResponse{
		metadata: md,
		items:    []SingleResponse{},
		ResourceLinks: ResourceLinks{
			router:  router,
			request: req,
		},
		resourceType: resourceType,
	}, nil
}

// AddItems add a slice of items to the collection response
func (cr *CollectionResponse) AddItem(sr SingleResponse) {
	cr.items = append(cr.items, sr)
}

func (cr CollectionResponse) MarshalJSON() ([]byte, error) {
	responseObj := make(map[string]interface{}, 0)
	responseObj["items"] = cr.items
	responseObj["metadata"] = cr.metadata

	if len(cr.GetLinks()) > 0 {
		responseObj["links"] = cr.GetLinks()
	}

	return json.Marshal(responseObj)
}

// NewModelSingleResponse creates a new SingleResponse specifically for instance objects
func NewModelSingleResponse(model interface{}, rm map[string]string, resourceType string, router *mux.Router, req *http.Request) (SingleResponse, error) {
	errs := make([]string, 0)
	sr, _ := CreateSingleResponse(model, resourceType, router, req)

	routeParams := mux.Vars(req)
	if routeParams == nil {
		routeParams = make(map[string]string)
	}
	s := reflect.ValueOf(model)
	if s.Kind() == reflect.Struct {
		routeParams["id"] = s.FieldByName("Id").String()
	}


	if _, ok := rm[route.CGET_ROUTE]; ok {
		if linkErr := sr.AddLink("GET", "parent", rm[route.CGET_ROUTE], routeParams); linkErr != nil {
			errs = append(errs, linkErr.Error())
		}
	}

	if _, ok := rm[route.GET_ROUTE]; ok {
		if linkErr := sr.AddLink("GET", "self", rm[route.GET_ROUTE], routeParams); linkErr != nil {
			errs = append(errs, linkErr.Error())
		}
	}

	if _, ok := rm[route.PATCH_ROUTE]; ok {
		if linkErr := sr.AddLink("PATCH", "update", rm[route.PATCH_ROUTE], routeParams); linkErr != nil {
			errs = append(errs, linkErr.Error())
		}
	}

	if _, ok := rm[route.DELETE_ROUTE]; ok {
		if linkErr := sr.AddLink("DELETE", "delete", rm[route.DELETE_ROUTE], routeParams); linkErr != nil {
			errs = append(errs, linkErr.Error())
		}
	}

	if len(errs) > 0 {
		return sr, errors.New(strings.Join([]string(errs), " || "))
	}

	return sr, nil
}

// NewConfigCollectionResponse creates a new CollectionResponse specifically for instance objects
func NewModelCollectionResponse(cm CollectionMetadata, rm map[string]string, resourceType string, router *mux.Router, req *http.Request) (CollectionResponse, error) {
	errs := make([]string, 0)
	cr, _ := CreateCollectionResponse(cm, resourceType, router, req)

	routeParams := mux.Vars(req)

	if _, ok := rm[route.CGET_ROUTE]; ok {
		if linkErr := cr.AddLink("GET", "self", rm[route.CGET_ROUTE], routeParams); linkErr != nil {
			errs = append(errs, linkErr.Error())
		}

		if linkErr := cr.AddLink("GET", "first", rm[route.CGET_ROUTE], routeParams); linkErr != nil {
			errs = append(errs, linkErr.Error())
		}

		if linkErr := cr.AddLink("GET", "last", rm[route.CGET_ROUTE], routeParams); linkErr != nil {
			errs = append(errs, linkErr.Error())
		}
	}

	if _, ok := rm[route.POST_ROUTE]; ok {
		if linkErr := cr.AddLink("POST", "create", rm[route.POST_ROUTE], routeParams); linkErr != nil {
			errs = append(errs, linkErr.Error())
		}
	}

	if len(errs) > 0 {
		return cr, errors.New(strings.Join([]string(errs), " || "))
	}

	return cr, nil
}
