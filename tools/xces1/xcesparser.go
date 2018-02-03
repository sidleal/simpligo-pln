package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"

	"golang.org/x/text/encoding/charmap"
)

func main() {

	path := "/home/sidleal/usp/PROPOR2018/"

	for i := 1; i <= 165; i++ {
		production := fmt.Sprintf("production%v", i)

		alignPath := path + production + "/" + production + "-align.xml"
		if !fileExists(alignPath) {
			log.Println("ERRO------------------, não existe: ", alignPath)
			continue
		}

		align := readFile(alignPath)
		// log.Println(align)

		var retCmdRegEx *regexp.Regexp = regexp.MustCompile(`fromDoc="(?P<from>[^"]+)".*toDoc="(?P<to>[^"]+)"`)
		match := retCmdRegEx.FindStringSubmatch(align)
		from := match[1]
		to := match[2]
		log.Println(production, "From: ", from, "To:", to)

		fromMap := readFile(path + production + "/" + from)
		// log.Println(fromMap)

		retCmdRegEx = regexp.MustCompile(`struct type="s" from="(?P<from>[^"]+).*to="(?P<to>[^"]+)".*\n.*value="(?P<val>[^"]+)"`)
		matches := retCmdRegEx.FindAllStringSubmatch(fromMap, -1)
		for _, match := range matches {
			sFrom := match[1]
			sTo := match[2]
			sVal := match[3]

			x, _ := strconv.Atoi(sFrom)
			y, _ := strconv.Atoi(sTo)
			if x > y {
				log.Println("-------------------------------Erro ALIGN", production, sFrom, sTo, sVal)
			}
			log.Println("From: ", sFrom, "To:", sTo, "Val:", sVal)
		}
	}

	// from := readFile("/home/sidleal/usp/PROPOR2018/production108/production108.txt")
	// log.Println(from)

	// to := readFile("/home/sidleal/usp/PROPOR2018/production108/production108-natural.txt")
	// log.Println(to)

	// align := readFile("/home/sidleal/usp/PROPOR2018/production108/production108-align.xml")
	// log.Println(align)

	// from, err := ioutil.ReadFile("/home/sidleal/usp/PROPOR2018/production108/production108.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(from))

	// to, err := ioutil.ReadFile("/home/sidleal/usp/PROPOR2018/production108/production108.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(to))

	// align, err := ioutil.ReadFile("/home/sidleal/usp/PROPOR2018/production108/production108-align.xml")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(align))

	// log.Println("testex")
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
