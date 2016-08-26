package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func words(s string) []string {
	strip := func(r rune) rune {
		if strings.ContainsRune(",.?!:;,", r) {
			return -1
		}
		return r
	}
	return strings.Fields(strings.Map(strip, s))
}

func groupByWords(subs []subTitle) map[string][]subTitle {
	out := make(map[string][]subTitle)
	for _, sub := range subs {
		for _, word := range words(sub.Text) {
			out[word] = append(out[word], sub)
		}
	}
	return out
}

func constructPhrase(s string, store map[string][]subTitle) ([]subTitle, error) {
	var out []subTitle
	for _, word := range words(s) {
		bucket, found := store[word]
		if !found {
			return nil, errors.New("can't satisfy")
		}
		out = append(out, bucket[0])
	}
	return out, nil
}

var (
	fileName = flag.String("fileName", "", "the .srt file to parse")
	phrase   = flag.String("phrase", "", "the phrase that you want to construct")
)

func main() {
	flag.Parse()
	if len(*fileName) == 0 || !strings.HasSuffix(*fileName, ".srt") {
		log.Fatal("You must specify a .srt file to parse.")
	}
	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	subs, err := parseSubs(file)
	if err != nil {
		log.Fatal(err)
	}
	if len(*phrase) == 0 {
		log.Fatal("You must specify a phrase to construct.")
	}
	fmt.Println(constructPhrase(*phrase, groupByWords(subs)))
}
