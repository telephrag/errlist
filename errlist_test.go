package errlist

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestDev(t *testing.T) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	top := New(errors.New("a usefull error msg")).Set("location", "this location")
	middle := New(errors.New("another usefull error msg")).Set("code", 1337)
	bot := errors.New("yet again usefull error message")
	top.Wrap(middle).Wrap(bot).Wrap(nil)

	top.Set("timestamp", time.Unix(234567890000, 0))

	log.Print(top)

	const top_Error = "" +
		"{\"error\": \"a usefull error msg\", \"data\": {\"location\": \"this location\", \"timestamp\": \"9403-03-01T07:13:20+04:00\"}}" +
		"\n  L {\"error\": \"another usefull error msg\", \"data\": {\"code\": 1337}}" +
		"\n      L {\"error\": \"yet again usefull error message\"}\n"
	top.Unwrap().Error() // segfault test
	if top.Error() != top_Error {
		t.Errorf("\nexpected: %s \nreceived: %s\n", top_Error, top.Error())
	}

	top.UnwrapAsNode().UnwrapAsNode().UnwrapAsNode().Error() // segfault test
}
