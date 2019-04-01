package senter

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

// func main() {
// 	a := `
// 	Muito interessante. A coisa funciona. O número 23.123, e o http://coisa.com. Isso mesmo.
// 	Também o sid.leal@gmail.com. Pois bem... Assim que aconteceu? Foi o Prof. João?
// 	Acho que foi! E ainda assim funciona. Ele disse: "Aqui nao pode ter quebra. Mesmo tendo ponto."
// 	E nem (aqui. Ignore isso.).
// 	Nascido em 230 A.C. com louvor.
// 	O Neil M. Ferguson apareceu naquela época.

// 	Ferguson disse: "Você precisa de uma combinação de estratégias. Nenhuma sozinha conseguiria prevenir com sucesso uma epidemia."`

// 	parsedText := ParseText(a)
// 	log.Println(parsedText)

// 	jsonText := ParseTextToJson(a)
// 	log.Println(jsonText)

// 	for _, p := range parsedText.Paragraphs {
// 		for _, s := range p.Sentences {
// 			log.Println(s.Text)
// 		}
// 	}

// }

type ParsedText struct {
	Paragraphs      []ParsedParagraph `json:"paragraphs"`
	TotalParagraphs int64             `json:"totp"`
	TotalSentences  int64             `json:"tots"`
	TotalTokens     int64             `json:"tott"`
	TotalWords      int64             `json:"totw"`
}

type ParsedParagraph struct {
	Idx          int64            `json:"idx"`
	Sentences    []ParsedSentence `json:"sentences"`
	Text         string           `json:"txt"`
	QtyTokens    int64            `json:"qtt"`
	QtyWords     int64            `json:"qtw"`
	QtySentences int64            `json:"qts"`
}

type ParsedSentence struct {
	Idx       int64         `json:"idx"`
	Tokens    []ParsedToken `json:"tokens"`
	Text      string        `json:"txt"`
	QtyTokens int64         `json:"qtt"`
	QtyWords  int64         `json:"qtw"`
}

type ParsedToken struct {
	Idx    int64  `json:"idx"`
	Token  string `json:"token"`
	IsWord int64  `json:"w"`
}

var abbreviations []string = []string{
	"Prof.", "A.C.", "Jr.", "a.C.",
}

func ParseText(rawText string) ParsedText {
	preText := preProcesText(rawText)
	processedText := processText(preText)
	parsedText := postProcess(processedText)
	return parsedText
}

func ParseTextToJson(rawText string) string {
	parsedText := ParseText(rawText)
	ret, err := json.Marshal(parsedText)
	if err != nil {
		log.Println(err)
	}
	return string(ret)
}

/*
Rules:
1. Delimita-se uma sentença sempre que uma marca de nova linha (carriage return e line
feed) é encontrada, independentemente de um sinal de fim de sentença ter sido encontrado
anteriormente;
2. Não se delimitam sentenças dentro de aspas, parênteses, chaves e colchetes;
3. Delimita-se uma sentença quando os símbolos de interrogação (?) e exclamação (!) são
encontrados;
4. Delimita-se uma sentença quando o símbolo de ponto (.) é encontrado e este não é um
ponto de número decimal, não pertence a um símbolo de reticências (...), não faz parte de
endereços de e-mail e páginas da Internet e não é o ponto que segue uma abreviatura;
5. Delimita-se uma sentença quando uma letra maiúscula é encontrada após o sinal de
reticências ou de fecha-aspas.
(Pardo, 2006)
6. Abreviações de nome próprio, espaçoMaisculaPonto
*/

func preProcesText(rawText string) string {
	out := rawText + "\n"

	// rule 2 - " { [ ( ) ] } "
	out = applyGroupRule(out, `"([^"]+?)"[^A-z]`)
	out = applyGroupRule(out, `“(.+?)”`)
	out = applyGroupRule(out, `\{(.+?)\}`)
	out = applyGroupRule(out, `\[(.+?)\]`)
	out = applyGroupRule(out, `\((.+?)\)`)

	// rule 4 - abbreviations
	for _, abbrev := range abbreviations {
		abbrevNew := strings.Replace(abbrev, ".", "|dot|", -1)
		out = strings.Replace(out, abbrev, abbrevNew, -1)
	}

	// rule 4 - internet
	regEx := regexp.MustCompile(`(http|ftp|www|@)(?P<dot>.+?)(,|\.)*\s`)
	matches := regEx.FindAllStringSubmatch(out, -1)
	for _, match := range matches {
		addressOld := match[2]
		addressNew := strings.Replace(addressOld, ".", "|dot|", -1)
		addressNew = strings.Replace(addressNew, ":", "|col|", -1)
		out = strings.Replace(out, addressOld, addressNew, -1)
	}

	regEx = regexp.MustCompile(`([A-z0-9_\-\.]+)@`)
	matches = regEx.FindAllStringSubmatch(out, -1)
	for _, match := range matches {
		emailOld := match[0]
		emailNew := strings.Replace(emailOld, ".", "|dot|", -1)
		out = strings.Replace(out, emailOld, emailNew, -1)
	}

	//rule 4 - decimals
	regEx = regexp.MustCompile(`([0-9]+)\.([0-9]+),([0-9]+)`)
	out = regEx.ReplaceAllString(out, `$1|dot|$2|virg|$3`)

	regEx = regexp.MustCompile(`([0-9]+),([0-9]+)`)
	out = regEx.ReplaceAllString(out, `$1|virg|$2`)

	regEx = regexp.MustCompile(`([0-9]+)\.([0-9]+)`)
	out = regEx.ReplaceAllString(out, `$1|dot|$2`)

	//rule 5 - quotes
	regEx = regexp.MustCompile(`"[\s\n]+?([A-Z])`)
	out = regEx.ReplaceAllString(out, `". $1`)

	//rule 5 - reticences
	regEx = regexp.MustCompile(`\.\.\.\s*([A-Z])`)
	out = regEx.ReplaceAllString(out, `|dot||dot||dot| |||$1`)

	//rule 4 - reticences
	regEx = regexp.MustCompile(`\.\.\.`)
	out = regEx.ReplaceAllString(out, `|dot||dot||dot|`)

	//rule 6 - nomes
	regEx = regexp.MustCompile(`(\s[A-Z]{1})\.(\s)`)
	out = regEx.ReplaceAllString(out, `$1|dot|$2`)

	// rule 3
	out = strings.Replace(out, ".", ".|||", -1)
	out = strings.Replace(out, "?", "?|||", -1)
	out = strings.Replace(out, "!", "!|||", -1)

	return out
}

func applyGroupRule(rawText string, regexGroup string) string {
	regEx := regexp.MustCompile(regexGroup)
	matches := regEx.FindAllStringSubmatch(rawText, -1)
	for _, match := range matches {
		sentenceOld := match[1]
		sentenceNew := strings.Replace(sentenceOld, ".", "|gdot|", -1)
		sentenceNew = strings.Replace(sentenceNew, "?", "|gint|", -1)
		sentenceNew = strings.Replace(sentenceNew, "!", "|gexc|", -1)
		rawText = strings.Replace(rawText, sentenceOld, sentenceNew, -1)
	}
	return rawText
}

func tokenizeText(rawText string) []string {
	regEx := regexp.MustCompile(`([A-z]+)-([A-z]+)`)
	rawText = regEx.ReplaceAllString(rawText, "$1|hyp|$2")

	regEx = regexp.MustCompile(`\|gdot\|`)
	rawText = regEx.ReplaceAllString(rawText, ".")

	regEx = regexp.MustCompile(`\|gint\|`)
	rawText = regEx.ReplaceAllString(rawText, "?")

	regEx = regexp.MustCompile(`\|gexc\|`)
	rawText = regEx.ReplaceAllString(rawText, "!")

	regEx = regexp.MustCompile(`([\.\,"\(\)\[\]\{\}\?\!;:-]{1})`)
	rawText = regEx.ReplaceAllString(rawText, "  $1 ")

	regEx = regexp.MustCompile(`\s+`)
	rawText = regEx.ReplaceAllString(rawText, " ")

	return strings.Split(rawText, " ")
}

func processText(rawText string) ParsedText {
	parsedText := ParsedText{}
	paragraphs := strings.Split(rawText, "\n")

	var idxParagraphs int64 = 0
	var idxSentences int64 = 0
	var idxTokens int64 = 0
	var idxWords int64 = 0

	parsedText.Paragraphs = []ParsedParagraph{}

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p != "" {
			idxParagraphs++
			parsedParagraph := ParsedParagraph{}
			parsedParagraph.Idx = idxParagraphs
			parsedParagraph.Sentences = []ParsedSentence{}
			parsedParagraph.Text = p
			var qtwP int64 = 0
			var qttP int64 = 0
			sentences := strings.Split(p, "|||")
			for _, s := range sentences {
				s = strings.TrimSpace(s)
				if s != "" {
					idxSentences++
					parsedSentence := ParsedSentence{idxSentences, []ParsedToken{}, s, 0, 0}

					tokens := tokenizeText(s)
					var qtw int64 = 0
					var qtt int64 = 0
					for _, t := range tokens {
						if len(t) > 0 {
							idxTokens++
							qtt++
							qttP++
							token := ParsedToken{idxTokens, t, 0}
							if len(t) > 1 || strings.Index(`{[()]}.,"?!;:-'#&`, t) < 0 {
								qtw++
								qtwP++
								idxWords++
								token.IsWord = 1
							}
							parsedSentence.Tokens = append(parsedSentence.Tokens, token)
						}
					}
					parsedSentence.QtyTokens = qtt
					parsedSentence.QtyWords = qtw
					parsedParagraph.Sentences = append(parsedParagraph.Sentences, parsedSentence)
				}
			}
			parsedParagraph.QtyTokens = qttP
			parsedParagraph.QtyWords = qtwP
			parsedParagraph.QtySentences = int64(len(parsedParagraph.Sentences))
			parsedText.Paragraphs = append(parsedText.Paragraphs, parsedParagraph)
		}

	}
	parsedText.TotalParagraphs = idxParagraphs
	parsedText.TotalSentences = idxSentences
	parsedText.TotalTokens = idxTokens
	parsedText.TotalWords = idxWords

	return parsedText
}

func postProcess(parsedText ParsedText) ParsedText {

	for i, p := range parsedText.Paragraphs {
		parsedText.Paragraphs[i].Text = punctuateBack(p.Text)

		for j, s := range p.Sentences {
			parsedText.Paragraphs[i].Sentences[j].Text = punctuateBack(s.Text)
			for k, t := range s.Tokens {
				parsedText.Paragraphs[i].Sentences[j].Tokens[k].Token = punctuateBack(t.Token)
			}
		}
	}
	return parsedText
}

func punctuateBack(text string) string {
	text = strings.Replace(text, "|dot|", ".", -1)
	text = strings.Replace(text, "|int|", "?", -1)
	text = strings.Replace(text, "|exc|", "!", -1)
	text = strings.Replace(text, "|gdot|", ".", -1)
	text = strings.Replace(text, "|gint|", "?", -1)
	text = strings.Replace(text, "|gexc|", "!", -1)
	text = strings.Replace(text, "|hyp|", "-", -1)
	text = strings.Replace(text, "|col|", ":", -1)
	text = strings.Replace(text, "|||", "", -1)
	text = strings.Replace(text, "|virg|", ",", -1)
	return text
}
