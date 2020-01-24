# Middleware
Middleware are used to run code before and after the handler code.

Client Id
---
Will validate that a `x-ied-client-id` header is set and that it is a valid UUIDv4.

Content-Type
---
Will validate that a `Content-Type` header is set to the specified type.  It will also set the Content-Type on the 
response to the configured Content-Type.  There is a special handler `JsonContentType` that is for convenience of 
checking and setting the Content-Type to Json.

Token
---
Will validate that a `x-ied-service-token` header is set and that it is equal to the environment variable 
`SERVICE_TOKEN`.

Recover
---
Will catch any panics that occur and will return a 500 error json response and recover to keep the application 
running. 

Validation
---
Will attempt to pre-validate a request based on the model that is passed in.