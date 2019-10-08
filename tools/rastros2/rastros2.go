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
)

func mainx() {

	log.Println("starting")

	f, err := os.OpenFile("/home/sidleal/sid/usp/rastros/maisum/rastros150.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	raw := readFile("/home/sidleal/sid/usp/rastros/maisum/rastros150.txt")
	lines := strings.Split(raw, "\n")

	for i, line := range lines {
		if line == "" {
			break
		}
		log.Println(line)

		// retRanker := callRanker(line)

		fRet := callMetrix(line)

		feats := strings.Split(fRet, ",")

		header := ""
		ret := ""
		for _, feat := range feats {
			kv := strings.Split(feat, ":")
			if len(kv) > 1 {
				if i == 0 {
					header += kv[0] + ","
				}
				ret += kv[1] + ","
			}
		}

		if i == 0 {
			header = strings.TrimRight(header, ",")
			// header += "complexity"
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

func callMetrix(text string) string {

	resp, err := http.Post("https://simpligo.sidle.al/api/v1/metrix/all/m3tr1x01", "text", bytes.NewReader([]byte(text)))
	if err != nil {
		return fmt.Sprintf("Error extracting metrics: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error parsing response. %v", err)
	}

	ret := string(body)

	return ret
}

func callRanker(text string) string {

	resp, err := http.Post("https://simpligo.sidle.al/api/v1/sentence-ranker/m3tr1x01", "text", bytes.NewReader([]byte(text)))
	if err != nil {
		return fmt.Sprintf("Error extracting metrics: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error parsing response. %v", err)
	}

	ret := string(body)

	return ret
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
