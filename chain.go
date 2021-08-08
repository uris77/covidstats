package covidstats

import "net/http"

// Middleware is a type alias for http.HandlerFunc
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain acts as a list of http.HandlerFunc constructors.
// Chain is immutable.
type Chain struct {
	middlewares []Middleware
}

// NewChain creates a new Chain,
// memorizing the given list of middlewares.
// Middlewares are executed upon a call to Then()
//func NewChain(middlewares ...Middleware) Chain {
//	return Chain{append(([]Middleware)(nil), middlewares...)}
//}

// Chainz applies middlewares to a http.HandlerFunc
func Chainz(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// Then chains the middleware and returns the final http.HandlerFunc.
// A chain can be safely reused by calling Then() several times.
func (c Chain) Then(h http.HandlerFunc) http.HandlerFunc {
	for i := range c.middlewares {
		h = c.middlewares[len(c.middlewares)-1-i](h)
	}
	return h
}

// Append extends a chain, adding the specified middlewares as the last ones
// in the request flow.
// Append returns a new Chain, leaving the original one untouched.
func (c Chain) Append(middlewares ...Middleware) Chain {
	newMiddleware := make([]Middleware, 0, len(c.middlewares)+len(middlewares))
	newMiddleware = append(newMiddleware, c.middlewares...)
	newMiddleware = append(newMiddleware, middlewares...)

	return Chain{newMiddleware}
}

// Extend extends a chain by adding the specified chain as the last on in the request flow.
// Extend returns a new chain, leaving the original one untouched.
func (c Chain) Extend(chains Chain) Chain {
	return c.Append(chains.middlewares...)
}
