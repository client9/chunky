#!sh

#FILES=data/enwiki-20260401-multistream1.pos.jsonl
FILES=data/*.pos.jsonl

go run ./cmd/lexrules -feat prevtag -fmt matrix $FILES | \
	sqlite3 rules.db ".import --csv /dev/stdin features"

go run ./cmd/lexrules -feat nexttag -fmt matrix $FILES | sqlite3 rules.db ".import --csv /dev/stdin features"

go run ./cmd/lexrules -feat prevtag+nexttag -fmt matrix $FILES | sqlite3 rules.db ".import --csv /dev/stdin features"

go run ./cmd/lexrules -feat prev2tag+prevtag+nexttag -fmt matrix $FILES | sqlite3 rules.db ".import --csv /dev/stdin features"

go run ./cmd/lexrules -feat prev2tag+prevtag+nexttag+next2tag -fmt matrix $FILES | sqlite3 rules.db ".import --csv /dev/stdin features"
