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

type Align struct {
	Production string
	Level      string
	From       AlignDetail
	To         AlignDetail
}
type AlignDetail struct {
	Idx string
	Raw string
}

func main() {

	aligns := []Align{}

	path := "/home/sidleal/usp/PROPOR2018/"

	for i := 1; i <= 165; i++ { //165
		production := fmt.Sprintf("production%v", i)
		prodAligns := processProdution(path, production)
		for _, align := range prodAligns {
			aligns = append(aligns, align)
		}
	}

	f, err := os.Create("/home/sidleal/usp/align9.txt")
	if err != nil {
		log.Println("ERRO", err)
	}

	defer f.Close()

	validations := []Align{}
	for i, align := range aligns {
		if align.From.Raw != align.To.Raw && (align.From.Raw+".") != align.To.Raw {
			n, err := f.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", align.Production, align.Level, align.From.Idx, align.From.Raw, align.To.Idx, align.To.Raw))
			if err != nil {
				log.Println("ERRO", err)
			}
			fmt.Printf("wrote %d bytes\n", n)

			if i%54 == 0 {
				validations = append(validations, align)
			}

		}
	}

	log.Println("--------- VALIDACAO ALINHAMENTOS -------")
	for i, align := range validations {
		log.Println("----------", i, "-----------")
		log.Println(align.From.Raw)
		log.Println(align.To.Raw)
	}

}

func processProdution(path string, production string) []Align {
	aligns := []Align{}

	original := Text{}
	original.Path = path + production + "/" + production + ".txt"

	natural := Text{}
	natural.Path = path + production + "/" + production + "_natural.txt"

	strong := Text{}
	strong.Path = path + production + "/" + production + "_strong.txt"

	if !fileExists(original.Path) {
		log.Println("ERRO------------------, não existe: ", original.Path)
		return []Align{}
	}

	original.Raw = readFile(original.Path)
	original.Raw = original.Raw[strings.Index(original.Raw, "\n")+1 : len(original.Raw)]
	log.Println(original.Raw)

	log.Print("\n\n-------------------------------------------------\n\n")

	natural.Raw = readFile(natural.Path)
	natural.Raw = natural.Raw[strings.Index(natural.Raw, "\n")+1 : len(natural.Raw)]
	log.Println(natural)

	log.Print("\n\n-------------------------------------------------\n\n")

	strong.Raw = readFile(strong.Path)
	strong.Raw = strong.Raw[strings.Index(strong.Raw, "\n")+1 : len(strong.Raw)]
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
		log.Println(sentMap)
		ret := orderByDesc(sentMap)
		// log.Println("-------------------->", ret[0].Key, ret[0].Value)
		log.Println("yoooooooooooooooooooooooooo")
		log.Println(rawSentFromMap[ret[0].Key])
		log.Println(sn.Raw)

		align := Align{}
		align.Level = "ORI->NAT"
		align.Production = production
		align.From = AlignDetail{ret[0].Key, rawSentFromMap[ret[0].Key]}
		align.To = AlignDetail{fmt.Sprintf("%d", sn.Idx), sn.Raw}
		aligns = append(aligns, align)
		log.Println("yoooooooooooooooooooooooooo")
		// log.Println("yyyyyyyyyyyy")

	}

	for _, sn := range strong.Sentences {
		// log.Println("xxxxxxxxxx")
		// log.Println(sn.Idx, sn.TF, sn.Raw)
		// log.Println("xxxxxxxxxx")

		sentMap := make(map[string]int64)
		rawSentFromMap := make(map[string]string)

		for _, s := range natural.Sentences {
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

		align := Align{}
		align.Level = "NAT->STR"
		align.Production = production
		align.From = AlignDetail{ret[0].Key, rawSentFromMap[ret[0].Key]}
		align.To = AlignDetail{fmt.Sprintf("%d", sn.Idx), sn.Raw}
		aligns = append(aligns, align)
		log.Println("oooooooooooooooooooooooooo")
		// log.Println("yyyyyyyyyyyy")

	}

	for _, align := range aligns {
		log.Println("\n\n-----------------------", align.Production, " - ", align.Level, "------------------------")
		log.Println(align.From.Idx, " -- ", align.From.Raw)
		log.Println(align.To.Idx, " -- ", align.To.Raw)
	}

	return aligns

}

func fillSentencesTF(text *Text) {
	var idx int64 = 0
	for _, p := range text.Parsed.Paragraphs {
		for _, s := range p.Sentences {
			if s.QtyWords > 2 {
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
