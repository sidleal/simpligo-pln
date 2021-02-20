package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	log.Println("starting")

	fout, err := os.OpenFile("/home/sidleal/sid/usp/mjfinatto/metricas.tsv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	path := "/home/sidleal/sid/usp/mjfinatto"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	first := true
	for _, f := range files {
		fileName := f.Name()
		log.Println("---------------------------------------------------------------")
		log.Println(fileName)

		if !strings.HasSuffix(fileName, ".txt") {
			continue
		}

		raw := readFile(path + "/" + fileName)

		log.Println(raw)

		text := cleanText(raw)

		fRet := callMetrix(text)
		log.Println(fRet)

		tRet := strings.Split(fRet, "-------------------------------")

		feats := strings.Split(tRet[1], ",")

		header := "arquivo\ttexto\t"
		text = strings.ReplaceAll(text, "\r", "")
		text = strings.ReplaceAll(text, "\t", "")
		text = strings.ReplaceAll(text, "\n", "\\n")
		ret := fmt.Sprintf("%v\t%v\t", fileName, text)
		for _, feat := range feats {
			kv := strings.Split(feat, ":")
			if len(kv) > 1 {
				if first {
					header += kv[0] + "\t"
				}
				ret += kv[1] + "\t"
			}
		}

		if first {
			header = strings.TrimRight(header, "\t")
			_, err := fout.WriteString(header + "\n")
			if err != nil {
				log.Println("ERRO", err)
			}
			first = false
		}

		ret = strings.TrimRight(ret, "\t")
		_, err := fout.WriteString(ret + "\n")
		if err != nil {
			log.Println("ERRO", err)
		}

	}

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
	text = strings.ReplaceAll(text, "`", "\"")
	text = strings.ReplaceAll(text, "´", "\"")
	return text
}
