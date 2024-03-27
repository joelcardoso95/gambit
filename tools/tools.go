package tools

import (
	"fmt"
	"time"
)

func DateMySQL() string {
	time := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		time.Year(), time.Month(), time.Day(), time.Hour(), time.Minute(), time.Second())

}
