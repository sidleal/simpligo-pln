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

func main() {

	concatSent := readFile("/home/sidleal/usp/coling2018/1_align_concat_ori.tsv")
	lines := strings.Split(concatSent, "\n")

	concatSentences := []Sentence{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := Sentence{}
		sentence.RawText = tokens[4]
		concatSentences = append(concatSentences, sentence)

	}

	sizeSent := readFile("/home/sidleal/usp/coling2018/1_align_size_ori.tsv")
	lines = strings.Split(sizeSent, "\n")

	sizeSentences := []Sentence{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		tokens := strings.Split(line, "\t")
		sentence := Sentence{}
		sentence.RawText = tokens[4]
		sizeSentences = append(sizeSentences, sentence)

	}

	log.Println(len(concatSentences), len(sizeSentences), len(sizeSentences)-len(concatSentences))

	for _, item := range sizeSentences {
		achou := false
		for _, itemC := range concatSentences {
			if itemC.RawText == item.RawText {
				achou = true
				break
			}
		}
		if !achou {
			log.Println("-----------------------", item.RawText)
		}
	}

	lastSent := ""
	count := 0
	for _, item := range sizeSentences {
		if lastSent == item.RawText {
			count++
			log.Println(item.RawText)
		}
		lastSent = item.RawText
	}

	log.Println("----> ", count)

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
