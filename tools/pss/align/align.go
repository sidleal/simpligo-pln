package main

import (
	"io"
	"log"
	"os"
	"strings"
)

type Sentence struct {
	Producao  string
	Level     string
	Text      string
	Sentence  string
	Paragraph string
	RawText   string
}

type Align struct {
	Producao  string
	Level     string
	TextA     string
	SentenceA string
	TextB     string
	SentenceB string
}

func main() {

	rawSentences := readFile("/home/sidleal/usp/coling2018/1_sentences.txt")
	lines := strings.Split(rawSentences, "\n")

	sentences := []Sentence{}
	for _, line := range lines {
		if line == "" {
			break
		}
		tokens := strings.Split(line, "\t")
		// log.Println(tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5])
		sentence := Sentence{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]}
		sentences = append(sentences, sentence)

	}

	rawAligns := readFile("/home/sidleal/usp/coling2018/1_aligns.txt")
	lines = strings.Split(rawAligns, "\n")

	aligns := []Align{}
	for _, line := range lines {
		if line == "" {
			break
		}
		tokens := strings.Split(line, "\t")
		// log.Println(tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5])
		align := Align{tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], tokens[5]}
		aligns = append(aligns, align)

	}

	// for i, item := range sentences {
	// 	log.Println(i, item.Level, item.Producao, item.Text, item.Sentence, item.RawText)
	// }

	// for i, item := range aligns {
	// 	log.Println(i, item.Level, item.Producao, item.TextA, item.SentenceA, item.TextB, item.SentenceB)
	// }

	count := 1
	for _, item := range sentences {
		if item.Level == "ORI" && item.Producao == "93" {
			log.Println("--------")
			log.Println(count, item.Level, item.Producao, item.Text, item.Sentence, item.RawText)
			pairs := getPairs(aligns, sentences, item.Sentence)
			for _, pair := range pairs {
				log.Println(count, pair.Level, pair.Producao, pair.Text, pair.Sentence, pair.RawText)
			}

			count++
		}

	}

}

func getPairs(aligns []Align, sentences []Sentence, id string) []Sentence {
	ret := []Sentence{}
	for _, item := range aligns {
		if item.SentenceA == id {
			if ok, sentence := getSentence(sentences, item.SentenceB); ok {
				ret = append(ret, sentence)
			}
		}
	}
	return ret

}

func getSentence(sentences []Sentence, id string) (bool, Sentence) {
	for _, item := range sentences {
		if item.Sentence == id {
			return true, item
		}
	}
	return false, Sentence{}

}

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	// r := charmap.ISO8859_1.NewDecoder().Reader(f)
	r := io.Reader(f)

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
