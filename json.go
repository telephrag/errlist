package errlist

import (
	"encoding/json"
	"fmt"
)

func (e *ErrNode) MarshalJSON() ([]byte, error) {
	type temp struct {
		Err  string                 `json:"error"`
		Data map[string]interface{} `json:"data"`
	}

	obj := temp{
		Err:  e.Err.Error(),
		Data: e.Data,
	}

	return json.Marshal(obj)
}

func (e *ErrNode) MarshalList() ([]byte, error) {
	err := e
	count := 1
	for err.next != nil {
		count++
		err = err.next
	}

	err = e
	listAsArr := make([]*ErrNode, count)
	for i := range listAsArr {
		listAsArr[i] = err
		err = err.next
	}

	return json.Marshal(listAsArr)
}

func (e *ErrNode) ErrorList() string {
	res := ""
	err := e
	for err != nil {
		errBytes, _ := json.Marshal(err)
		res += "    " + string(errBytes) + ",\n"
		err = err.next
	}

	return fmt.Sprintf("[\n%s]", res)
}
