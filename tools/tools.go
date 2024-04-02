package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func DateMySQL() string {
	time := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second())

}

func SkipString(value string) string {
	description := strings.ReplaceAll(value, "'", "")
	description = strings.ReplaceAll(description, "\"", "")
	return description
}

func AdjustQuery(query string, fieldName string, typeField string, valueN int, valueF float64, valueS string) string {
	if (typeField == "S" && len(valueS) < 1) || (typeField == "F" && valueF == 0) || (typeField == "N" && valueN == 0) {
		return query
	}

	if !strings.HasSuffix(query, "SET ") {
		query += ", "
	}

	switch typeField {
	case "S":
		query += fieldName + " = '" + SkipString(valueS) + "'"
	case "N":
		query += fieldName + " = " + strconv.Itoa(valueN)
	case "F":
		query += fieldName + " = " + strconv.FormatFloat(valueF, 'e', -1, 64)
	}

	return query
}
