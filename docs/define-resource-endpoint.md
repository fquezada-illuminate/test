# Define Resource Endpoint

1. [Create Service For Resource](#create-service-for-resource)
2. [Create Handlers](#create-handlers)
3. [Add AttachRoutes Function To Service](#add-attachroutes-function-to-service)
4. [Add Routes to Main Router](#add-routes-to-main-router)


<a name="create-service-for-resource">Create Service For Resource</a>
---
A service struct is used to define all of the controller actions and its dependencies.  The struct should be defined in
the `internal/[resoure]` directory.  An example definition would be.

`internal/instance/svc.go`
```
type Svc struct {
	router     *mux.Router
	repository db.Repository
	validator  *validation.Validator
}
```

<a name="create-handlers">Create Handlers</a>
---
Define all of the handlers (controller actions) that belong to this service.

`internal/instance/svc.go`
```
func (is Svc) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}
func (is Svc) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
	}
}
...
```


<a name="add-attachroutes-function-to-service">Add AttachRoutes Function To Service</a>
---
Create a map of RouteNames.  This will be used to automatically generate links in our responses.

`internal/instance/svc.go`
```
var RouteNames = map[string]string{
	route.CGET_ROUTE:   "cget_instance",
	route.GET_ROUTE:    "get_instance",
	route.POST_ROUTE:   "post_instance",
	route.PATCH_ROUTE:  "patch_instance",
	route.DELETE_ROUTE: "delete_instance",
}
```

The svc class should be responsible for generating its own routes.  As such add a function to the svc called `AttachRoutes`

`internal/instance/svc.go`
```
func (is Svc) AttachRoutes(r *mux.Router) *mux.Router {
    // Create a subrouter that will have route handlers added to them
    sr := r.PathPrefix("/instances").Subrouter()
    sr.Use(middleware.Token)
	
    // add route handlers here
    sr.Path("/{id}").Methods("GET").Handler(is.GetHandler()).Name(RouteNames[route.GET_ROUTE])
    sr.Path("").Methods("POST").Handler(is.PostHandler()).Name(RouteNames[route.POST_ROUTE])

    // we want to return the subrouter so that we can chain sub-resources to it if necessary 
    return sr
}

```
Within this function you can add middleware that will apply to all of routes.  `sr.Use(middleware.Token)` is an example 
of that. 

<a name="add-routes-to-main-router">Add Routes to Main Router</a>
---
In main.go instantiate the newly created service and attach the routes to the global router.

`main.go`
```
r := mux.NewRouter()

instanceSvc := instance.svc {
    //dependencies
}

instanceSvc.AttachRoutes(r)
```
