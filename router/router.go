package router

const (
	OKResponse          = "+OK\r\n"
	NullResponse        = "$-1\r\n"
	EmptyStringResponse = "+\r\n"
)

type Router interface {
	Lookup(tokens []string) (Callback, bool)
}
