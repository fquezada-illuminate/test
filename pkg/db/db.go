package db

import (
	"github.com/fatih/structs"
	"strings"
)

// GetJsonToDbMap returns map of the json field name as the key with the db columns as the value based on the model that
// passed in.
func GetJsonToDbMap(s interface{}) map[string]string {
	columns := make(map[string]string)

	fields := structs.Fields(s)

	for i := range fields {
		dbField := strings.Split(fields[i].Tag("db"), ",")[0]
		if dbField == "" || dbField == "-" {
			continue
		}

		jsonField := strings.Split(fields[i].Tag("json"), ",")[0]
		if jsonField == "" || jsonField == "-" {
			continue
		}

		columns[jsonField] = dbField
	}

	return columns
}

// GetColumns returns an array of the db columns on the model.
func GetColumns(s interface{}) []string {
	var columns []string

	fields := structs.Fields(s)
	for i := range fields {
		dbField := fields[i].Tag("db")
		dbField = strings.Split(dbField, ",")[0]
		if dbField == "" || dbField == "-" {
			continue
		}
		columns = append(columns, dbField)
	}

	return columns
}

// StructToMap returns the a map that is based on the struct that is passed in.  It will use the `structs` tag to
// determine what the key is for the map.
func StructToMap(s interface{}) map[string]interface{} {
	smap := structs.New(s)

	return smap.Map()
}
