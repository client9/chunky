package main

import (
	"fmt"
	//	"github.com/client9/chunky"
	"github.com/client9/chunky/tok"
	"io"
	"log"
	"os"
)

func main() {

	// Read from Stdin
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("unable read input: %s", err)
	}

	sentences := tok.Parse(string(b))
	format := "md"

	switch format {
	case "plain":
		for _, s := range sentences {
			for _, t := range s.Tokens {
				fmt.Printf("%s\n", t)
			}
		}
		fmt.Println("")

	case "md":
		fmt.Println("| sent | offset | word | tags | rule |")
		fmt.Println("|------|--------|------|------|------|")

		for sn, s := range sentences {
			for _, t := range s.Tokens {
				tags := t.String()[len(t.Word)+1:] // strip "word/" prefix
				if t.IsUnknownTag() {
					tags = "**UNK**"
				}
				fmt.Printf("| %d | %d | %s | %s | %s |\n", sn+1, t.Offset, t.Word, tags, t.Rule)
			}
		}
		fmt.Println("")
	}
}
