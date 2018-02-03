package main

import (
	"fmt"
	"io"
	"log"
	"os"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
	"golang.org/x/text/encoding/charmap"
)

func main() {

	path := "/home/sidleal/usp/PROPOR2018/"

	for i := 1; i <= 1; i++ { //165
		production := fmt.Sprintf("production%v", i)

		originalPath := path + production + "/" + production + ".txt"
		naturalPath := path + production + "/" + production + "_natural.txt"
		strongPath := path + production + "/" + production + "_strong.txt"

		if !fileExists(originalPath) {
			log.Println("ERRO------------------, não existe: ", originalPath)
			continue
		}

		original := readFile(originalPath)
		log.Println(original)

		log.Print("\n\n-------------------------------------------------\n\n")

		natural := readFile(naturalPath)
		log.Println(natural)

		log.Print("\n\n-------------------------------------------------\n\n")

		strong := readFile(strongPath)
		log.Println(strong)

		parsedOriginal := senter.ParseText(original)
		parsedNatural := senter.ParseText(natural)
		parsedStrong := senter.ParseText(strong)

		log.Println(parsedOriginal)
		log.Println(parsedNatural)
		log.Println(parsedStrong)

	}
}

func fileExists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	r := charmap.ISO8859_1.NewDecoder().Reader(f)

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
