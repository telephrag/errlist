package errlist

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// Generally checking possible SEGFAULTs by chaining errors of various types
// and probing different use cases
func TestSegfaults(t *testing.T) {
	top := New(fmt.Errorf("a usefull error msg")).Set("location", "this location")
	middle := New(fmt.Errorf("another usefull error msg")).Set("code", "1337")
	bot := fmt.Errorf("yet again usefull error message")
	var hell error = nil

	top.Wrap(middle).Wrap(bot).Wrap(hell) // wrap various kinds of `errNode`
	log.Print(top.UnwrapAsNode().json())

	top.UnwrapAsNode().Unwrap().Error() // call `Error()` on underlying `error`
	top.Error()                         // call `Error()` on childless `errNode`

	uw := top.UnwrapAsNode() // try getting value at non-existent key
	uw.Get("code")

	// Check if UnwrapAsNode() actually returns top
	top.UnwrapAsNode().UnwrapAsNode().UnwrapAsNode().UnwrapAsNode().UnwrapAsNode()

}

func TestMarshalJSON(t *testing.T) {
	top := New(fmt.Errorf("a usefull error msg")).Set("location", "this location")
	middle := New(fmt.Errorf("another usefull error msg")).Set("code", 1337)
	bot := fmt.Errorf("yet again usefull error message")
	var hell error = nil

	top.Wrap(middle).Wrap(bot).Wrap(hell) // wrap various kinds of `errNode`

	const top_MarshalList = "[{\"error\":\"a usefull error msg\",\"data\":{\"location\":\"this location\"}},{\"error\":\"another usefull error msg\",\"data\":{\"code\":1337}},{\"error\":\"yet again usefull error message\",\"data\":{}},{\"error\":\"\",\"data\":{}}]"
	bytes, err := top.MarshalList()
	if err != nil {
		t.Errorf("expected `err` to be nil, got: %v", err)
	}
	if string(bytes) != top_MarshalList {
		t.Errorf("expected `string(bytes)` to be equal to `top_MarshalList`, got: %s", string(bytes))
	}

	const top_PrettyStr = "[\n" +
		"    {\"error\":\"a usefull error msg\",\"data\":{\"location\":\"this location\"}},\n" +
		"    {\"error\":\"another usefull error msg\",\"data\":{\"code\":1337}},\n" +
		"    {\"error\":\"yet again usefull error message\",\"data\":{}},\n" +
		"    {\"error\":\"\",\"data\":{}},\n]"
	prettyStr := top.ErrorList()
	if prettyStr != top_PrettyStr {
		t.Errorf("expected `prettyStr` to be equal to `top_PrettyStr`, got: %s", prettyStr)
	}
}

func generateValue() int64 {
	val := rand.Int63() % 10000
	ng := int64(runtime.NumGoroutine())

	work := md5.Sum([]byte(fmt.Sprint(ng + val)))    // h(m)
	work = md5.Sum([]byte(fmt.Sprint(work, ng+val))) // h(h(m) || m)
	return int64(binary.BigEndian.Uint16(work[:2])) % 10000
}

func TestBenchmark(t *testing.T) {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		err := New(fmt.Errorf("error message")).Set("value", generateValue())
		err.Error()
	}
	elapsed := time.Since(start)
	fmt.Printf("finished in %d us\n", elapsed.Microseconds())

	start = time.Now()
	for i := 0; i < 1000; i++ {
		err := New(fmt.Errorf("error message")).Set("value", generateValue())
		err.ErrorList()
	}
	elapsed = time.Since(start)
	fmt.Printf("finished in %d us\n", elapsed.Microseconds())
}
