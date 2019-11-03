package postgres

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

func formatInBindVars(start, number int) string {
	var numbers []string
	for i := start; i < start+number; i++ {
		numbers = append(numbers, fmt.Sprintf("$%s", strconv.Itoa(i)))
	}
	return fmt.Sprintf("(%s)", strings.Join(numbers, ","))
}

// BuildFilterString builds a postgres where clause from a map of filter params
//
// The allowedFilters arg is to ensure only known filters can be applied and to map to an
// internal postgres filter arg name such as tag.id if needed.
func BuildFilterString(filters map[string][]string, allowedFilters map[string]string) (string, []interface{}) {
	queryString := "WHERE 1=1"
	var varList []interface{}

	i := 1
	for k, valList := range filters {
		if realFilterName, ok := allowedFilters[k]; ok {
			if len(valList) == 0 {
				continue
			}

			queryString = fmt.Sprintf(
				"%s AND %s IN %s",
				queryString,
				realFilterName,
				formatInBindVars(i, len(valList)),
			)
			for _, item := range valList {
				varList = append(varList, item)
			}
			i += len(valList)
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
