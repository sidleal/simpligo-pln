package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	Text    string
}

type SentenceAlign struct {
	From Sentence
	To   Sentence
}

func main() {

	strong := true

	path := "/home/sidleal/usp/PROPOR2018/"

	a := senter.ParseText("bla.")
	log.Println(a)

	productions := []Production{}

	for i := 1; i <= 1; i++ { //165

		sentencesAlign := []SentenceAlign{}
		production := fmt.Sprintf("production%v", i)

		alignPath := path + production + "/" + production + "-align.xml"
		if strong {
			alignPath = path + production + "/" + production + "_natural-align.xml"
		}
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
					sentencesAlign = append(sentencesAlign, SentenceAlign{Sentence{matchesFrom[0][0], 0, 0, ""}, Sentence{matchesTo[0][0], 0, 0, ""}})
				} else { //1 pra n
					beg, _ := strconv.Atoi(matchesTo[0][2])
					end, _ := strconv.Atoi(matchesTo[1][2])
					for i := beg; i <= end; i++ {
						sentencesAlign = append(sentencesAlign, SentenceAlign{Sentence{matchesFrom[0][0], 0, 0, ""}, Sentence{fmt.Sprintf("%v%v", matchesTo[0][1], i), 0, 0, ""}})
						log.Println("----", matchesTo[0][1], beg, end, i)
					}
				}
			} else { //n pra 1
				beg, _ := strconv.Atoi(matchesFrom[0][2])
				end, _ := strconv.Atoi(matchesFrom[1][2])
				for i := beg; i <= end; i++ {
					sentencesAlign = append(sentencesAlign, SentenceAlign{Sentence{fmt.Sprintf("%v%v", matchesFrom[0][1], i), 0, 0, ""}, Sentence{matchesTo[0][0], 0, 0, ""}})

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
			// log.Println("------------xxx------------", prd.Name)
			// log.Println(`struct type="s" from="(?P<from>[^"]+).*to="(?P<to>[^"]+)".*\n.*value="` + s.From.Id + `"`)

			// log.Println("------------xxx------------")
			fromT, _ := strconv.ParseInt(match[1], 10, 64)
			toT, _ := strconv.ParseInt(match[2], 10, 64)
			productions[i].Sentences[j].From.FromPos = fromT
			productions[i].Sentences[j].From.ToPos = toT

			regEx = regexp.MustCompile(`struct type="s" from="(?P<from>[^"]+).*to="(?P<to>[^"]+)".*\n.*value="` + s.To.Id + `"`)
			match = regEx.FindStringSubmatch(toMap)

			if len(match) == 0 {
				log.Println("ERRO", prd.Name, prd.FromMap, prd.ToMap, `struct type="s" from="(?P<from>[^"]+).*to="(?P<to>[^"]+)".*\n.*value="`+s.To.Id+`"`)
				continue
			}

			fromT, _ = strconv.ParseInt(match[1], 10, 64)
			toT, _ = strconv.ParseInt(match[2], 10, 64)
			productions[i].Sentences[j].To.FromPos = fromT
			productions[i].Sentences[j].To.ToPos = toT

		}
	}

	for i, prd := range productions {
		fromTokensFrom := readFile(path + prd.Name + "/" + strings.Replace(prd.FromMap, "-s.xml", "-token.xml", -1))
		fromTokensTo := readFile(path + prd.Name + "/" + strings.Replace(prd.ToMap, "-s.xml", "-token.xml", -1))

		for j, s := range prd.Sentences {
			log.Println("----------------------------")

			//FROM
			if s.From.ToPos < s.From.FromPos {
				s.From.ToPos = productions[i].Sentences[j+1].From.ToPos
				if s.From.ToPos < s.From.FromPos {
					s.From.ToPos = productions[i].Sentences[j+2].From.ToPos
				}
				if s.From.ToPos < s.From.FromPos {
					s.From.ToPos = productions[i].Sentences[j+3].From.ToPos
				}
				if s.From.ToPos < s.From.FromPos {
					s.From.ToPos = productions[i].Sentences[j+4].From.ToPos
				}
			}
			if j > 0 && s.From.FromPos == 37 {
				s.From.FromPos = productions[i].Sentences[j-1].From.FromPos
				if s.From.FromPos == 37 {
					s.From.FromPos = productions[i].Sentences[j-2].From.FromPos
				}
			}

			if j > 0 && productions[i].Sentences[j-1].From.ToPos < productions[i].Sentences[j-1].From.FromPos && s.From.FromPos == (productions[i].Sentences[j-1].From.ToPos+1) {
				s.From.FromPos = productions[i].Sentences[j-1].From.FromPos
			}

			//TO
			if s.To.ToPos < s.To.FromPos && j+1 < len(productions[i].Sentences) {
				s.To.ToPos = productions[i].Sentences[j+1].To.ToPos
				if s.To.ToPos < s.To.FromPos && j+2 < len(productions[i].Sentences) {
					s.To.ToPos = productions[i].Sentences[j+2].To.ToPos
				}
				if s.To.ToPos < s.To.FromPos && j+3 < len(productions[i].Sentences) {
					s.To.ToPos = productions[i].Sentences[j+3].To.ToPos
				}
				if s.To.ToPos < s.To.FromPos && j+4 < len(productions[i].Sentences) {
					s.To.ToPos = productions[i].Sentences[j+4].To.ToPos
				}
			}
			if j > 0 && s.To.FromPos == 37 {
				s.To.FromPos = productions[i].Sentences[j-1].To.FromPos
				if s.To.FromPos == 37 {
					s.To.FromPos = productions[i].Sentences[j-2].To.FromPos
				}
			}

			if j > 0 && productions[i].Sentences[j-1].To.ToPos < productions[i].Sentences[j-1].To.FromPos && s.To.FromPos == (productions[i].Sentences[j-1].To.ToPos+1) {
				s.To.FromPos = productions[i].Sentences[j-1].To.FromPos
			}

			// if j > 0 && s.From.ToPos < productions[i].Sentences[j-1].From.FromPos {
			// 	s.From.ToPos = productions[i].Sentences[j-1].From.ToPos + 1
			// }

			regExStr := fmt.Sprintf(`(?s)from='%v'(?P<c>.*to='%v'.*?<\/struct>)`, s.From.FromPos, s.From.ToPos)
			regEx := regexp.MustCompile(regExStr)
			matchFrom := regEx.FindStringSubmatch(fromTokensFrom)

			regExStr = fmt.Sprintf(`(?s)from='%v'(?P<c>.*to='%v'.*?<\/struct>)`, s.To.FromPos, s.To.ToPos)
			regEx = regexp.MustCompile(regExStr)
			matchTo := regEx.FindStringSubmatch(fromTokensTo)

			textFrom := ""
			textTo := ""
			regEx = regexp.MustCompile(`name='base' value='(?P<val>.*)'`)

			if len(matchFrom) == 0 {
				log.Println("---------ERRO - não encontrado: from ", s.From.FromPos, s.From.ToPos)
				textFrom = "ERRO"
			} else {
				matchesFrom := regEx.FindAllStringSubmatch(matchFrom[1], -1)
				for _, match := range matchesFrom {
					textFrom += match[1] + " "
				}
			}
			productions[i].Sentences[j].From.Text = textFrom

			if len(matchTo) == 0 {
				log.Println("---------ERRO - não encontrado: from ", s.To.FromPos, s.To.ToPos)
				textTo = "ERRO"
			} else {
				matchesTo := regEx.FindAllStringSubmatch(matchTo[1], -1)
				for _, match := range matchesTo {
					textTo += match[1] + " "
				}
			}
			productions[i].Sentences[j].To.Text = textTo

			log.Println(i, j, "\n", textFrom, "\n", textTo)
		}

	}

	regEx1 := regexp.MustCompile(`(?P<t>[A-zà\$áã]+)(?P<s>=)`)
	regEx2 := regexp.MustCompile(`(?P<t>[A-z0-9\)"]+) (?P<s>[\)\.,])`)
	regEx3 := regexp.MustCompile(`(?P<s>[A-z0-9\s]+[\(\-]) (?P<t>[A-z0-9]+)`)
	regEx4 := regexp.MustCompile(`(")\s(.*?)\s(")`)
	for i, prd := range productions {
		log.Println("===========================", prd.Name, "=====================")
		for j, s := range prd.Sentences {
			from := regEx1.ReplaceAllString(s.From.Text, `$1 `)
			to := regEx1.ReplaceAllString(s.To.Text, `$1 `)

			from = regEx2.ReplaceAllString(from, `$1$2`)
			to = regEx2.ReplaceAllString(to, `$1$2`)

			from = regEx3.ReplaceAllString(from, `$1$2`)
			to = regEx3.ReplaceAllString(to, `$1$2`)

			from = regEx4.ReplaceAllString(from, `$1$2$3`)
			to = regEx4.ReplaceAllString(to, `$1$2$3`)

			from = regEx2.ReplaceAllString(from, `$1$2`)
			to = regEx2.ReplaceAllString(to, `$1$2`)

			productions[i].Sentences[j].From.Text = from
			productions[i].Sentences[j].To.Text = to

			log.Println(i, j, s.From.FromPos, s.From.ToPos, s.To.FromPos, s.To.ToPos, ":\n\n", from, "\n", to, "\n")

		}
	}

	f, err := os.Create("/home/sidleal/usp/align4-str.txt")
	if err != nil {
		log.Println("ERRO", err)
	}

	defer f.Close()

	for i, prd := range productions {
		lastSentFrom := ""
		lastSentTo := ""
		for j, s := range prd.Sentences {
			if s.From.Text != "ERRO" && s.To.Text != "ERRO" && s.From.Text != s.To.Text {

				if lastSentFrom == s.From.Text && lastSentTo == s.To.Text {
					//ignore dupls
				} else {

					if senter.ParseText(s.From.Text).TotalSentences > 1 || senter.ParseText(s.To.Text).TotalSentences > 1 {
						log.Println("-----> mais de uma sentença")
						log.Println(s.From.Text)
						log.Println(s.To.Text)
					} else {
						n, err := f.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", i, prd.Name, j, s.From.Text, s.To.Text))
						if err != nil {
							log.Println("ERRO", err)
						}
						fmt.Printf("wrote %d bytes\n", n)
					}
					lastSentFrom = s.From.Text
					lastSentTo = s.To.Text
				}

			}
		}
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
