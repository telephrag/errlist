package errlist

type errData interface {
	~int | ~string
}

func GetData[T errData](list *ErrNode, k string) (v T, ok bool) {
	v, ok = list.Data[k].(T)
	return v, ok
}
