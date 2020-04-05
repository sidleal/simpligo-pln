package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
)

// ------------------------------------

func maint() {

	path := "/home/sidleal/sid/usp/cloze_exps/"

	data := readFile(path + "50pars.txt")

	f, err := os.Create(fmt.Sprintf("%v/50pars_features.tsv", path))
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	lines := strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}
		line = strings.ReplaceAll(line, "\\", "")
		log.Println(i, line)

		sent := senter.ParseText(line)

		log.Println("Tot words:", sent.TotalWords)
		for _, s := range sent.Paragraphs[0].Sentences {
			log.Println(s.Idx, s.Text)

			fRet := callMetrix(s.Text)

			log.Println(fRet)
			fRet2 := strings.Split(fRet, "\n")
			feats := strings.Split(fRet2[2], ",")

			header := "id\tsent\ttext\t"
			ret := fmt.Sprintf("%v\t%v\t%v\t", i, s.Idx, s.Text)
			for _, feat := range feats {
				kv := strings.Split(feat, ":")
				if len(kv) > 1 {
					if i == 1 {
						header += kv[0] + "\t"
					}
					ret += kv[1] + "\t"
				}
			}

			if i == 1 && s.Idx == 1 {
				header = strings.TrimRight(header, "\t")
				_, err := f.WriteString(header + "\n")
				if err != nil {
					log.Println("ERRO", err)
				}
			}

			ret = strings.TrimRight(ret, "\t")
			_, err := f.WriteString(ret + "\n")
			if err != nil {
				log.Println("ERRO", err)
			}

		}

		// if i > 2 {
		// 	break
		// }
	}

}

func callMetrix(text string) string {

	timeout := time.Duration(300 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	url := "http://fw.nilc.icmc.usp.br:23380/api/v1/metrix/189/m3tr1x01?format=plain"

	resp, err := client.Post(url, "text", bytes.NewReader([]byte(text)))
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error parsing response. %v", err)
	}

	ret := string(body)

	return ret
}
