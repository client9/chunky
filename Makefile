SHELL := sh

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## build module and CLI
	go build ./...
bench:
	go test ./tok/ -bench=BenchmarkParse -benchmem

# aslso -errors NP/VP/etc in second command
accuracy:
	go run ./cmd/eval-pos data/conll2000-test.txt
	go run ./cmd/eval-chunks data/conll2000-test.txt

tok/rules_gen.go: rules.db rules.sh cmd/mkrules/main.go ## regenerate context disambiguation rules from corpus statistics
	{ \
	F1=NOUN  F2=VERB  MIN=4000 RATIO=10 bash rules.sh | sed 's/^/NOUN|VERB|/'; \
	F1=ADJ   F2=NOUN  MIN=4000 RATIO=20 bash rules.sh | sed 's/^/ADJ|NOUN|/'; \
	F1=ADP   F2=PART  MIN=4000 RATIO=10 bash rules.sh | sed 's/^/ADP|PART|/'; \
	F1=AUX   F2=VERB  MIN=4000 RATIO=10 bash rules.sh | sed 's/^/AUX|VERB|/'; \
	F1=DET   F2=PRON  MIN=4000 RATIO=10 bash rules.sh | sed 's/^/DET|PRON|/'; \
	F1=ADP   F2=SCONJ MIN=4000 RATIO=10 bash rules.sh | sed 's/^/ADP|SCONJ|/'; \
	F1=ADJ   F2=VERB  MIN=4000 RATIO=10 bash rules.sh | sed 's/^/ADJ|VERB|/'; \
	F1=ADJ   F2=ADV   MIN=4000 RATIO=10 bash rules.sh | sed 's/^/ADJ|ADV|/'; \
	F1=ADV   F2=DET   MIN=4000 RATIO=10 bash rules.sh | sed 's/^/ADV|DET|/'; \
	F1=ADP   F2=ADV   MIN=4000 RATIO=10 bash rules.sh | sed 's/^/ADP|ADV|/'; \
	F1=ADV   F2=NUM   MIN=4000 RATIO=10 bash rules.sh | sed 's/^/ADV|NUM|/'; \
	} | go run ./cmd/mkrules -allpairs > tok/rules_gen.go

tok/lexicon_gen.go: data/brown-penn-nltk.json cmd/brown-remap/main.go closed.go words.go ## generate compiled-in lexicon for tok package
	go run ./cmd/brown-remap/ -go < data/brown-penn-nltk.json > tok/lexicon_gen.go 2>/dev/null

data/brown-reduced.json: data/brown-penn-nltk.json cmd/brown-remap/main.go closed.go words.go ## generate brown corpus with UD tags
	go run ./cmd/brown-remap/ < data/brown-penn-nltk.json > data/brown-reduced.json

data/brown-penn-nltk.json: scripts/brown-penn-tags.py ## extract Brown lexicon with Brown Tags
	python3 ./scripts/brown-penn-tags.py > data/brown-penn-nltk.json

closed.go: cmd/closedforms/main.go ## generate closed form list
	go run ./cmd/closedforms/main.go -go data/enwiki-20260401-multistream1.pos.txt > closed.go

clean: ## clean
	rm -f data/brown-reduced.json

eval: ## score chunker against CoNLL-2000 test set
	go run ./cmd/eval-chunks data/conll2000-test.txt

test: ## tests
	go test ./...

cover: ## generate code coverage report
	rm -f cover.out
	go test -run='^Test' -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -func=cover.out

## NOTE: this downloads it's schema over the network
lintverify:
	golangci-lint config verify

fmt: ## reformat source code
	go mod tidy
	gofmt -w -s .

lint: ## lint and verify repo is already formatted
	go mod tidy
	git diff --exit-code -- go.mod go.sum
	test -z "$$(gofmt -l *.go)"
	golangci-lint run .
