// Copyright (c) 2017, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package stanza

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

var blob = `
# Country data facts
Name:	República Argentina
Common:	Argentina
ISO3166: AR
Capital: Buenos Aires
Population: 42669500
Anthem:	Ya su trono dignísimo abrieron
	las Provincias Unidas del Sud
	y los libres del mundo responden:
	"¡Al gran pueblo argentino, salud!"
%
Name:	대한민국
Common:	South Korea
ISO3166: KR
Capital: Seoul
Population: 51302044
Anthem:	무궁화 삼천리 화려강산
	대한 사람, 대한으로 길이 보전하세
%%
Name:	中华人民共和国
Common:	China
ISO3166: CN
Capital: Beijing
Population: 1339724852
%
Name:	Росси́я
Common:	Russia
ISO3166: RU
Capital: Moscow
Population: 144192450
Anthem: Славься, Отечество наше свободное,
	Братских народов союз вековой,
	Предками данная мудрость народная!
	Славься, страна! Мы гордимся тобой!
%%
`

func TestRead(t *testing.T) {
	r := NewReader(strings.NewReader(blob))
	i := 0
	for {
		rec, err := r.Read()
		if err != nil {
			if errors.Cause(err) == io.EOF {
				break
			}
			t.Errorf("read: reading error: %v", err)
		}
		if _, ok := rec["common"]; !ok {
			t.Errorf("read: field %q not found", "common")
		}
		if _, ok := rec["iso3166"]; !ok {
			t.Errorf("read: field %q not found", "iso3166")
		}
		if rec["common"] != "China" {
			if len(rec) != 6 {
				t.Errorf("read: expecting 6 fields, found: %d", len(rec))
			}
			an, ok := rec["anthem"]
			if !ok {
				t.Errorf("read: field %q not found", "anthem")
			}
			if v := len(strings.Split(an, "\n")); v < 2 {
				t.Errorf("read: field %q should be multiline", "anthem")
			}
		} else if len(rec) != 5 {
			t.Errorf("read: expecting 5 fields, found: %d", len(rec))
		}
		i++
	}
	if i != 4 {
		t.Errorf("read: expecting 4 records, found: %d", i)
	}
}

func TestWrite(t *testing.T) {
	r := NewReader(strings.NewReader(blob))
	country := make(map[string]map[string]string)
	out := &bytes.Buffer{}
	w := NewWriter(out)
	w.SetFields([]string{"name", "common", "iso3166", "capital", "population", "anthem"})
	for {
		rec, err := r.Read()
		if err != nil {
			if errors.Cause(err) == io.EOF {
				break
			}
			t.Errorf("write: reading error: %v", err)
		}
		err = w.Write(rec)
		if err != nil {
			t.Errorf("write: %v", err)
		}
		country[rec["common"]] = rec
	}
	if err := w.Flush(); err != nil {
		t.Errorf("write: flushing error: %v", err)
	}

	r = NewReader(strings.NewReader(out.String()))
	i := 0
	for {
		rec, err := r.Read()
		if err != nil {
			if errors.Cause(err) == io.EOF {
				break
			}
			t.Errorf("write: reading error: %v", err)
		}
		p, ok := country[rec["common"]]
		if !ok {
			t.Errorf("write: country %q not found", rec["common"])
			continue
		}
		for f, v := range rec {
			if p[f] != v {
				t.Errorf("write: country %q: expecting %q, found %q", rec["common"], p["common"], v, p[f])
			}
		}
		i++
	}
	if i != 4 {
		t.Errorf("write: expecting 4 records, found: %d", i)
	}
}
