package postgres

import (
	"fmt"
	"io"
	"log"

	"github.com/jmoiron/sqlx"
)

// BuildFilterString builds a postgres where clause from a map of filter params
//
// The allowedFilters arg is to ensure only known filters can be applied and to map to an
// internal postgres filter arg name such as tag.id if needed.
func BuildFilterString(query string, filters map[string][]string, allowedFilters map[string]string) (string, []interface{}, error) {
	filterString := "WHERE 1=1"
	var inputArgs []interface{}

	// filter key is the url name of the filter used as the lookup for the allowed filters list
	for filterKey, filterValList := range filters {
		if realFilterName, ok := allowedFilters[filterKey]; ok {
			if len(filterValList) == 0 {
				continue
			}

			filterString = fmt.Sprintf("%s AND %s IN (?)", filterString, realFilterName)
			inputArgs = append(inputArgs, filterValList)
		}
	}
	// template the where clause into the original query and then expand the IN clauses with sqlx
	query, args, err := sqlx.In(fmt.Sprintf(query, filterString), inputArgs...)
	if err != nil {
		return "", nil, err
	}

	// not ideal but since we are using txdb as the driver name in tests, need to
	// specify what to rebind ? to
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	return query, args, nil
}

// Close ensures that a deferred close call's error is checked.
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
