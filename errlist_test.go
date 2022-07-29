package errlist

import (
	"fmt"
	"log"
	"testing"
)

// Checking methods responsible for converting `errNode` into string.
func TestError(t *testing.T) {
	const top_Error = "{\"error\": \"a usefull error msg\", \"data\": {\"location\": \"this location\"}}" +
		"\n  L {\"error\": \"another usefull error msg\", \"data\": {\"code\": \"1337\"}}" +
		"\n      L {\"error\": \"yet again usefull error message\"}\n"

	top := New(fmt.Errorf("a usefull error msg")).Set("location", "this location")
	middle := New(fmt.Errorf("another usefull error msg")).Set("code", "1337")
	bot := fmt.Errorf("yet again usefull error message")

	top.Wrap(middle).Wrap(bot)

	if top.Error() != top_Error {
		t.Errorf("\nexpected: %s \nreceived: %s\n", top_Error, top.Error())
	}
}

// Generally checking possible SEGFAULTs by chaining errors of various types
// and probing different use cases
func TestSegfaults(t *testing.T) {
	top := New(fmt.Errorf("a usefull error msg")).Set("location", "this location")
	middle := New(fmt.Errorf("another usefull error msg")).Set("code", "1337")
	bot := fmt.Errorf("yet again usefull error message")
	var hell error = nil

	top.Wrap(middle).Wrap(bot).Wrap(hell) // wrap various kinds of `errNode`
	log.Print(top.JSON())
	log.Print(top.UnwrapAsNode().json())

	top.UnwrapAsNode().Unwrap().Error() // call `Error()` on underlying `error`
	top.Error()                         // call `Error()` on childless `errNode`

	uw := top.UnwrapAsNode() // try getting value at non-existent key
	str, _ := uw.Get("code")
	str += "some string"

	// Check if UnwrapAsNode() actually returns top
	top.UnwrapAsNode().UnwrapAsNode().UnwrapAsNode().UnwrapAsNode().UnwrapAsNode()

}
