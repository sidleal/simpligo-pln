package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var regEx *regexp.Regexp = regexp.MustCompile(`(?ms)Local:(.*?)\nData da coleta:(.*?)\nResponsável pela coleta:(.*?)\nProposta:(.*?)\n.*?Código:(.*?)\nSérie:(.*?)\nEdição:(.*?)\n(.*)`)
var regEx2 = regexp.MustCompile(`\\(\w+)\/`)

func main_textus_arquivo_unico() {

	log.Println("starting")

	f3, err := os.Create("/home/sidleal/sid/usp/lrev/textus.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f3.Close()

	_, err = f3.WriteString("texto_id\tc1\tc2\tlocal\tdata\tresponsavel\tproposta\tcodigo\tserie\tedicao\ttexto\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	path := "/home/sidleal/sid/usp/textus/out"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	for _, f := range files {
		fileName := f.Name()
		log.Println("---------------------------------------------------------------")
		log.Println(fileName)

		if !strings.HasSuffix(fileName, "txt") {
			continue
		}

		// if i > 10 {
		// 	break
		// }

		log.Println(path + "/" + fileName)

		tokens := strings.Split(strings.TrimRight(fileName, ".txt"), "-")
		textID := tokens[0]
		tokens2 := strings.Split(tokens[1], "_")
		c1 := tokens2[0]
		c3 := tokens2[1] // serie
		c2 := tokens2[2]
		c4 := tokens2[3] //proposta

		log.Println(textID, c1, c2, c3, c4)

		text := readFile(path+"/"+fileName, false)

		log.Println(text)

		match := regEx.FindStringSubmatch(text)
		local := strings.TrimSpace(match[1])
		data := strings.TrimSpace(match[2])
		resp := strings.TrimSpace(match[3])
		prop := strings.TrimSpace(match[4])
		cod := strings.TrimSpace(match[5])
		serie := strings.TrimSpace(match[6])
		edicao := strings.TrimSpace(match[7])
		rawText := strings.TrimSpace(match[8])

		rawText = strings.ReplaceAll(rawText, "\n\n", "|par|")
		rawText = strings.ReplaceAll(rawText, "_\n", "")
		rawText = strings.ReplaceAll(rawText, "-\n", "")
		rawText = strings.ReplaceAll(rawText, "\n", " ")
		rawText = strings.ReplaceAll(rawText, "|par|", "<br>")
		rawText = strings.ReplaceAll(rawText, "\\(rasura)*/", "")
		rawText = strings.ReplaceAll(rawText, "(rasura)*", "")
		rawText = strings.ReplaceAll(rawText, "●", "- ")
		rawText = strings.ReplaceAll(rawText, "\t", " ")

		rawText = regEx2.ReplaceAllString(rawText, `$1`)

		log.Println(local, data, resp, prop, cod, serie, edicao)
		log.Println(rawText)

		_, err = f3.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			textID, c1, c2, local, data, resp, prop, cod, serie, edicao, rawText))
		if err != nil {
			log.Println("ERRO", err)
		}

	}

}

func main_textus_xtract() { //

	log.Println("starting")

	fout, err := os.OpenFile("/home/sidleal/sid/usp/lrev/textus_metrics_rest3.tsv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	mapProcessed := map[string]int{}
	data := readFile("/home/sidleal/sid/usp/lrev/textus_metrics.tsv", false)

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		tokens := strings.Split(line, "\t")
		textID := tokens[0]
		mapProcessed[textID] = 1
	}

	data = readFile("/home/sidleal/sid/usp/lrev/textus.tsv", false)

	lines = strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}

		tokens := strings.Split(line, "\t")
		textID := tokens[0]
		c1 := tokens[1]
		c2 := tokens[2]
		local := tokens[3]
		data := tokens[4]
		resp := tokens[5]
		prop := tokens[6]
		cod := tokens[7]
		serie := tokens[8]
		edicao := tokens[9]
		rawText := tokens[10]

		if _, found := mapProcessed[textID]; found {
			continue
		}

		text := strings.ReplaceAll(rawText, "<br>", "\n")
		text = cleanText(text)

		fRet := callMetrix(text)
		log.Println("-------------------------------------------------------")
		log.Println("-->", i)
		log.Println(fRet)

		feats := strings.Split(fRet, ",")

		header := "texto_id\tc1\tc2\tlocal\tdata\tresponsavel\tproposta\tcodigo\tserie\tedicao\ttexto\t"
		ret := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t", textID, c1, c2, local, data, resp, prop, cod, serie, edicao, rawText)
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
