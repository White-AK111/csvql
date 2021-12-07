package parser

import (
	"fmt"
	"strings"
)

type IParser interface {
}

type Parser struct {
	notation []string
}

func (a *Parser) GetConditions(query string) error {
	var priorities = map[string]uint8{
		"||":  1,
		"&&":  1,
		"and": 1,
		"or":  1,
		">":   2,
		"<":   2,
		">=":  2,
		"<=":  2,
		"=":   2,
		"==":  2,
		"!=":  2,
	}

	query = strings.ToLower(query)
	query = strings.TrimSpace(query)

	nop, n := 0, len(query)
	operators := make([]string, n)
	var notation []string

	i := 0
	popPushOp := func(op string) {
		priority := priorities[op]
		for nop > 0 && priorities[operators[nop-1]] >= priority {
			nop--
			notation = append(notation, operators[nop])
		}
		operators[nop] = op
		nop++
		i++
	}

	for i < n {
		c := query[i]
		switch c {
		case ' ':
			i++
		case ')':
			for nop > 0 && operators[nop-1] != "(" {
				nop--
				notation = append(notation, operators[nop])
			}
			if nop == 0 || operators[nop-1] != "(" {
				return fmt.Errorf("'%v' has no '(' found for ')' at %v", query, i)
			}
			nop--
			if nop > 0 && operators[nop-1] == "!" {
				notation = append(notation, "!")
				nop--
			}
			i++
		case '(':
			operators[nop] = "("
			nop++
			i++
		case '*', '/', '%', '+', '-':
			popPushOp(string(c))
		case '!':
			next := query[i+1]
			if next == '(' {
				operators[nop] = string(c)
				nop++
				i++
			} else if next == '=' {
				i++
				popPushOp(string([]byte{c, next}))
			} else {
				return fmt.Errorf("'%v' has invalid token at %v: %v", query, i+1, next)
			}
		case '>', '<':
			op := []byte{c}
			if query[i+1] == '=' {
				op = append(op, '=')
				i++
			}
			popPushOp(string(op))
		case '|', '&', '=':
			next := query[i+1]
			if next != c {
				return fmt.Errorf("'%v' has invalid token at %v: %v", query, i+1, next)
			}
			i++
			popPushOp(string([]byte{c, next}))
		default:
			var word []byte
			for i < n {
				c = query[i]
				word = append(word, c)
				i++
			}
			notation = append(notation, string(word))
		}
	}

	for nop > 0 {
		nop--
		if op := operators[nop]; op != "(" {
			notation = append(notation, op)
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
