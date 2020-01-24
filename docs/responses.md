# Responses
The following are general purpose classes to generate a standard output for resource endpoints.

SingleResponse
---
**`CreateSingleResponse(modelResp interface{}, resourceType string, router *mux.Router, req *http.Request)`**   
This is used to create a response for a single resource.  It will most often be used for Get(one), Post, and Patch.  It 
will also be used but the `CollectionResponse` to generate the list of items.

**modelResp**: \
This is a struct that represents the response for a model.  It is advised that if there are fields
that should be conditionally shown, or fields that need to be formatted in a particular way, a specific response struct 
should be made to separate the display logic from the model.  For example:

`internal/user/model.go`
```
type user struct {
    Id string 
    Username string
    Password string
    Salt string
}
``` 

`internal/user/response.go`
```
type response struct {
    Id string
    Username string 
}
```

In this instance we will use `response` as the **modelResp** instead of the `user`.  The json tag `json:"-"` won't work 
because it won't allow the value to be set through the API.


**resourceType** \
Is the string to be used as the envelope key for the response.

**router** \ 
The main router for the entire main.  It is used to produce the URI's based on the route name.

**request** \ 
Is also required so that the response can generate the full URI in the links. 

CollectionResponse
---
 **`CreateCollectionResponse()`** \
Has near identical parameters as `CreateSingleResponse` for the same reasons.  However the  `collectionMetadata` 
parameter is used to pass information about the collection search such as paging, and search criteria.

Error
---
**`NewErrorResponse(statusCode int, message string)`** \
Is used to generate 400 error response to be encoded in the response.


