package router

import "fmt"

type tokensLimits struct {
	min, max int
	callback Callback
}

func (t *tokensLimits) Lookup(tokens []string) (Callback, error) {
	if 0 < t.min && len(tokens) < t.min {
		return nil, fmt.Errorf("expected minimum number of tokens %d received %d: %w", t.min, len(tokens), ErrIncorrectTokens)
	}

	if 0 < t.max && len(tokens) > t.max {
		return nil, fmt.Errorf("expected maximum number of tokens %d received %d: %w", t.max, len(tokens), ErrIncorrectTokens)
	}

	return t.callback, nil
}

func minMaxTokens(min, max int, callback Callback) *tokensLimits {
	return &tokensLimits{
		min:      min,
		max:      max,
		callback: callback,
	}
}

var _ Router = &tokensLimits{}
