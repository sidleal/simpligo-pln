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
)

func main() {

	log.Println("starting")

	// f, err := os.OpenFile("/home/sidleal/sid/usp/arthur/out.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f, err := os.OpenFile("out-esic.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	raw := readFile("esic.csv")
	lines := strings.Split(raw, "\n")

	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		tokens := strings.Split(line, ";")

		textID := tokens[0]
		text := tokens[1]

		text = strings.ReplaceAll(text, "\\r\\n", "{{enter}}")
		text = strings.ReplaceAll(text, "\\n", "{{enter}}")
		text = strings.ReplaceAll(text, "\\\"", "\"")
		text = strings.ReplaceAll(text, "è", "e")
		text = strings.ReplaceAll(text, "ì", "i")
		text = strings.ReplaceAll(text, "ò", "o")
		text = strings.ReplaceAll(text, "ù", "u")

		text = strings.ReplaceAll(text, " à ", "{{crase}}")
		text = strings.ReplaceAll(text, "à", "a")
		text = strings.ReplaceAll(text, "{{crase}}", " à ")
		text = strings.ReplaceAll(text, "`", "\"")
		text = strings.ReplaceAll(text, "´", "\"")

		log.Println(i, textID, text)
		log.Println("---------------------------")

		fRet := callMetrixNilc(text)
		log.Println(fRet)

		fTokens := strings.Split(fRet, "-------------------------------\n")
		if len(fTokens) < 2 {
			continue
		}

		feats := strings.Split(fTokens[1], ",")

		header := "id_texto,"
		ret := textID + ","
		for _, feat := range feats {
			kv := strings.Split(feat, ":")
			if len(kv) > 1 {
				if i == 1 {
					header += kv[0] + ","
				}
				ret += kv[1] + ","
			}
		}

		if i == 1 {
			header = strings.TrimRight(header, ",")
			_, err := f.WriteString(header + "\n")
			if err != nil {
				log.Println("ERRO", err)
			}
		}

		ret = strings.TrimRight(ret, ",")
		// ret += retRanker
		_, err := f.WriteString(ret + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}
	}

	// log.Println(header, ret)

}

func callMetrixNilc(text string) string {

	timeout := time.Duration(300 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	url := "http://fw.nilc.icmc.usp.br:23380/api/v1/metrix/all/m3tr1x01?format=plain"

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
