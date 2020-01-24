# Creating A Model
Overall creating a model is no different than defining a struct.  Go tags are needed on each field to 
determine database columns, validation etc.

```
type Config struct {
	Id        string             `json:"id" db:"id" structs:"id" validate:"required,uuid4"`
	Type      types.NullString   `json:"type" db:"type" structs:"type,omitnested" validate:"min=1,max=255"`
	Key       types.NullString   `json:"key" db:"key" structs:"key,omitnested" validate:"notblank=1 255"`
	Value     types.NullString   `json:"value" db:"value" structs:"value,omitnested" validate:"notblank=1 255"`
	ClientId  types.NullString   `json:"clientId" db:"client_id" structs:"client_id,omitnested" validate:"required,uuid4"`
	CreatedAt types.NullDatetime `json:"createdAt" db:"created_at" structs:"created_at,omitnested"`
	UpdatedAt types.NullDatetime `json:"updatedAt" db:"updated_at" structs:"updated_at,omitnested"`
}
``` 

Tags
---
**[json](https://golang.org/pkg/encoding/json/#Marshal)** \
Defines the key within a json payload that will be mapped to the tagged struct member.

**[db](https://github.com/gocraft/dbr)** \
Defines the column name that the struct will use to map back and forth from the database.

**[structs](https://github.com/fatih/structs)** \
Used to map structs to map\[string\]interface{} for update statements.  There is a lot of overlap between db and structs.  
We would like to combine them together but at this point they have to be separate.  When using `driver.valuer` interfaces
such as *types.NullString* or *types.NullDatetime*,  `omitnested` must be added to the tag otherwise the library will try 
decode the embedded struct into a map as well.

**[validate](https://github.com/go-playground/validator)** \
Used to define the validation for the struct member.  See library documentation for more detail.


Types
---
In general when the struct member is settable through the api, the type should be types.Null*.  When a *null* passed in 
through json, a non-Null type will be ignored.  Because of this we cannot return an error to the consumer.

**Datetime, NullDatetime, Date, and NullDate**: Are used to properly marshal and unmarshal into the appropriate format. 