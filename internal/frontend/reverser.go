package frontend

type URLReverser interface {
	Reverse(name string, params ...interface{}) string
}
