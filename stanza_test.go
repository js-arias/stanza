// Copyright (c) 2016, J. Salvador Arias <jsalarias@gmail.com>
// All rights reserved.
// Distributed under BSD2 license that can be found in the LICENSE file.

package stanza

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

var arBlob = `
ISO3166-2: AR-C
Name:      Ciudad Autónoma de Buenos Aires
Category:  City
%%
ISO3166-2: AR-B
Name:      Buenos Aires
Category:  Province
Anthem:
	Llanura, sierra y camino
	hicieron de ti el destino
	del hombre, el pan y la paz,
	la Nación Argentina ennoblece
	tu alma con libertad.

%%
ISO3166-2: AR-K
Name:      Catamarca
Category:  Province
%%
ISO3166-2: AR-H
Name:      Chaco
Category:  Province
%%
ISO3166-2: AR-U
Name:      Chubut
Category:  Province
%%
ISO3166-2: AR-X
Name:      Córdoba
Category:  Province
%%
ISO3166-2: AR-W
Name:      Corrientes
Category:  Province
%%
ISO3166-2: AR-E
Name:      Entre Ríos
Category:  Province
%%
ISO3166-2: AR-P
Name:      Formosa
Category:  Province
%%
ISO3166-2: AR-Y
Name:      Jujuy
Category:  Province
%%
ISO3166-2: AR-L
Name:      La Pampa
Category:  Province
%%
ISO3166-2: AR-F
Name:      La Rioja
Category:  Province
%%
ISO3166-2: AR-M
Name:      Mendoza
Category:  Province
%%
ISO3166-2: AR-N
Name:      Misiones
Category:  Province
%%
ISO3166-2: AR-Q
Name:      Neuquén
Category:  Province
%%
ISO3166-2: AR-R
Name:      Río Negro
Category:  Province
%%
ISO3166-2: AR-A
Name:      Salta
Category:  Province
%%
ISO3166-2: AR-J
Name:      San Juan
Category:  Province
%%
ISO3166-2: AR-D
Name:      San Luis
Category:  Province
%%
ISO3166-2: AR-Z
Name:      Santa Cruz
Category:  Province
%%
ISO3166-2: AR-S
Name:      Santa Fe
Category:  Province
%%
ISO3166-2: AR-G
Name:      Santiago del Estero
Category:  Province
%%
ISO3166-2: AR-V
Name:      Tierra del Fuego
Category:  Province
%%
ISO3166-2: AR-T
Name:      Tucumán
Category:  Province
%%
`

func TestRead(t *testing.T) {
	r := NewReader(strings.NewReader(arBlob))
	i := 0
	for i = 0; ; i++ {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("Error while reading (rec: %d): %v", i, err)
		}
		if len(rec.Get("Name")) == 0 {
			t.Errorf("Empty record value (rec: %d)", i)
		}
		if rec.Get("Category") != rec.Get("category") {
			t.Errorf("Different values for the same key (rec: %d), %s, %s", i, rec.Get("Category"), rec.Get("category"))
		}
		if rec.Get("Name") == "Buenos Aires" {
			if len(rec.Get("anthem")) == 0 {
				t.Errorf("Expecting anthem value for 'Buenos Aires' province")
			}
			continue
		}
		if len(rec) != 3 {
			t.Errorf("Too many fields: %d (rec: %d-%s), expecting: %d", len(rec), i, rec.Get("Name"), 3)
		}
	}
	if i != 24 {
		t.Errorf("Wrong number of reads: %d, expecting %d", i, 24)
	}
	if len(r.Keys()) != 4 {
		t.Errorf("Wrong number of keys: %d, expecting %d", len(r.Keys()), 4)
	}
}

func testWrite(t *testing.T) {
	var b bytes.Buffer
	r := NewReader(strings.NewReader(arBlob))
	w := NewWriter(&b)
	for i := 0; ; i++ {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("Error while reading (rec: %d): %v", i, err)
		}
		err = w.Write(rec)
		if err != nil {
			t.Errorf("Error while writting (rec: %d): %v", i, err)
		}
	}
	err := w.Flush()
	if err != nil {
		t.Errorf("Error while flushing): %v", err)
	}

	// test that the result of the writting
	r = NewReader(strings.NewReader(b.String()))
	i := 0
	for i = 0; ; i++ {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("Error while reading write result (rec: %d): %v", i, err)
		}
		if len(rec.Get("Name")) == 0 {
			t.Errorf("Empty record value in write result (rec: %d)", i)
		}
		if rec.Get("Category") != rec.Get("category") {
			t.Errorf("Different values for the same key for a write result (rec: %d), %s, %s", i, rec.Get("Category"), rec.Get("category"))
		}
		if rec.Get("Name") == "Buenos Aires" {
			if len(rec.Get("anthem")) == 0 {
				t.Errorf("Expecting anthem value for 'Buenos Aires' province in write result")
			}
			continue
		}
		if len(rec) != 3 {
			t.Errorf("Too many fields in write result: %d (rec: %d-%s), expecting: %d", len(rec), i, rec.Get("Name"), 3)
		}
	}
	if i != 24 {
		t.Errorf("Wrong number of reads in write result: %d, expecting %d", i, 24)
	}
	if len(r.Keys()) != 4 {
		t.Errorf("Wrong number of keys in write result: %d, expecting %d", len(r.Keys()), 4)
	}
}
