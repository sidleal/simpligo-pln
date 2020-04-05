package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ------------------------------------

func mainy() {

	path := "/home/sidleal/sid/usp/cloze_exps/"

	data := readFile(path + "50pars.txt")

	lines := strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}
		log.Println(i, line)

		f, err := os.Create(fmt.Sprintf("%v/palavras/par%02d.txt", path, i))
		if err != nil {
			log.Println("ERRO", err)
		}

		flat := callPalavras(line, "flat")
		flat = strings.ReplaceAll(flat, "</ß>\n\n</ß>", "</ß>")
		log.Println(flat)

		_, err = f.WriteString(flat)
		if err != nil {
			log.Println("ERRO", err)
		}

		f.Close()

		f2, err := os.Create(fmt.Sprintf("%v/palavras/par%02d.xml", path, i))
		if err != nil {
			log.Println("ERRO", err)
		}

		tigerxml := callPalavras(line, "tigerxml")
		tigerxml = strings.ReplaceAll(tigerxml, "</ß>\n</ß>", "</ß>")
		log.Println(tigerxml)

		_, err = f2.WriteString(tigerxml)
		if err != nil {
			log.Println("ERRO", err)
		}

		f2.Close()

	}

}

func callPalavras(text string, retType string) string {

	options := "--dep-fuse"

	resp, err := http.PostForm("http://fw.nilc.icmc.usp.br:23380/api/v1/palavras/"+retType+"/m3tr1x01",
		url.Values{"content": {text}, "options": {options}})
	if err != nil {
		log.Printf("Error: %v\n", err)
		return fmt.Sprintf("Error: %v\n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading response: %v.", err)
	}

	bodyString := string(body)

	return bodyString
}
