package main

import (
	"log"
	"sort"
	"strings"
)

type Paragraph struct {
	Index int
	Group int
	Text  string
	Genre string
}

var grupoJornalistico = map[int]int{
	1: 1, 2: 1, 3: 6, 4: 3, 5: 4, 6: 5, 8: 3, 9: 6, 10: 1, 11: 0, 12: 2, 13: 5, 14: 5, 15: 1, 16: 0, 17: 0, 18: 1, 19: 2, 20: 1, 21: 4, 22: 1, 23: 4, 24: 5, 25: 1, 26: 5, 28: 5, 29: 4, 39: 1, 40: 5, 41: 0, 43: 1, 47: 4, 61: 1, 62: 1, 63: 1, 64: 1, 65: 1, 66: 3, 67: 6, 69: 5, 70: 5, 71: 5, 72: 0,
}

var grupoLiterario = map[int]int{
	30: 3, 31: 1, 32: 2, 33: 3, 34: 2, 35: 0, 45: 1, 51: 1, 68: 1,
}

var grupoDivulgacao = map[int]int{
	//	7: 1, 27: 3, 36: 6, 37: 3, 38: 6, 42: 0, 44: 3, 46: 0, 48: 6, 49: 0, 50: 0, 52: 3, 53: 3, 54: 6, 55: 6, 56: 6, 57: 0, 58: 3, 59: 3, 60: 3, 73: 6, 74: 0, 75: 0, 76: 4, 77: 1, 78: 1, 79: 1, 80: 1, 81: 2, 82: 0, 83: 1, 84: 1, 85: 4, 86: 1, 87: 1, 88: 1, 89: 2, 90: 2, 91: 1, 92: 5, 93: 1, 94: 2, 95: 5, 96: 1, 97: 1, 98: 5, 99: 0, 100: 2,
	7: 2, 27: 7, 36: 5, 37: 2, 38: 5, 42: 0, 44: 2, 46: 0, 48: 5, 49: 0, 50: 0, 52: 7, 53: 7, 54: 5, 55: 5, 56: 5, 57: 0, 58: 7, 59: 7, 60: 2, 73: 7, 74: 0, 75: 0, 76: 4, 77: 3, 78: 3, 79: 0, 80: 3, 81: 0, 82: 3, 83: 3, 84: 3, 85: 2, 86: 0, 87: 3, 88: 3, 89: 1, 90: 3, 91: 3, 92: 7, 93: 3, 94: 1, 95: 3, 96: 3, 97: 0, 98: 6, 99: 0, 100: 3,
}

type ParagraphOrder []Paragraph

func (a ParagraphOrder) Len() int      { return len(a) }
func (a ParagraphOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ParagraphOrder) Less(i, j int) bool {
	if a[i].Genre == a[j].Genre {
		return a[i].Group < a[j].Group
	} else {
		return a[i].Genre < a[j].Genre
	}
}

var paragraphList = []Paragraph{}

func main() {

	log.Println("starting")

	raw := readFile("/home/sidleal/sid/usp/rastros/rastros100.txt")
	lines := strings.Split(raw, "\n")

	for i, line := range lines {
		if line == "" {
			break
		}
		idx := i + 1

		group := -1
		genre := ""
		if _, found := grupoJornalistico[idx]; found {
			group = grupoJornalistico[idx]
			genre = "Jornalístico"
		} else if _, found := grupoLiterario[idx]; found {
			group = grupoLiterario[idx]
			genre = "Literário"
		} else if _, found := grupoDivulgacao[idx]; found {
			group = grupoDivulgacao[idx]
			genre = "Divulgação Científica"
		}
		paragraphList = append(paragraphList, Paragraph{idx, group, line, genre})
	}

	sort.Sort(ParagraphOrder(paragraphList))

	lastGroup := -1
	for _, p := range paragraphList {
		if lastGroup != p.Group {
			log.Println("\n=============================================", "Gênero:", p.Genre, "Grupo", p.Group, "=============================================")
		} else {
			log.Println("---------------------")
		}
		log.Println(p.Index, p.Text)
		lastGroup = p.Group

	}

}
