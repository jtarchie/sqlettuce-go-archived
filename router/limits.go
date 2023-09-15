package router

import "fmt"

type TokensLimitsRouter struct {
	min, max int
	callback Callback
}

func (t *TokensLimitsRouter) Lookup(tokens []string) (Callback, bool) {
	if 0 < t.min && len(tokens) < t.min {
		return staticResponseCallback(
			fmt.Sprintf(
				"expected minimum number of tokens %d received %d",
				t.min,
				len(tokens),
			),
		), false
	}

	if 0 < t.max && len(tokens) > t.max {
		return staticResponseCallback(
			fmt.Sprintf(
				"expected maximum number of tokens %d received %d",
				t.max,
				len(tokens),
			),
		), false
	}

	return t.callback, true
}

func MinMaxTokensRouter(min, max int, callback Callback) *TokensLimitsRouter {
	return &TokensLimitsRouter{
		min:      min,
		max:      max,
		callback: callback,
	}
}

var _ Router = &TokensLimitsRouter{}
