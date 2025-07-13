package assert

import (
	"fmt"
	"log"
)

func Assert(condition bool, message any) {
	if condition {
		return
	}

	switch msg := message.(type) {
	case error:
		log.Fatalln(msg.Error())
	default:
		log.Fatalln(msg)
	}
}

func Assertf(condition bool, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	Assert(condition, message)
}

func AssertErr(err error) {
	Assert(err == nil, err)
}
