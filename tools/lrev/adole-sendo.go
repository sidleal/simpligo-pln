package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func main_arquivo_unico() {

	log.Println("starting ---")

	mapParts := map[string]string{}

	meta := readFile("/home/sidleal/sid/usp/adole-sendo/metadados.tsv", false)
	lines := strings.Split(meta, "\n")
	for i, line := range lines {
		tokens := strings.Split(line, "\t")
		if i == 0 {
			for j, t := range tokens {
				log.Println(j, t)
			}
			continue
		}
		partID := tokens[0]
		partYear := tokens[9]
		mapParts[fmt.Sprintf("VF_%v", partID)] = partYear
	}

	// for k, v := range mapParts {
	// 	log.Println(k, v)
	// }

	f, err := os.Create("/home/sidleal/sid/usp/lrev/adole-sendo.tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	_, err = f.WriteString("text_id\tyear\ttext\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	raw := readFile("/home/sidleal/sid/usp/adole-sendo/versao_entrega_corrigida_sandra.txt", false)
	lines = strings.Split(raw, "\n")

	textID := ""
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.Index(line, " ") < 1 {
			textID = strings.ReplaceAll(line, "\r", "")
			continue
		}

		text := line

		// log.Println(i, textID, text)

		year := strings.TrimSpace(mapParts[textID])

		log.Println(textID, year, text)

		_, err = f.WriteString(fmt.Sprintf("%v\t%v\t%v\n", textID, year, text))
		if err != nil {
			log.Println("ERRO", err)
		}

		log.Println("---------------------------")

	}

}

func readFile(path string, iso bool) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	r := io.Reader(f)
	if iso {
		r = charmap.ISO8859_1.NewDecoder().Reader(f)
	}

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

func main_extract_adolesendo() {

	log.Println("starting")

	fout, err := os.OpenFile("/home/sidleal/sid/usp/lrev/adole-sendo_metrics.tsv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	data := readFile("/home/sidleal/sid/usp/lrev/adole-sendo.tsv", false)

	lines := strings.Split(data, "\n")
	for i, line := range lines {
		if line == "" || i == 0 {
			continue
		}
		tokens := strings.Split(line, "\t")
		textID := tokens[0]
		year := tokens[1]
		text := tokens[2]

		text = cleanText(text)

		// if textID != "VF_A066" && textID != "VF_A089" {
		// 	continue
		// }
		log.Println(text)

		fRet := callMetrix(text)
		log.Println("-------------------------------------------------------")
		log.Println("-->", i)
		log.Println(fRet)

		feats := strings.Split(fRet, ",")

		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.ReplaceAll(text, "\r", " ")
		header := "text_id\tyear\ttext\t"
		ret := fmt.Sprintf("%v\t%v\t%v\t", textID, year, text)
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

func callMetrix(text string) string {

	timeout := time.Duration(15 * time.Minute)
	client := http.Client{
		Timeout: timeout,
	}

	// resp, err := client.Post("http://fw.nilc.icmc.usp.br:23380/api/v1/metrix/all/m3tr1x01?format=plain", "text", bytes.NewReader([]byte(text)))
	// resp, err := client.Post("http://simpligo.sidle.al:23380/api/v1/metrix/all/m3tr1x01?format=plain", "text", bytes.NewReader([]byte(text)))
	resp, err := client.Post("http://localhost:23380/api/v1/metrix/all/m3tr1x01?format=plain", "text", bytes.NewReader([]byte(text)))

	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error parsing response. %v", err)
	}

	ret := string(body)
	log.Println("->", ret, "<-")

	retTokens := strings.Split(ret, "++")
	if len(retTokens) < 2 {
		return ""
	}
	ret = strings.TrimSpace(retTokens[1])

	return ret
}

func cleanText(text string) string {
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
	text = strings.ReplaceAll(text, " À ", "{{CRASE}}")
	text = strings.ReplaceAll(text, "À", "A")
	text = strings.ReplaceAll(text, "{{CRASE}}", " À ")

	text = strings.ReplaceAll(text, "``", "\"")
	text = strings.ReplaceAll(text, "''", "\"")
	text = strings.ReplaceAll(text, "`", "\"")
	text = strings.ReplaceAll(text, "´", "\"")
	text = strings.ReplaceAll(text, "º", "o")
	text = strings.ReplaceAll(text, "ª", "a")
	text = strings.ReplaceAll(text, "ĉ", "c")
	text = strings.ReplaceAll(text, "ý", "y")
	text = strings.ReplaceAll(text, "\\", "/")
	text = strings.ReplaceAll(text, "♪", "")
	text = strings.ReplaceAll(text, "ă", "ã")
	text = strings.ReplaceAll(text, "ò", "o")
	text = strings.ReplaceAll(text, "Ò", "O")
	text = strings.ReplaceAll(text, "å", "a")
	text = strings.ReplaceAll(text, "ř", "r")
	text = strings.ReplaceAll(text, "ő", "o")
	text = strings.ReplaceAll(text, "Û", "U")
	text = strings.ReplaceAll(text, "û", "u")
	text = strings.ReplaceAll(text, "ẽ", "e")

	return text
}
