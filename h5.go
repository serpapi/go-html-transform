package h5

import (
	"bufio"
	"fmt"
	"os"
	"io"
)

type Attribute struct {
	Name string
	Value string
}

type NodeType int
const (
	TextNode NodeType = iota // zero value so the default
	ElementNode NodeType = iota
	)

type Node struct {
	Type NodeType
	data []int
	Attr []*Attribute
	Parent *Node
	Children []*Node
}

func (n *Node) Data() string {
	return string(n.data)
}

type TokenConsumer func(*Parser, []int)

type Parser struct {
	In *bufio.Reader
	Top *Node
	curr *Node
	consumer TokenConsumer
}

// Handles the various tokenization states
type stateHandler func(p *Parser) (stateHandler, os.Error)

func NewParser(r io.Reader) *Parser {
	return &Parser{In: bufio.NewReader(r)}
}

func (p *Parser) nextInput() (int, os.Error) {
	r, _, err := p.In.ReadRune()
	return r, err
}

func (p *Parser) Parse() os.Error {
	// we start in the data state
	h := handleChar(handleData)
	for h != nil {
		h2, err := h(p)
		if err == os.EOF {
			return nil
		}
		if err != nil {
			// parse error
			return os.NewError(fmt.Sprintf("Parse error: %s", err))
		}
		h = h2
	}
	return nil
}

func textConsumer(p *Parser, chars... int) {
	p.curr.data = append(p.curr.data, chars...) // ugly but safer
}

func handleChar(h func(*Parser, int) stateHandler) stateHandler {
	return func(p *Parser) (stateHandler, os.Error) {
			c, err := p.nextInput()
			if err != nil {
				return nil, err
			}
			return h(p, c), nil
		}
}

// Section 11.2.4.1
func handleData(p *Parser, c int) stateHandler {
	switch c {
	//case '&': // TODO(jwall): do we actually care for this parser?
		//return handleChar(charRefHandler)
	case '<':
		return handleChar(tagOpenHandler)
	default:
		// consume the token
		textConsumer(p, c)
		return handleChar(handleData)
	}
	panic("Unreachable")
}

// Section 11.2.4.2
func charRefHandler(p *Parser, c int) stateHandler {
	switch c {
	case '\t', '\n', '\u000C', ' ', '<', '&':
		// TODO
	case '#':
		// TODO
	default:
		// TODO
	}
	panic("Unreachable")
}

// Section 11.2.4.8
func tagOpenHandler(p *Parser, c int) stateHandler {
	switch c {
	case '!': // markup declaration state
		// TODO
	case '/': // end tag open state
		// TODO
	case '?': // parse error // bogus comment state
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		// TODO
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
		// TODO
	default: // parse error // recover using Section 11.2.4.8 rules
	}
	return nil
}
