package strings

import (
	"fmt"
	"regexp"
)

const (
	Empty         = ""
	CommandPrefix = "/"
)

func IsEmpty(value string) bool {
	return value == Empty
}

func ErrorConcat(err error, layer, origin string) (string, string) {
	return err.Error(), fmt.Sprintf("%s.%s", layer, origin)
}

func ParseStockCodeFromMessage(message string) (string, string) {
	var action, stockCode string
	re := regexp.MustCompile(`/(\w+)=([\w\.]+)`)
	matches := re.FindStringSubmatch(message)

	if len(matches) > 2 {
		action = matches[1]
		stockCode = matches[2]
	}

	return action, stockCode
}
