SHELL := sh

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## build module and CLI
	go build ./...

tok/rules_gen.go: rules.db rules.sh cmd/mkrules/main.go ## regenerate context disambiguation rules from corpus statistics
	{ \
	F1=NOUN F2=VERB          bash rules.sh | go run ./cmd/mkrules -tag1 NOUN -tag2 VERB  -var nounVerbRules; \
	F1=ADJ  F2=NOUN RATIO=20 bash rules.sh | go run ./cmd/mkrules -tag1 ADJ  -tag2 NOUN  -var adjNounRules  -noheader; \
	F1=ADP  F2=PART          bash rules.sh | go run ./cmd/mkrules -tag1 ADP  -tag2 PART  -var adpPartRules  -noheader; \
	F1=AUX  F2=VERB          bash rules.sh | go run ./cmd/mkrules -tag1 AUX  -tag2 VERB  -var auxVerbRules  -noheader; \
	F1=DET  F2=PRON          bash rules.sh | go run ./cmd/mkrules -tag1 DET  -tag2 PRON  -var detPronRules  -noheader; \
	F1=ADP  F2=SCONJ         bash rules.sh | go run ./cmd/mkrules -tag1 ADP  -tag2 SCONJ -var adpSconjRules -noheader; \
	F1=ADJ  F2=VERB          bash rules.sh | go run ./cmd/mkrules -tag1 ADJ  -tag2 VERB  -var adjVerbRules  -noheader; \
	F1=ADJ  F2=ADV           bash rules.sh | go run ./cmd/mkrules -tag1 ADJ  -tag2 ADV   -var adjAdvRules   -noheader; \
	F1=ADV  F2=DET           bash rules.sh | go run ./cmd/mkrules -tag1 ADV  -tag2 DET   -var advDetRules   -noheader; \
	F1=ADP  F2=ADV           bash rules.sh | go run ./cmd/mkrules -tag1 ADP  -tag2 ADV   -var adpAdvRules   -noheader; \
	F1=ADV  F2=NUM           bash rules.sh | go run ./cmd/mkrules -tag1 ADV  -tag2 NUM   -var advNumRules   -noheader; \
	} > tok/rules_gen.go

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
