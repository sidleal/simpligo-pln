package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main_unico_pors_file() {

	log.Println("starting")

	f3, err := os.Create("/home/sidleal/sid/usp/lrev/porsimples.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f3.Close()

	_, err = f3.WriteString("text_id\tproduction\tlevel\ttext\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	path := "/home/sidleal/sid/usp/porsimples_text/porsimples_text_all"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	for i, f := range files {
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

		tokens := strings.Split(strings.TrimRight(fileName, ".txt"), "_")
		production := tokens[0][10:len(tokens[0])]
		level := "original"
		if len(tokens) > 1 {
			level = tokens[1]
		}

		log.Println(production, level)

		text := readFile(path+"/"+fileName, true)

		log.Println(text)

		text = strings.ReplaceAll(text, "\r", "")
		text = strings.ReplaceAll(text, "\t", " ")
		text = strings.ReplaceAll(text, "\n", "<br>")

		_, err = f3.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\n",
			i, production, level, text))
		if err != nil {
			log.Println("ERRO", err)
		}

	}

}

func main_xtrac_porsimples() {

	log.Println("starting")

	fout, err := os.OpenFile("/home/sidleal/sid/usp/lrev/porsimples_metrics_x.tsv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	data := readFile("/home/sidleal/sid/usp/lrev/porsimples.tsv", false)

	lines := strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}
		tokens := strings.Split(line, "\t")
		textID := tokens[0]
		production := tokens[1]
		level := tokens[2]
		rawText := tokens[3]

		if production != "158" && production != "160" {
			continue
		}

		text := strings.ReplaceAll(rawText, "<br>", "\n")
		text = cleanText(text)

		fRet := callMetrix(text)
		log.Println("-------------------------------------------------------")
		log.Println("-->", i)
		log.Println(fRet)

		feats := strings.Split(fRet, ",")

		header := "text_id\tproduction\tlevel\ttext\t"
		ret := fmt.Sprintf("%v\t%v\t%v\t%v\t", textID, production, level, rawText)
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
