// Copyright (c) 2016, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

// Package stanza reads and writes record-jar/stanza files as described by
// E.S. Raymond, "The Art of Unix Programming", Chapter 5.
//
// A stanza file comprise multi-line records made up key/value pairs of the
// form:
//	key: value
// Each line contains an key/value pair. Long values may be folded over
// multiple lines, provided that each continuation line starts with one or
// more spaces.
//
// Record delimiter is a line consisting of "%%\n".
//
// Blank lines, or lines cosisting only of whitespaces are ignored.
//
// Lines starting with # are comments.
//
// An example is:
//	ISO3166-2: AR-C
//	Name:      Ciudad Autónoma de Buenos Aires
//	Category:  City
//	%%
//	ISO3166-2: AR-B
//	Name:      Buenos Aires
//	Category:  Province
//	%%
//	ISO3166-2: AR-K
//	Name:      Catamarca
//	Category:  Province
//	%%
//	ISO3166-2: AR-H
//	Name:      Chaco
//	Category:  Province
//	%%
package stanza

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// A Record is a map of string keys to string values. The keys are
// case-insensitive.
type Record map[string]string

// Set sets the key to value. It replaces any existing value.
func (r Record) Set(key, value string) {
	r[strings.ToLower(key)] = value
}

// Get gets the value associated with the given key. If there is no value
// associated with the key, Get returns the empty string.
func (r Record) Get(key string) string {
	if r == nil {
		return ""
	}
	v, ok := r[strings.ToLower(key)]
	if !ok {
		return ""
	}
	return v
}

// A Reader reads records from a record-jar/stanza encoded file.
type Reader struct {
	keys  []string
	line  int
	r     *bufio.Reader
	field bytes.Buffer
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

// Read reads one record from r.
func (r *Reader) Read() (record Record, err error) {
	for {
		record, err = r.parseRecord()
		if len(record) != 0 {
			break
		}
		if err != nil {
			if err != io.EOF {
				err = fmt.Errorf("line %d: %v", err)
			}
			return nil, err
		}
	}
	return record, nil
}

// Keys returns the list of all keys read (in order of appareance while
// reading). It is useful when we need to know the keys in order because
// a map keys are returned in random order.
func (r *Reader) Keys() []string {
	return r.keys
}

// skip reads runes up to and including the rune delim or until error.
func (r *Reader) skip(delim rune) error {
	for {
		r1, _, err := r.r.ReadRune()
		if err != nil {
			return err
		}
		if r1 == delim {
			return nil
		}
	}
}

// parseRecord reads and parses a single record-jar/stanza record from r.
func (r *Reader) parseRecord() (record Record, err error) {
	record = Record{}
	// check out for empty records, empty lines, or comment lines
	for {
		r.line++
		r1, _, err := r.r.ReadRune()
		if err != nil {
			return nil, err
		}
		if r1 == '%' {
			err = r.skip('\n')
			if err != nil {
				return nil, err
			}
			break
		}
		if r1 == '\n' {
			continue
		}
		if r1 == '#' {
			err = r.skip('\n')
			if err != nil {
				return nil, err
			}
			continue
		}
		r.r.UnreadRune()
		r.line--
		break
	}
	for {
		r.line++
		key, err := r.parseKey()
		if err != nil {
			return nil, err
		}
		value, next, err := r.parseValue()
		if err != nil {
			return nil, err
		}
		if (len(value) != 0) && (len(key) != 0) {
			nk := true
			kv := strings.ToLower(key)
			for _, k := range r.keys {
				if kv == k {
					nk = false
					break
				}
			}
			record.Set(key, value)
			if nk {
				r.keys = append(r.keys, kv)
			}
		}
		if next == 0 {
			break
		}
		if next == '%' {
			r.line++
			err = r.skip('\n')
			if err != nil {
				if len(record) == 0 {
					return nil, err
				}
				break
			}
			break
		}
	}
	return record, nil
}

// parseKey parses the key string.
func (r *Reader) parseKey() (key string, err error) {
	r.field.Reset()
	for {
		r1, _, err := r.r.ReadRune()
		if err != nil {
			return "", err
		}
		if unicode.IsSpace(r1) {
			return "", errors.New("unexpected space in key field")
		}
		if r1 == ':' {
			return r.field.String(), nil
		}
		r.field.WriteRune(r1)
	}
}

// parseValue parses the value string. Next rune is the next valid rune.
func (r *Reader) parseValue() (value string, next rune, err error) {
	r.field.Reset()
	space := false
	first := true
	nline := false
	for {
		r1, _, err := r.r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", 0, err
		}
		if r1 == '\n' {
			r1, _, err = r.r.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				return "", 0, err
			}
			if r1 == '\n' {
				r.line++
				// unreads the ending '\n' so it can read the
				// next line
				r.r.UnreadRune()
				continue
			}
			if r1 == '#' {
				err = r.skip('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					return "", 0, err
				}
				// unreads the ending '\n' so it can read the
				// next line
				r.r.UnreadRune()
				continue
			}
			if r1 == '%' {
				next = r1
				break
			}
			if unicode.IsSpace(r1) {
				space = false
				nline = true
				continue
			}
			next = r1
			r.r.UnreadRune()
			break
		}
		if unicode.IsSpace(r1) {
			if (!nline) && (!first) {
				space = true
			}
			continue
		}
		if nline {
			if !first {
				r.field.WriteRune('\n')
			}
			nline = false
		}
		if space {
			r.field.WriteRune(' ')
			space = false
		}
		r.field.WriteRune(r1)
		first = false
	}
	return r.field.String(), next, nil
}
