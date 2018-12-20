package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"

	"github.com/sidleal/simpligo-pln/tools/senter"
)

type Sentence struct {
	FileName string
	Text     string
	QtyWords int64
}

type SentenceOrder []Sentence

func (a SentenceOrder) Len() int      { return len(a) }
func (a SentenceOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SentenceOrder) Less(i, j int) bool {
	return a[i].QtyWords > a[j].QtyWords //desc
}

func main() {

	sentList := []Sentence{}

	path := "/home/sidleal/usp/magali/CORPUS_SIMPLIFICAR2"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	fout, err := os.OpenFile("/home/sidleal/usp/magali/corpus2.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	_, err = fout.WriteString("File\tParagraphs\tSentences\tWords\tTokens\tBytes\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	listaSentAspas := [][]string{}
	for _, f := range files {
		log.Println(f.Name())
		data := readFile(path + "/" + f.Name())

		//tratar texto
		regEx := regexp.MustCompile(`<[a-z]+>.*>`)
		data = regEx.ReplaceAllString(data, "")

		regEx = regexp.MustCompile(`"([A-z0-9À-ú\(])`)
		data = regEx.ReplaceAllString(data, `“$1`)

		regEx = regexp.MustCompile(`([A-z0-9À-ú\.\)\?!,])"`)
		data = regEx.ReplaceAllString(data, `$1”`)

		regEx = regexp.MustCompile(`“[^”]+\.[^”]+”`)
		matches := regEx.FindAllStringSubmatch(data, -1)
		for _, match := range matches {
			listaSentAspas = append(listaSentAspas, []string{f.Name(), match[0]})
		}

		//parser senter
		descSenter := senter.ParseText(data)

		fmt.Printf("Paragraphs: %v - Sentences: %v - Words: %v - Tokens: %v - Bytes: %v\n", descSenter.TotalParagraphs, descSenter.TotalSentences, descSenter.TotalWords, descSenter.TotalTokens, f.Size())

		_, err = fout.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", f.Name(), descSenter.TotalParagraphs, descSenter.TotalSentences, descSenter.TotalWords, descSenter.TotalTokens, f.Size()))
		if err != nil {
			log.Println("ERRO", err)
		}

		for _, p := range descSenter.Paragraphs {
			for _, s := range p.Sentences {
				// log.Println(s.QtyWords)
				sentList = append(sentList, Sentence{f.Name(), s.Text, s.QtyWords})
			}
		}

	}

	sort.Sort(SentenceOrder(sentList))

	_, err = fout.WriteString("\n\nText\tWords\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	for _, sent := range sentList {
		//		log.Println(sent.QtyWords)
		_, err = fout.WriteString(fmt.Sprintf("%v\t%v\t%v\n", sent.FileName, sent.Text, sent.QtyWords))
		if err != nil {
			log.Println("ERRO", err)
		}

	}

	_, err = fout.WriteString("\n\nSeq\tFile\tSentence\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	for i, sent := range listaSentAspas {
		//log.Println(i, sent)
		_, err = fout.WriteString(fmt.Sprintf("%v\t%v\t%v\n", i, sent[0], sent[1]))
		if err != nil {
			log.Println("ERRO", err)
		}

	}

}

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	//r := charmap.ISO8859_1.NewDecoder().Reader(f)

	ret := ""

	buf := make([]byte, 32*1024)
	for {
		//n, err := r.Read(buf)
		n, err := f.Read(buf)
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

func readFileLines(path string) []string {
	ret := []string{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ret

}
