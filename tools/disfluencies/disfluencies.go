package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/unidoc/unioffice/document"
)

var regEx = regexp.MustCompile(`([A-Z]+[0-9_]*[0-9_]+[0-9]+)(.*)`)

func main() {

	log.Println("starting")

	fout, err := os.OpenFile("/home/sidleal/sid/usp/datasetlucia/tudo.tsv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	path := "/home/sidleal/sid/usp/datasetlucia"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	for _, f := range files {
		fileName := f.Name()
		log.Println("---------------------------------------------------------------")
		log.Println(fileName)

		files2, err := ioutil.ReadDir(path + "/" + fileName)
		if err != nil {
			log.Println("Erro", err)
		}

		for _, f2 := range files2 {
			fileName2 := f2.Name()
			log.Println("--->", fileName2)

			files3, err := ioutil.ReadDir(path + "/" + fileName + "/" + fileName2)
			if err != nil {
				log.Println("Erro", err)
			}

			i := 0
			for _, f3 := range files3 {
				fileName3 := f3.Name()
				log.Println("--------->", fileName3)

				str, err := document.ExtractText(path+"/"+fileName+"/"+fileName2+"/"+fileName3, []int{1})
				// doc, err := document.Open(path + "/" + fileName + "/" + fileName2 + "/" + fileName3)
				if err != nil {
					log.Fatalf("error opening Windows Word 2016 document: %s", err)
				}
				log.Println(str)
				matchContent := regEx.FindStringSubmatch(str)

				if fileName == "Sem marcação de disfluências" {

					text := matchContent[2]

					fRet := callMetrix(text)
					log.Println(fRet)

					feats := strings.Split(fRet, ",")

					header := "tipo1,tipo2,arquivo,codigo,texto,"
					ret := fmt.Sprintf("%v\t%v\t%v\t%v\t%v,", fileName, fileName2, fileName3, matchContent[1], matchContent[2])
					for _, feat := range feats {
						kv := strings.Split(feat, ":")
						if len(kv) > 1 {
							if i == 0 {
								header += kv[0] + "\t"
							}
							ret += kv[1] + "\t"
						}
					}

					if i == 0 {
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

					i++

				} else {

					strtowrite := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", fileName, fileName2, fileName3, matchContent[1], matchContent[2])
					_, err = fout.WriteString(strtowrite)
					if err != nil {
						log.Println("ERRO", err)
					}
				}
			}

		}

	}

}

func callMetrix(text string) string {

	timeout := time.Duration(300 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Post("https://simpligo.sidle.al/api/v1/metrix/all/m3tr1x01", "text", bytes.NewReader([]byte(text)))
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
