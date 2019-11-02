package postgres

import (
	"fmt"
	"io"
	"log"
)

// BuildFilterString builds a postgres where clause from a map of filter params
//
// The allowedFilters arg is to ensure only known filters can be applied and to map to an
// internal postgres filter arg name such as tag.id if needed.
func BuildFilterString(filters map[string]string, allowedFilters map[string]string) (string, []interface{}) {
	queryString := "WHERE 1=1"
	var varList []interface{}

	i := 0
	for k, v := range filters {
		if val, ok := allowedFilters[k]; ok {
			i++
			queryString = fmt.Sprintf(" %s AND %s = $%v", queryString, val, i)
			varList = append(varList, v)
		}
	}
	return queryString, varList
}

// Close ensures that a deferred close call's error is checked.
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
