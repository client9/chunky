# Lexicon Generation

## Brown Corpus 

The current lexicon is generated from the Brown corpus, using python's NLTK.

The script to extract all the words, along with their tags, is in /scripts/brown-penn-tags.py

(Note: perhaps misnamed since the Brown Corpus uses Brown Tags, not Penn tags.).

This dumps out a large map (word to tags), encoded as JSON.

## Remap and Cleanup

In Go, these are remapped to UD tags using /cmd/brown-remap/main.go

In addition a number of words are deleted (anything with capitals, numbers, etc).

The output is tok/lexicon_gen.go

(Note: the output is not formated with 'gofmt'.. it should be and it's possible to do this in source code).

## Future Work

Items for future work.  All are not urgent and not critical.

The Go datastructure is a `map[string][]Tag`.  The value `[]Tag` while clear and understandable has number of problems:

- memory allocations
- order does not matter
- is actually a set -- tags are not duplicated.

It could be replaced with a bitfield.


Many words in the lexicon can be removed since they can reconsisted with the rules  in /tok/unknown.go.  For example, all nouns ending in "-tion" can be removed.


