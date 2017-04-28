# Stanza

Package stanza reads and write records in a list ('stanza') format.

## Format

Stanza files have the following format:

1. Each line containing a field must starts with the field name and
   separated from its content by ':' character. If the field name
   ends with a new line rather than ':', the field is considered as
   empty.
2. Field names are case insensitive (always read as lower caps),
   without spaces, and should be unique.
3. A field ends with a new line. If the content of the field extends
   more than one line, the next line should start with at least one
   space or tab character.
4. A record ends with a line that start with '%' character. Any
   character after '%' will be ignored (Usually "%%" is used to
   increase visibility of end-of-record).
5. Lines starting with '#' are taken as comments.
6. Empty lines are ignored.

## Example

An example of a stanza list is:

```
# Country data facts
name:	República Argentina
common:	Argentina
iso3166: AR
capital: Buenos Aires
population: 42669500
anthem:	Ya su trono dignísimo abrieron
	las Provincias Unidas del Sud
	y los libres del mundo responden:
	"¡Al gran pueblo argentino, salud!"
%%
name:	대한민국
common:	South Korea
iso3166: KR
capital: Seoul
population: 51302044
anthem:	무궁화 삼천리 화려강산
	대한 사람, 대한으로 길이 보전하세
%%
name:	中华人民共和国
common:	China
iso3166: CN
capital: Beijing
population: 1339724852
%%
name:	Росси́я
common:	Russia
iso3166: RU
capital: Moscow
population: 144192450
anthem: Славься, Отечество наше свободное,
	Братских народов союз вековой,
	Предками данная мудрость народная!
	Славься, страна! Мы гордимся тобой!
%%
```

## Source

Stanza file format are inspired by the record-jar/stanza format described
by E. Raymond
"[The Art of UNIX programming](http://www.catb.org/esr/writings/taoup/html/ch05s02.html#id2906931)"
(2003) Addison-Wesley , and C. Strozzi
[NoSQL list format](http://www.strozzi.it/cgi-bin/CSA/tw7/I/en_US/NoSQL/Table%20structure)
(2007).

## Authorship and license

Copyright (c) 2017, J. Salvador Arias <jsalarias@gmail.com>
All rights reserved.
Distributed under BSD-style license that can be found in the LICENSE file.

