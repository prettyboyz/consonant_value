package toml

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type parser struct {
	mapping map[string]interface{}
	types   map[string]tomlType
	lx      *lexer

	// A list of keys in the order that they appear in the TOML data.
	ordered []Key

	// the full key for the current hash in scope
	context Key

	// the base key name for everything except hashes
	currentKey string

	// rough approximation of line number
	approxLine int

	// A map of 'key.group.names' to whether they were created implicitly.
	implicits map[string]bool
}

type parseError string

func (pe parseError) Error() string {
	return string(pe)
}

func parse(data string) (p *parser, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(parseError); ok {
				return
			}
			panic(r)
		}
	}()

	p = &parser{
		mapping:   make(map[string]interface{}),
		types:     make(map[string]tomlType),
		lx:        lex(data),
		ordered:   make([]Key, 0),
		implicits: make(map[string]bool),
	}
	for {
		item := p.next()
		if item.typ == itemEOF {
			break
		}
		p.topLevel(item)
	}

	return p, nil
}

func (p *parser) panicf(format string, v ...interface{}) {
	msg := fmt.Sprintf("Near line %d (last key parsed '%s'): %s",
		p.approxLine, p.current(), fmt.Sprintf(format, v...))
	panic(parseError(msg))
}

func (p *parser) next() item {
	it := p.lx.nextItem()
	if it.typ == itemError {
		p.panicf("%s", it.val)
	}
	return it
}

func (p *parser) bug(format string, v ...interface{}) {
	panic(fmt.Sprintf("BUG: "+format+"\n\n", v...))
}

func (p *parser) expect(typ itemType) item {
	it := p.next()
	p.assertEqual(typ, it.typ)
	return it
}

func (p *parser) assertEqual(expected, got itemType) {
	if expected != got {
		p.bug("Expected '%s' but got '%s'.", expected, got)
	}
}

func (p *parser) topLevel(item item) {
	switch item.typ {
	case itemCommentStart:
		p.approxLine = item.line
		p.expect(itemText)
	case itemTableStart:
		kg := p.next()
		p.approxLine = kg.line

		var key Key
		for ; kg.typ != itemTableEnd && kg.typ != itemEOF; kg = p.next() {
			key = append(key, p.keyString(kg))
		}
		p.assertEqual(itemTableEnd, kg.typ)

		p.establishContext(key, false)
		p.setType("", tomlHash)
		p.ordered = append(p.ordered, key)
	case itemArrayTableStart:
		kg := p.next()
		p.approxLine = kg.line

		var key Key
		for ; kg.typ != itemArrayTableEnd && kg.typ != itemEOF; kg = p.next() {
			key = append(key, p.keyString(kg))
		}
		p.assertEqual(itemArrayTableEnd, kg.typ)

		p.establishContext(key, true)
		p.setType("", tomlArrayHash)
		p.ordered = append(p.ordered, key)
	case itemKeyStart:
		kname := p.next()
		p.approxLine = kname.line
		p.currentKey = p.keyString(kname)

		val, typ := p.value(p.next())
		p.setValue(p.currentKey, val)
		p.setType(p.currentKey, typ)
		p.ordered = append(p.ordered, p.context.add(p.currentKey))
		p.currentKey = ""
	default:
		p.bug("Unexpected type at top level: %s", item.typ)
	}
}

// Gets a string for a key (or part of a key in a table name).
func (p *parser) keyString(it item) string {
	switch it.typ {
	case itemText:
		return it.val
	case itemString, itemMultilineString,
		itemRawString, itemRawMultilineString:
		s, _ := p.value(it)
		return s.(string)
	default:
		p.bug("Unexpected key type: %s", it.typ)
		panic("unreachable")
	}
}

// value translates an expected value from the lexer into a Go value wrapped
// as an empty interface.
func (p *parser) value(it item) (interface{}, tomlType) {
	switch it.typ {
	case itemString:
		return p.replaceEscapes(it.val), p.typeOfPrimitive(it)
	case itemMultilineString:
		trimmed := stripFirstNewline(stripEscapedWhitespa