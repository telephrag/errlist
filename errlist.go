package errlist

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrEmpty = errors.New("")

type ErrNode struct {
	Data map[string]interface{}
	err  error
	next *ErrNode
}

// Constructor. If `nil` is passed internal `err` will be substituted with ErrEmpty.
func New(err error) (self *ErrNode) {
	if errAsErr, ok := err.(*ErrNode); ok {
		return errAsErr
	}

	// to prevent segfault on Unwrap().Error() of childless node with this error inside
	if err == nil {
		err = ErrEmpty
	}

	return &ErrNode{
		Data: make(map[string]interface{}),
		err:  err,
	}
}

// Replaces underlying `err` with `new`.
// If `new` is `nil`, `ErrEmpty` will be used instead.
func (e *ErrNode) ReplaceErr(new error) {
	if new == nil {
		e.err = ErrEmpty
	}
	e.err = new
}

// Checks if `e` has data and non-empty error.
func (e *ErrNode) Empty() bool {
	return len(e.Data) == 0 && e.err == ErrEmpty
}

// Returns true if underlying `err` of some node in the chain is of the same kind as given `err`.
func (e *ErrNode) Has(err error) bool {
	if e.err == err {
		return true
	}

	tail := e
	if tail.next == nil {
		return false
	}

	for tail.next != nil {
		tail = tail.next
		if tail.err == err {
			return true
		}
	}

	return false
}

// Checks if `e` is standalone node or has children.
func (e *ErrNode) HasChildren() bool {
	return e.next == nil
}

// Sets data inside underlying map at `k`.
// Use this to store context of error e.g. timestamp, location etc.
func (e *ErrNode) Set(k string, v interface{}) (self *ErrNode) {
	e.Data[k] = v
	return e
}

// Gets data from underlying map at `k`.
func (e *ErrNode) Get(k string) (v interface{}, ok bool) {
	v, ok = e.Data[k]
	return v, ok
}

// Pushes back `child` to list with head `e`.
// If `child` is not of type `Err`, `New()` is called.
// Should be preferred to standard library methods when it comes to wrapping errors.
func (e *ErrNode) Wrap(child error) (self *ErrNode) {
	tail := e
	for tail.next != nil {
		tail = tail.next
	}

	if childAsErr, ok := child.(*ErrNode); ok {
		tail.next = childAsErr
		return e
	}
	tail.next = New(child)
	return e
}

// Pops back element from the list with head `e`
// If `e.next` is not `nil` returns `next`. Otherwise returs `e`'s underlying `error`.
func (e *ErrNode) Unwrap() error {
	tail := e
	if tail.next == nil {
		return tail.err
	}

	var prev *ErrNode
	for tail.next != nil {
		prev = tail
		tail = tail.next
	}

	res := *tail
	prev.next = nil

	return &res
}

// Same as `Unwrap()` but returns `e` itself if it has no children.
func (e *ErrNode) UnwrapAsNode() *ErrNode {
	tail := e
	if tail.next == nil {
		return tail
	}

	var prev *ErrNode
	for tail.next != nil {
		prev = tail
		tail = tail.next
	}

	res := *tail
	prev.next = nil

	return &res
}

// Returns `e`'s represented as JSON string.
// If `e` is empty returns empty string.
func (e *ErrNode) JSON() string {
	var res string

	if e.Empty() {
		return ""
	}

	if e.err != ErrEmpty {
		res = fmt.Sprintf("\"error\": \"%v\"", e.err) // TODO: use json.Marshal()
	}

	if len(e.Data) > 0 {
		var data string
		for k, v := range e.Data {
			vBytes, err := json.Marshal(v)
			if err != nil {
				data += fmt.Sprintf("\"%s\": \"%s\", ", k, err)
			} else {
				data += fmt.Sprintf("\"%s\": %s, ", k, string(vBytes))
			}
		}
		data = fmt.Sprintf("{%s}", data[:len(data)-2])
		if res != "" {
			res = fmt.Sprintf("%s, \"data\": %s", res, data)
		} else {
			res = fmt.Sprintf("\"data\": %s", data)
		}
	}

	return fmt.Sprintf("{%s}", res)
}

// Proceed to errlist_test.go to see what output will be like.
func (e ErrNode) Error() string {
	res := e.JSON() + "\n"
	err := e
	depth := 0
	for err.next != nil {
		if err.next.Empty() {
			err = *err.next
			continue
		}

		for i := 0; i < depth; i++ {
			res += "    "
		}
		depth++

		res += "  L " + err.next.JSON() + "\n"
		err = *err.next
	}

	return res
}
