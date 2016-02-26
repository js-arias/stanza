// Copyright (c) 2016, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package stanza

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// A Writer writes records to a record-jar/stanza encoded file.
type Writer struct {
	Keys []string // list of keys
	w    *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

// Write writes a record to w.
func (w *Writer) Write(record Record) (err error) {
	ks := w.Keys
	if len(ks) == 0 {
		for key := range record {
			ks = append(ks, key)
		}
	}
	fields := false
	for _, k := range ks {
		v := strings.TrimSpace(record.Get(k))
		if len(v) == 0 {
			continue
		}
		err = w.writeKey(k)
		if err != nil {
			return err
		}
		err = w.writeValue(v)
		if err != nil {
			return err
		}
		fields = true
	}
	if fields {
		_, err = w.w.WriteString(`%%`)
		if err != nil {
			return
		}
		err = w.w.WriteByte('\n')
		if err != nil {
			return
		}
	}
	return nil
}

// Flush writes any bufferend data to the underlying io.Writer.
func (w *Writer) Flush() error {
	return w.w.Flush()
}

// writeKey writes a key.
func (w *Writer) writeKey(key string) (err error) {
	space := false
	for _, r1 := range key {
		if unicode.IsSpace(r1) || (r1 == ':') {
			space = true
			continue
		}
		if space {
			err = w.w.WriteByte('-')
			if err != nil {
				return
			}
		}
		_, err = w.w.WriteRune(r1)
		if err != nil {
			return
		}
	}
	_, err = w.w.WriteString(": ")
	return
}

// writeValue writes a value.
func (w *Writer) writeValue(value string) (err error) {
	space := false
	line := false
	for _, r1 := range value {
		if r1 == '\n' {
			line = true
			space = false
			continue
		}
		if unicode.IsSpace(r1) {
			space = true
			continue
		}
		if line {
			err = w.w.WriteByte('\n')
			if err != nil {
				return
			}
			err = w.w.WriteByte('\t')
			if err != nil {
				return
			}
			space = false
		}
		if space {
			err = w.w.WriteByte(' ')
			if err != nil {
				return
			}
		}
		_, err = w.w.WriteRune(r1)
	}
	return w.w.WriteByte('\n')
}
