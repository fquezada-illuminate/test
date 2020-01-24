# go-rest-service-lib
A library for creating rest services in Golang

How To
---
[Define A Resource Endpoint](docs/define-resource-endpoint.md) \
[Creating A Model](docs/creating-models.md) \
[Middleware](docs/middleware.md) \
[Responses](docs/responses.md) \
[Misc](docs/misc.md)

Project Structure
---
We follow the basic structure projects as defined by [Project Structure](https://github.com/golang-standards/project-layout).

* **cmd**: Application binaries that are broken out by folder.
* **internal**:   Any code that should not be imported by other services.  Most business logic specific to each service should reside in this directory
* **pkg**: functions that can imported by other programs:


3rd Party Libraries
---
**[Pq](https://github.com/lib/pq):** \
A postgresql driver for golang.

**[godotenv](https://github.com/joho/godotenv):** \
Read a .env file into the environment to be used in the application.
 
**[Gorrilla Mux](https://github.com/gorilla/mux):** \
A router to handle requests and direct requests to the correct handlers

**[Ksuid](https://github.com/segmentio/ksuid):** \
A library that is used to generate a partially sequential UUIDS.

**[Logrus](https://github.com/sirupsen/logrus):** \
A library used for logging similar to monolog.

**[Go-playground Validator](https://github.com/go-playground/validator):** \
A struct validation library.  Used for validating requests and models. 

**[DBR](github.com/gocraft/dbr):** \
A lightweight database abstraction layer.    

**[Faith Structs](https://github.com/fatih/structs):** \
Used to convert structs into maps based on tags.  Nested structs such as nullable fields will also be converted to a 
nested map unless tagged otherwise.  The main use of this package is to convert structs to something the DBAL can use
to insert and update values. 

**[Null Library](https://github.com/guregu/null):** \
Used to take allow values to be null.  It properly decodes true/false and int from json when unmarshalled.
