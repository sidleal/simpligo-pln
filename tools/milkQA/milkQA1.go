package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var regEx = regexp.MustCompile(`<assunto>.*<corpo>(.*)`)

func mainx() {

	log.Println("starting")

	mapWords := map[string]int{}
	qtTokens := 0
	maxTokens := 0
	minTokens := 1000

	path := "/home/sidleal/sid/usp/TopicosPLN/PerguntasMilkQA"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	fout, err := os.OpenFile(path+"/all_questions.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	for _, f := range files {
		fileName := f.Name()
		if strings.HasSuffix(fileName, "txt") && fileName != "all_questions.txt" {
			log.Println("---------------------------------------------------------------")
			log.Println(fileName)
			tokens := strings.Split(fileName, "-")
			fileID := tokens[0]
			raw := readFile(path + "/" + fileName)
			log.Println(raw)
			textFull := strings.ReplaceAll(raw, "\n", " ")
			textFull = strings.ReplaceAll(textFull, "\r", "")
			matchContent := regEx.FindStringSubmatch(textFull)

			text := matchContent[1]
			log.Println(fileID, text)
			_, err := fout.WriteString(fileID + "\t" + text + "\t\t\t\n")
			if err != nil {
				log.Println("ERRO", err)
			}

			text = strings.ToLower(text)
			ttok := strings.Split(text, " ")
			if maxTokens < len(ttok) {
				maxTokens = len(ttok)
			}
			if minTokens > len(ttok) {
				minTokens = len(ttok)
			}
			for _, t := range ttok {
				qtTokens++
				if _, found := mapWords[t]; found {
					mapWords[t]++
				} else {
					mapWords[t] = 1
				}
			}

		}
	}

	log.Println("--------------------------------------------")
	log.Println("Tot tokens", qtTokens)
	log.Println("Max tokens", maxTokens)
	log.Println("Min tokens", minTokens)
	log.Println("Words", len(mapWords))
	log.Println("--------------------------------------------")

	// raw := readFile("/home/sidleal/sid/usp/rastros/maisum/rastros150.txt")
	// lines := strings.Split(raw, "\n")

	// for i, line := range lines {
	// 	if line == "" {
	// 		break
	// 	}
	// 	log.Println(line)

	// 	// retRanker := callRanker(line)

	// 	fRet := callMetrix(line)

	// 	feats := strings.Split(fRet, ",")

	// 	header := ""
	// 	ret := ""
	// 	for _, feat := range feats {
	// 		kv := strings.Split(feat, ":")
	// 		if len(kv) > 1 {
	// 			if i == 0 {
	// 				header += kv[0] + ","
	// 			}
	// 			ret += kv[1] + ","
	// 		}
	// 	}

	// 	if i == 0 {
	// 		header = strings.TrimRight(header, ",")
	// 		// header += "complexity"
	// 		_, err := f.WriteString(header + "\n")
	// 		if err != nil {
	// 			log.Println("ERRO", err)
	// 		}
	// 	}

	// 	ret = strings.TrimRight(ret, ",")
	// 	// ret += retRanker
	// 	_, err := f.WriteString(ret + "\n")
	// 	if err != nil {
	// 		log.Println("ERRO", err)
	// 	}

	// }

	// log.Println(header, ret)

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
