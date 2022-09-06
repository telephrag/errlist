# Errlist
`errlist` is implementation of `error` interface as a singly-linked list. Error wrapping is attaching a new element to the end of the list and unwrapping is popping the last element. Additional data can be stored in a map included into each list's element. 

For example usage see `errlist_test.go` and repo called `peary` in my profile. Albeit in the later it's used incorrectly in some places (it's desirable to use package's own functionality instead of storing wrapping error created with `fmt.Errorf()` inside list elements) since I haven't completed refactoring yet.

# Known problems
1. In `Error()` method `json.Marshal()` is used to transform additional data elements into strings while giving quotation marks only to those types that actually need them. This can result in error which will be written to ouput string instead of value that was meant to be converted.
2. As of now `Error()` outputs list containing more than one element as pretty pyramid that is easy to read (run `go test` to get better understanding). This form can be hard to parse so, I'll probably "stringify" list as json array in the future.   