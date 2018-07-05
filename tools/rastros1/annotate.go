package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
)

func main() {

	regEx := regexp.MustCompile(`(.+) \t \[(.+)\] (<.*>)* ([A-Z0-9 ]+) @`)

	log.Println("carregando freq list")
	freq := map[string][]float64{}
	path := "/home/sidleal/usp/dispersion_br_v7.tsv"
	lines := readFileLines(path)
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}
		line := strings.Replace(line, ".", "", -1)
		line = strings.Replace(line, ",", ".", -1)

		cols := strings.Split(line, "\t")

		word := cols[0]
		dp, _ := strconv.ParseFloat(cols[12], 64)
		freqCor, _ := strconv.ParseFloat(cols[13], 64)

		freq[word] = []float64{dp, freqCor}
	}
	log.Println("freq em memoria")

	rawText := readFile("/home/sidleal/usp/rastros/CHC/Química que colore o céu.txt")

	sent := senter.ParseText(rawText)

	outputPalavras := ""
	output := ""
	output += fmt.Sprintf("Total Parágrafos: %v \n", sent.TotalParagraphs)
	output += fmt.Sprintf("Total Sentenças: %v \n", sent.TotalSentences)
	output += fmt.Sprintf("Total Tokens: %v \n", sent.TotalTokens)
	output += fmt.Sprintf("Total Palavras: %v \n", sent.TotalWords)

	output += "\n\n"

	for _, p := range sent.Paragraphs {
		for _, s := range p.Sentences {
			log.Println("processing", p.Idx, s.Idx)
			output += fmt.Sprintf("\n\nSentença %v (Parágrafo %v): \n\n", s.Idx, p.Idx)
			res := execPalavras(s.Text)
			outputPalavras += res + "\n"

			lines := strings.Split(res, "\n")

			output += s.Text + "\n\n"
			output += fmt.Sprintf("Palavras: %v - Tokens: %v.\n\n", s.QtyWords, s.QtyTokens)

			for _, line := range lines {
				if line != "" && line != "<ß>" && line != "</ß>" {
					matches := regEx.FindAllStringSubmatch(line, -1)
					for _, match := range matches {
						var dp float64
						var freqCor float64
						if wordFreq, found := freq[strings.ToLower(match[1])]; found {
							dp = wordFreq[0]
							freqCor = wordFreq[1]
						}

						output += fmt.Sprintf("%v - %v - %v - %v - %v \n", match[1], match[2], match[4], int(freqCor), dp)
					}

				}
			}

		}
	}

	output += "-----------------------------------------------------------------------------------------------------\n"
	output += "                                 PALAVRAS PARSER COMPLETO\n"
	output += "-----------------------------------------------------------------------------------------------------\n"

	output += outputPalavras

	f1, err := os.OpenFile("/home/sidleal/usp/rastros/teste.info", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f1.Close()

	_, err = f1.WriteString(output)
	if err != nil {
		log.Println("ERRO", err)
	}

	// fmt.Print(output)
	// fmt.Println("-------------------------------------------------")
	// fmt.Print(outputPalavras)

}

func getWordPalavras(word string, parseList []string) string {
	for _, line := range parseList {
		if strings.HasPrefix(line, word) {
			return line
		}
	}

	return ""
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

func execPalavras(content string) string {

	retType := "flat"
	options := "--dep-fuse"

	palavrasIP := ""
	palavrasPort := "23380"

	resp, err := http.PostForm("http://"+palavrasIP+":"+palavrasPort+"/"+retType,
		url.Values{"sentence": {content}, "options": {options}})
	if err != nil {
		log.Printf("Error: %v\n", err)
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("Error reading response: %v", err))
	}

	bodyString := string(body)

	return bodyString
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
