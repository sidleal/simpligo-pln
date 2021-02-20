package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var regExTitle = regexp.MustCompile(`<title>(.*?)<\/title>`)
var regExSubtitle = regexp.MustCompile(`<subtitle>(.*?|.*?\n)<\/subtitle>`)
var regExElse = regexp.MustCompile(`<.*?>.*?<\/.*?>`)
var regExTag = regexp.MustCompile(`<(.*?)>`)

func main_adapt_unico() {

	log.Println("starting")

	f3, err := os.Create("/home/sidleal/sid/usp/lrev/adapt2kids.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f3.Close()

	_, err = f3.WriteString("text_id\tgroup\tfile\ttext\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	path := "/home/sidleal/sid/usp/lrev/Adapt2Kids-Corpus"

	files := listAllSubTxts(path)
	for i, f := range files {
		log.Println(f)

		filePath := f[45:len(f)]
		fileTokens := strings.Split(filePath, "/")
		fileName := fileTokens[len(fileTokens)-1]
		fileGroup := filePath[0 : len(filePath)-len(fileName)-1]

		text := readFile(f, false)

		log.Println(text)

		text = regExTitle.ReplaceAllString(text, `$1`)
		text = regExSubtitle.ReplaceAllString(text, `$1`)
		text = regExElse.ReplaceAllString(text, "")
		text = strings.ReplaceAll(text, "<Não tem>", "")
		text = strings.ReplaceAll(text, "<não tem>", "")
		text = regExTag.ReplaceAllString(text, `$1`)
		text = strings.ReplaceAll(text, "\r", "")
		text = strings.ReplaceAll(text, "\n\n\n", "\n")
		text = strings.ReplaceAll(text, "\n\n\n", "\n")
		text = strings.ReplaceAll(text, "\t", " ")
		text = strings.ReplaceAll(text, "_", " ")
		text = strings.ReplaceAll(text, "\n", "<br>")

		_, err = f3.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\n", i+1, fileGroup, fileName, text))
		if err != nil {
			log.Println("ERRO", err)
		}

	}

}

func listAllSubTxts(path string) []string {
	ret := []string{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	for _, f := range files {
		fileName := f.Name()
		log.Println(fileName)
		if f.IsDir() {
			recRet := listAllSubTxts(path + "/" + f.Name())
			ret = append(ret, recRet...)
		}

		if strings.HasSuffix(fileName, "txt") {
			ret = append(ret, path+"/"+fileName)
		}
	}

	return ret
}

func main_adapt_xtract() { //

	log.Println("starting")

	fout, err := os.OpenFile("/home/sidleal/sid/usp/lrev/adapt2kids_metrics_5.4k.tsv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	data := readFile("/home/sidleal/sid/usp/lrev/adapt2kids.tsv", false)

	lines := strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}

		if i < 5400 || i > 5500 {
			continue
		}

		tokens := strings.Split(line, "\t")
		textID := tokens[0]
		group := tokens[1]
		fileName := tokens[2]
		rawText := tokens[3]

		text := strings.ReplaceAll(rawText, "<br>", "\n")
		text = cleanText(text)

		fRet := callMetrix(text)
		log.Println("-------------------------------------------------------")
		log.Println("-->", i)
		log.Println(fRet)

		feats := strings.Split(fRet, ",")

		header := "text_id\tgroup\tfile\ttext\t"
		ret := fmt.Sprintf("%v\t%v\t%v\t%v\t", textID, group, fileName, rawText)
		for _, feat := range feats {
			kv := strings.Split(feat, ":")
			if len(kv) > 1 {
				if i == 1 {
					header += kv[0] + "\t"
				}
				ret += kv[1] + "\t"
			}
		}

		if i == 1 {
			header = strings.TrimRight(header, "\t")
			_, err := fout.WriteString(header + "\n")
			if err != nil {
				log.Println("ERRO", err)
			}
		}

		ret = strings.TrimRight(ret, "\t")
		_, err := fout.WriteString(ret + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}

}
