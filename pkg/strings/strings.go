package strings

import (
	"fmt"
)

const Empty = ""

func IsEmpty(value string) bool {
	return value == Empty
}

func ErrorConcat(err error, layer, origin string) (string, string) {
	return err.Error(), fmt.Sprintf("%s.%s", layer, origin)
}
