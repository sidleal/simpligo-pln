package main

import (
	"log"
	"os"
	"strings"
)

func main() {

	log.Println("starting ---")

	// f, err := os.OpenFile("/home/sidleal/sid/usp/arthur/out.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f, err := os.OpenFile("/home/sidleal/sid/usp/adole-sendo/metrics.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	raw := readFile("/home/sidleal/sid/usp/adole-sendo/versao_entrega_corrigida_sandra.txt")
	lines := strings.Split(raw, "\n")

	textID := ""
	firstLine := true
	for i, line := range lines {
		if line == "" {
			continue
		}
		if strings.Index(line, " ") < 1 {
			textID = strings.ReplaceAll(line, "\r", "")
			continue
		}

		if textID != "VF_A071" && textID != "VF_A074" && textID != "VF_A077" {
			continue
		}

		text := line

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
				if firstLine {
					header += kv[0] + ","
				}
				ret += kv[1] + ","
			}
		}

		if firstLine {
			header = strings.TrimRight(header, ",")
			_, err := f.WriteString(header + "\n")
			if err != nil {
				log.Println("ERRO", err)
			}
			firstLine = false
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
