package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"

	senter "github.com/sidleal/simpligo-pln/tools/senter"
	"golang.org/x/text/encoding/charmap"
)

type Production struct {
	Name      string
	FromMap   string
	ToMap     string
	Sentences []SentenceAlign
}

type Sentence struct {
	Id      string
	FromPos int64
	ToPos   int64
}

type SentenceAlign struct {
	From Sentence
	To   Sentence
}

func main() {

	path := "/home/sidleal/usp/PROPOR2018/"

	a := senter.ParseText("bla.")
	log.Println(a)

	productions := []Production{}

	for i := 1; i <= 165; i++ {

		sentencesAlign := []SentenceAlign{}
		production := fmt.Sprintf("production%v", i)

		alignPath := path + production + "/" + production + "-align.xml"
		if !fileExists(alignPath) {
			log.Println("ERRO------------------, não existe: ", alignPath)
			continue
		}

		align := readFile(alignPath)
		// log.Println(align)

		var regEx *regexp.Regexp = regexp.MustCompile(`fromDoc="(?P<from>[^"]+)".*toDoc="(?P<to>[^"]+)"`)
		match := regEx.FindStringSubmatch(align)
		fromMap := match[1]
		toMap := match[2]
		log.Println(production, "From: ", fromMap, "To:", toMap)

		regEx = regexp.MustCompile(`href="#(?P<from>[^"]+).*\n.*xlink:href="#(?P<to>[^"]+)"`)
		matches := regEx.FindAllStringSubmatch(align, -1)
		for _, match := range matches {
			fromS := match[1]
			toS := match[2]
			log.Println(production, "From: ", fromS, "To:", toS)

			regEx = regexp.MustCompile(`(?P<a>p[0-9]+s)(?P<s>[0-9]+)`)
			matchesFrom := regEx.FindAllStringSubmatch(fromS, -1)
			matchesTo := regEx.FindAllStringSubmatch(toS, -1)
			if len(matchesFrom) == 1 {
				if len(matchesTo) == 1 {
					sentencesAlign = append(sentencesAlign, SentenceAlign{Sentence{matchesFrom[0][0], 0, 0}, Sentence{matchesTo[0][0], 0, 0}})
				} else { //1 pra n
					beg, _ := strconv.Atoi(matchesTo[0][2])
					end, _ := strconv.Atoi(matchesTo[1][2])
					for i := beg; i <= end; i++ {
						sentencesAlign = append(sentencesAlign, SentenceAlign{Sentence{matchesFrom[0][0], 0, 0}, Sentence{fmt.Sprintf("%v%v", matchesTo[0][1], i), 0, 0}})
						log.Println("----", matchesTo[0][1], beg, end, i)
					}
				}
			} else { //n pra 1
				beg, _ := strconv.Atoi(matchesFrom[0][2])
				end, _ := strconv.Atoi(matchesFrom[1][2])
				for i := beg; i <= end; i++ {
					sentencesAlign = append(sentencesAlign, SentenceAlign{Sentence{fmt.Sprintf("%v%v", matchesFrom[0][1], i), 0, 0}, Sentence{matchesTo[0][0], 0, 0}})

					log.Println("x----", matchesFrom[0][1], beg, end, i)
				}
			}

		}

		productions = append(productions, Production{production, fromMap, toMap, sentencesAlign})
	}

	for i, prd := range productions {
		fromMap := readFile(path + prd.Name + "/" + prd.FromMap)
		toMap := readFile(path + prd.Name + "/" + prd.ToMap)

		for j, s := range prd.Sentences {

			regEx := regexp.MustCompile(`struct type="s" from="(?P<from>[^"]+).*to="(?P<to>[^"]+)".*\n.*value="` + s.From.Id + `"`)
			match := regEx.FindStringSubmatch(fromMap)
			fromT, _ := strconv.ParseInt(match[1], 10, 64)
			toT, _ := strconv.ParseInt(match[2], 10, 64)
			productions[i].Sentences[j].From.FromPos = fromT
			productions[i].Sentences[j].From.ToPos = toT

			log.Println(prd.Name, prd.FromMap, prd.ToMap, `struct type="s" from="(?P<from>[^"]+).*to="(?P<to>[^"]+)".*\n.*value="`+s.To.Id+`"`)

			regEx = regexp.MustCompile(`struct type="s" from="(?P<from>[^"]+).*to="(?P<to>[^"]+)".*\n.*value="` + s.To.Id + `"`)
			match = regEx.FindStringSubmatch(toMap)
			fromT, _ = strconv.ParseInt(match[1], 10, 64)
			toT, _ = strconv.ParseInt(match[2], 10, 64)
			productions[i].Sentences[j].To.FromPos = fromT
			productions[i].Sentences[j].To.ToPos = toT

		}
	}

	for i, prd := range productions {
		log.Println(i, prd)
	}

	// for _, match := range matches {
	// 	fromS := match[1]
	// 	toS := match[2]
	// 	log.Println(production, "From: ", fromS, "To:", toS)

	// regEx = regexp.MustCompile(`struct type="s" from="(?P<from>[^"]+).*to="(?P<to>[^"]+)".*\n.*value="` + fromS + `"`)
	// match := regEx.FindStringSubmatch(fromMap)
	// fromT := match[1]
	// toT := match[2]
	// log.Println(production, "From: ", fromT, "To:", toT)

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
