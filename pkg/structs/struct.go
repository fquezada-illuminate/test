package structs

import (
	fstructs "github.com/fatih/structs"
	"strings"
)

type Helper struct{}

// Satisfy StructHelper Interface
func (h Helper) GetTagMap(s interface{}, key string, value string) (map[string]string, error) {
	columns := make(map[string]string)

	fields := fstructs.Fields(s)

	for i := range fields {
		keyField := strings.Split(fields[i].Tag(key), ",")[0]
		if keyField == "" || keyField == "-" {
			continue
		}

		valueField := strings.Split(fields[i].Tag(value), ",")[0]
		if valueField == "" || valueField == "-" {
			continue
		}

		columns[keyField] = valueField
	}

	return columns, nil
}

func (h Helper) GetTagValues(s interface{}, tag string) []string {
	var values []string

	fields := fstructs.Fields(s)
	for i := range fields {
		tagValue := fields[i].Tag(tag)
		if tagValue == "" || tagValue == "-" {
			continue
		}
		values = append(values, tagValue)
	}

	return values
}

func (h Helper) GetMapByTag(s interface{}, tag string) map[string]interface{} {
	structMap := fstructs.New(s)
	structMap.TagName = tag

	return structMap.Map()
}
