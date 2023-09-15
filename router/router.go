package router

const (
	OKResponse   = "+OK\r\n"
	NullResponse = "$-1\r\n"
)

type Router interface {
	Lookup(tokens []string) (Callback, bool)
}
