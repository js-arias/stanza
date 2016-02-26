Stanza
======

Package stanza reads and write [record-jar/stanza files](http://catb.org/esr/writings/taoup/html/ch05s02.html#id2906931)
as described by E.S. Raymond, "The Art of Unix Programming" (2003).

A stanza file comprise multi-line records made up key/value pairs of the form:
	key: value
Each line contains an key/value pair. Long values may be folded over multiple
lines, provided that each continuation line starts with one or more spaces.

Record delimiter is a line consisting of "%%\n".

Blank lines, or lines cosisting only of whitespaces are ignored.

Lines starting with # are comments.

An example is:

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

Authorship and license
----------------------

Copyright (c) 2016, J. Salvador Arias <jsalarias@csnat.unt.edu.ar>
All rights reserved.
Distributed under BSD2 license that can be found in the LICENSE file.

