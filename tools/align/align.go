package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
	"golang.org/x/text/encoding/charmap"
)

func main() {

	path := "/home/sidleal/usp/PROPOR2018/"

	for i := 1; i <= 1; i++ { //165
		production := fmt.Sprintf("production%v", i)
		processProdution(path, production)
	}
}

type Text struct {
	Path      string
	Raw       string
	Parsed    senter.ParsedText
	Sentences []Sentence
}
type Sentence struct {
	Idx int64
	TF  map[string]int64
	Raw string
}

func processProdution(path string, production string) {

	original := Text{}
	original.Path = path + production + "/" + production + ".txt"

	natural := Text{}
	natural.Path = path + production + "/" + production + "_natural.txt"

	strong := Text{}
	strong.Path = path + production + "/" + production + "_strong.txt"

	if !fileExists(original.Path) {
		log.Println("ERRO------------------, não existe: ", original.Path)
		return
	}

	original.Raw = readFile(original.Path)
	log.Println(original.Raw)

	log.Print("\n\n-------------------------------------------------\n\n")

	natural.Raw = readFile(natural.Path)
	log.Println(natural)

	log.Print("\n\n-------------------------------------------------\n\n")

	strong.Raw = readFile(strong.Path)
	log.Println(strong)

	original.Parsed = senter.ParseText(original.Raw)
	natural.Parsed = senter.ParseText(natural.Raw)
	strong.Parsed = senter.ParseText(strong.Raw)

	log.Println(original.Parsed)
	log.Println(natural.Parsed)
	log.Println(strong.Parsed)

	// words := make(map[string]int64)
	// words = countWords(original.Parsed, words)
	// words = countWords(natural.Parsed, words)
	// words = countWords(strong.Parsed, words)

	// orderByDesc(words)

	fillSentencesTF(&original)
	fillSentencesTF(&natural)
	fillSentencesTF(&strong)

	for _, sn := range natural.Sentences {
		// log.Println("xxxxxxxxxx")
		// log.Println(sn.Idx, sn.TF, sn.Raw)
		// log.Println("xxxxxxxxxx")

		sentMap := make(map[string]int64)
		rawSentFromMap := make(map[string]string)

		for _, s := range original.Sentences {
			var totWordsMatch int64 = 0
			for snk, snv := range sn.TF {
				totWordsMatch += s.TF[snk] * snv
			}
			sentMap[fmt.Sprintf("%d", s.Idx)] = totWordsMatch
			rawSentFromMap[fmt.Sprintf("%d", s.Idx)] = s.Raw
		}
		// log.Println("yyyyyyyyyyyy")
		// log.Println(sentMap)
		ret := orderByDesc(sentMap)
		// log.Println("-------------------->", ret[0].Key, ret[0].Value)
		log.Println("oooooooooooooooooooooooooo")
		log.Println(rawSentFromMap[ret[0].Key])
		log.Println(sn.Raw)
		log.Println("oooooooooooooooooooooooooo")
		// log.Println("yyyyyyyyyyyy")

	}

}

func fillSentencesTF(text *Text) {
	var idx int64 = 0
	for _, p := range text.Parsed.Paragraphs {
		for _, s := range p.Sentences {
			idx++
			words := make(map[string]int64)
			words = countSentenceWords(s, words)

			// log.Println("----------------")
			// orderByDesc(words)

			sentence := Sentence{idx, words, s.Text}
			text.Sentences = append(text.Sentences, sentence)
		}
	}

}

func countWords(text senter.ParsedText, words map[string]int64) map[string]int64 {
	for _, p := range text.Paragraphs {
		for _, s := range p.Sentences {
			words = countSentenceWords(s, words)
		}
	}
	return words
}

func countSentenceWords(s senter.ParsedSentence, words map[string]int64) map[string]int64 {
	for _, t := range s.Tokens {
		if t.IsWord == 1 {
			token := strings.ToLower(t.Token)
			words[token]++
		}
	}
	return words
}

type KeyValue struct {
	Key   string
	Value int64
}

func orderByDesc(words map[string]int64) []KeyValue {
	var ss []KeyValue
	for k, v := range words {
		ss = append(ss, KeyValue{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// for _, kv := range ss {
	// 	fmt.Printf("%s, %d\n", kv.Key, kv.Value)
	// }

	return ss
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	r := charmap.ISO8859_1.NewDecoder().Reader(f)

	ret := ""

	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			ret += string(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
	}

	return ret

}
