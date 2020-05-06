// Package parse provides a simple parser combinator library
package parse // import "github.com/spaskalev/misc/parse"

// The N struct is a parser result node.
type N struct {
	// Match indicates whether the parser succeeded.
	Matched bool
	// Content contains whatever was matched by the parser.
	Content string
	// Nodes contains any result nodes by nested parsers.
	Nodes []N
}

// Parser function type
// Takes a string and returns a result and the remaining part of the string.
type P func(string) (N, string)

// A sequence of parsers. Matches when all of p match.
func Seq(p ...P) P {
	return func(s string) (N, string) {
		result := N{Matched: true, Nodes: make([]N, 0, len(p))}
		for _, parser := range p {
			n, r := parser(s)
			result.Nodes = append(result.Nodes, n)
			if !n.Matched {
				result.Matched = false
				break
			}
			result.Content = result.Content + n.Content
			s = r
		}
		return result, s
	}
}

// Matches and returns on the first match of p.
func Any(p ...P) P {
	return func(s string) (N, string) {
		result := N{Matched: false, Nodes: nil}
		for _, parser := range p {
			n, r := parser(s)
			if n.Matched {
				return n, r
			}
			result.Nodes, s = append(result.Nodes, n), r
		}
		return result, s
	}
}

// Kleene star (zero or more). Always matches.
func K(p P) P {
	return func(s string) (N, string) {
		result := N{Matched: true, Nodes: nil}
		for n, r := p(s); n.Matched; n, r = p(r) {
			result.Content = result.Content + n.Content
			result.Nodes = append(result.Nodes, n)
			s = r
		}
		return result, s
	}
}

// Returns a delegating parser whose delegate can be set on later.
// Useful for recursive definitions.
func Defer() (P, *P) {
	var deferred P
	return func(s string) (N, string) {
		return deferred(s)
	}, &deferred
}

// Returns a parser that accepts the specified number of bytes
func Accept(count int) P {
	return func(s string) (N, string) {
		if len(s) < count {
			return N{Matched: false, Content: "", Nodes: nil}, s
		}
		return N{Matched: true, Content: s[:count], Nodes: nil}, s[count:]
	}
}

// Returns a parser that accepts the specified string
func String(value string) P {
	accept := Accept(len(value))
	return func(s string) (N, string) {
		n, r := accept(s)
		if !n.Matched {
			return n, s
		}
		if n.Content == value {
			return n, r
		}
		return N{Matched: false, Content: "", Nodes: []N{n}}, s
	}
}

// Returns a parser that matches any single digit
// It could have been written as Any(String("0"),String("1"),...,String("9")) as well :)
func Digit() P {
	accept := Accept(1)
	return func(s string) (N, string) {
		n, r := accept(s)
		if !n.Matched {
			return n, s
		}
		if n.Content[0] >= 48 && n.Content[0] <= 57 {
			return n, r
		}
		return N{Matched: false, Content: "", Nodes: []N{n}}, s
	}
}
