#!python3
import json
import nltk
from nltk.corpus import brown
nltk.download('brown')

twords = brown.tagged_words()
word_tags = {}
for word, tag in twords:
    word_tags.setdefault(word, [])
    if tag not in word_tags[word]:
        word_tags[word].append(tag)
print(json.dumps(word_tags, indent=2))

