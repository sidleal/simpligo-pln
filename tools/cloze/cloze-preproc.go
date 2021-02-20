package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type Word struct {
	Code     string // Código
	Name     string // Nome Teste
	QtdGenre string // Quantidade Gêneros
	Pars     string // Parágrafos por Participante
	Part     string // Nome Participante
	Email    string // Email
	Age      string // Age
	Gender   string // Gender
	Reg      string // Registro
	Sem      string // Semestre

	Org      string // Instituição
	Course   string // Curso
	Language string // Lingua
	Phone    string // Fone
	CPF      string // cpf

	ParsRead string // Parágrafos Lidos
	DTBegin  string // Data Início
	HRBegin  string // Hora Início
	ParID    int    // Parágrafo
	SentID   int    // Sentença
	WordID   int    // Índice Palavra
	RawWord  string // Palavra Crua
	Word     string // Palavra
	Resp     string // Resposta
	TBegin   int    // Tempo Início(ms)
	TDig     int    // Tempo Digitação(ms)
	TTot     int    // Tempo(ms)
	TPar     int    // Tempo Acumulado Parágrafo(ms)
	TTest    int    // Tempo Acumulado Teste(ms)
}

type WordOrder []Word

func (a WordOrder) Len() int      { return len(a) }
func (a WordOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a WordOrder) Less(i, j int) bool {
	if a[i].Reg == a[j].Reg && a[i].ParID == a[j].ParID {
		return a[i].WordID < a[j].WordID
	} else if a[i].Reg == a[j].Reg {
		return a[i].ParID < a[j].ParID
	} else {
		return a[i].Reg < a[j].Reg
	}
}

type FinalWordOrder []Word

func (a FinalWordOrder) Len() int      { return len(a) }
func (a FinalWordOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FinalWordOrder) Less(i, j int) bool {
	if a[i].WordID == a[j].WordID && a[i].ParID == a[j].ParID {
		return a[i].Reg < a[j].Reg
	} else if a[i].ParID == a[j].ParID {
		return a[i].WordID < a[j].WordID
	} else {
		return a[i].ParID < a[j].ParID
	}
}

func readFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	// r := charmap.ISO8859_1.NewDecoder().Reader(f)

	ret := ""

	buf := make([]byte, 32*1024)
	for {
		// n, err := r.Read(buf)
		n, err := f.Read(buf)
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

func readFileISO(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	r := charmap.ISO8859_1.NewDecoder().Reader(f)

	ret := ""

	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		// n, err := f.Read(buf)
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

func readFileLines(path string) []string {
	ret := []string{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ret

}

// ------------------------------------

func main_preproc() {

	rawFiles := []string{
		"cloze_baJwn24BmRI7xu8F5xxd_data25.csv", //puc
		"cloze_bqJxn24BmRI7xu8Fqhz6_data25.csv", //usp
		"cloze_bKJwn24BmRI7xu8FSBzU_data25.csv", //ufc
		"cloze_UgXZ-28BaYrNtDuxSkMf_data25.csv", //utfpr
		"cloze_ogW4nXABaYrNtDuxrk04_data25.csv", //ufabc
		"cloze_YwXxiHABaYrNtDux7Uh8_data25.csv", //uerj
	}

	path := "/home/sidleal/sid/usp/cloze_exps3/data/"

	exportDate := "2020_11_26"
	outFiles := []string{
		"cloze_puc_" + exportDate + ".csv",
		"cloze_usp_" + exportDate + ".csv",
		"cloze_ufc_" + exportDate + ".csv",
		"cloze_utfpr_" + exportDate + ".csv",
		"cloze_ufabc_" + exportDate + ".csv",
		"cloze_uerj_" + exportDate + ".csv",
	}

	f2, err := os.Create(path + "cloze_status_" + exportDate + ".tsv")
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f2.Close()

	_, err = f2.WriteString("arquivo\tnome\te-mail\torg\tidade\tgênero\tcurso\tsemestre\ttextos\tqtde\trespondidos\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	total_lines_original := 0

	for k, rf := range rawFiles {
		data := readFile(path + rf)

		mapParticipants := map[string]map[int]int{}

		f, err := os.Create(path + outFiles[k])
		if err != nil {
			log.Println("ERRO", err)
		}

		wordList := []Word{}

		lines := strings.Split(data, "\n")
		for i, line := range lines {
			if line == "" {
				continue
			}

			cols := strings.Split(line, "\t")
			if i == 0 {
				// log.Println(line)

				_, err = f.WriteString(line + "\n")
				if err != nil {
					log.Println("ERRO", err)
				}

				// for j, col := range cols {
				// 	log.Println(j, "->", col)
				// }
				continue
			}

			// log.Println("-----", line)

			word := Word{}
			word.Code = cols[0]
			word.Name = cols[1]
			word.QtdGenre = cols[2]
			word.Pars = cols[3]
			word.Part = cols[4]
			word.Email = cols[5]
			word.Age = cols[6]
			word.Gender = cols[7]
			word.Reg = cols[8]
			word.Sem = cols[9]

			word.Org = cols[10]
			word.Course = cols[11]
			word.Language = cols[12]
			word.Phone = cols[13]
			word.CPF = cols[14]

			word.ParsRead = cols[15]
			word.DTBegin = cols[16]
			word.HRBegin = cols[17]
			word.ParID, _ = strconv.Atoi(cols[18])
			word.SentID, _ = strconv.Atoi(cols[19])
			word.WordID, _ = strconv.Atoi(cols[20])
			word.RawWord = cols[21]
			word.Word = cols[22]
			word.Resp = cols[23]
			word.TBegin, _ = strconv.Atoi(cols[24])
			word.TDig, _ = strconv.Atoi(cols[25])
			word.TTot, _ = strconv.Atoi(cols[26])
			word.TPar, _ = strconv.Atoi(cols[27])
			word.TTest, _ = strconv.Atoi(cols[28])

			total_lines_original++

			//descarta participantes:
			discard := map[string]int{
				"362726343":    1, // teste UFABC
				"78093609":     1, // teste utfpr
				"93002194465":  1, // teste ufc
				"55.019.616-x": 1, // não entendeu a tarefa - UTFPR
				"098512262":    1, // não entendeu a tarefa  -  UERJ
				"350330001":    1, // não entendeu a tarefa  - USP
				"369128771":    1, // não entendeu a tarefa  - UFABC
			}
			if _, found := discard[word.Reg]; found {
				continue
			}

			//ajuste... ufabc e uerj não tem o primeiro paragrafo
			if strings.Index(outFiles[k], "ufabc") > 0 || strings.Index(outFiles[k], "uerj") > 0 {
				word.ParID++
			}

			discardParagraphs := map[string][]int{
				"12.294.189-1":  []int{2},  // falha alinhamento (erro browser)
				"17976454":      []int{2},  // falha alinhamento (erro browser)
				"V528596-D":     []int{32}, // falha alinhamento (erro browser)
				"2000012026485": []int{28}, // falha alinhamento (erro browser)

			}
			if _, found := discardParagraphs[word.Reg]; found {
				itemToDiscard := false
				for _, par := range discardParagraphs[word.Reg] {
					if word.ParID == par {
						itemToDiscard = true
					}
				}
				if itemToDiscard {
					continue
				}
			}

			//another cleaning
			word.Resp = strings.ReplaceAll(word.Resp, `"`, "")

			wordList = append(wordList, word)

			gender := word.Gender
			gender = strings.ToLower(gender)
			gender = strings.TrimSpace(gender)
			mapGender := map[string]string{
				"feminino":        "F",
				"femenino":        "F",
				"feminio":         "F",
				"femino":          "F",
				"fem":             "F",
				"f":               "F",
				"ferminino":       "F",
				"masculino":       "M",
				"homem":           "M",
				"masc":            "M",
				"m":               "M",
				"masculimo":       "M",
				"homem cisgenero": "O",
			}
			if val, found := mapGender[gender]; found {
				gender = val
			}

			partKey := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s", word.Part, word.Email, word.Org, word.Age, gender, word.Course, word.Sem)
			if _, found := mapParticipants[partKey]; !found {
				mapParticipants[partKey] = map[int]int{}
			}
			if _, found := mapParticipants[partKey][word.ParID]; !found {
				mapParticipants[partKey][word.ParID] = 0
			}
			mapParticipants[partKey][word.ParID]++

		}

		sort.Sort(WordOrder(wordList))

		cont := 0
		lastPar := 0
		for _, w := range wordList {
			if lastPar != w.ParID {
				cont = 0
			}
			lastPar = w.ParID

			newLine := ""
			// 2sid = 11joao, 32sid = 33joao, 39sid = 36joao
			// 4sid = 13joao,28sid = 26joao, 50sid = 26joao
			// 44sid = 44joao
			if w.ParID == 2 {
				if w.WordID < 40 {
					cont = w.WordID
				} else {
					cont++
				}

				if w.WordID == 40 || w.WordID == 43 {
					newLine = formatLine(w, cont)
					cont++
				}

				newLine += formatLine(w, cont)

			} else if w.ParID == 32 {
				if w.WordID < 43 {
					cont = w.WordID
				} else {
					cont++
				}

				if w.WordID == 43 || w.WordID == 48 {
					newLine = formatLine(w, cont)
					cont++
				}

				newLine += formatLine(w, cont)

			} else if w.ParID == 39 {
				if w.WordID < 40 {
					cont = w.WordID
				} else {
					cont++
				}

				if w.WordID == 40 || w.WordID == 45 {
					newLine = formatLine(w, cont)
					cont++
				}

				newLine += formatLine(w, cont)

			} else if w.ParID == 44 {
				if w.WordID < 34 {
					cont = w.WordID
				} else {
					cont++
				}

				if w.WordID == 34 {
					newLine = formatLine(w, cont)
					cont++
				}

				newLine += formatLine(w, cont)

			} else if w.ParID == 4 {
				if w.WordID < 59 {
					cont = w.WordID
				} else {
					cont++
				}

				if w.WordID == 58 {
					cont--
					continue
				}

				newLine = formatLine(w, cont)

			} else if w.ParID == 28 {
				if w.WordID < 41 {
					cont = w.WordID
				} else {
					cont++
				}

				if w.WordID == 41 {
					cont--
					continue
				}

				newLine = formatLine(w, cont)

			} else if w.ParID == 50 {
				if w.WordID < 19 {
					cont = w.WordID
				} else {
					cont++
				}

				if w.WordID == 19 {
					cont--
					continue
				}

				newLine = formatLine(w, cont)

			} else {
				newLine = formatLine(w, w.WordID)
			}

			// log.Print(newLine)
			_, err = f.WriteString(newLine)
			if err != nil {
				log.Println("ERRO", err)
			}
		}

		f.Close()

		// log.Println("============", outFiles[k], "============")
		for key, val := range mapParticipants {
			str := fmt.Sprintf("%v\t%v\t%v\t%v\n", outFiles[k], key, val, len(val))
			// fmt.Println(outFiles[k], "\t", key, "\t", val, "\t", len(val))
			// fmt.Println(str)
			_, err = f2.WriteString(str)
			if err != nil {
				log.Println("ERRO", err)
			}

		}
	}

	wordList := []Word{}
	header := ""
	for k, of := range outFiles {
		log.Println(k, of)
		data := readFile(path + of)

		lines := strings.Split(data, "\n")
		for i, line := range lines {
			if line == "" {
				continue
			}

			cols := strings.Split(line, "\t")
			if i == 0 {
				header = line
				// for j, col := range cols {
				// 	log.Println(j, "->", col)
				// }
				continue
			}

			// log.Println("-----", line)

			word := Word{}
			word.Code = cols[0]
			word.Name = cols[1]
			word.QtdGenre = cols[2]
			word.Pars = cols[3]
			word.Part = cols[4]
			word.Email = cols[5]
			word.Age = cols[6]
			word.Gender = cols[7]
			word.Reg = cols[8]
			word.Sem = cols[9]

			word.Org = cols[10]
			word.Course = cols[11]
			word.Language = cols[12]
			word.Phone = cols[13]
			word.CPF = cols[14]

			word.ParsRead = cols[15]
			word.DTBegin = cols[16]
			word.HRBegin = cols[17]
			word.ParID, _ = strconv.Atoi(cols[18])
			word.SentID, _ = strconv.Atoi(cols[19])
			word.WordID, _ = strconv.Atoi(cols[20])
			word.RawWord = cols[21]
			word.Word = cols[22]
			word.Resp = cols[23]
			word.TBegin, _ = strconv.Atoi(cols[24])
			word.TDig, _ = strconv.Atoi(cols[25])
			word.TTot, _ = strconv.Atoi(cols[26])
			word.TPar, _ = strconv.Atoi(cols[27])
			word.TTest, _ = strconv.Atoi(cols[28])

			wordList = append(wordList, word)

		}

	}

	sort.Sort(FinalWordOrder(wordList))

	finalOutFile := "cloze_all_" + exportDate + ".csv"
	f3, err := os.Create(path + finalOutFile)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f3.Close()

	_, err = f3.WriteString(header + "\n")
	if err != nil {
		log.Println("ERRO", err)
	}

	lastPart := ""
	lastWordID := 0
	for _, w := range wordList {

		// log.Println(i, w)
		if w.Reg == lastPart && w.WordID == lastWordID {
			continue
		}
		lastPart = w.Reg
		lastWordID = w.WordID

		_, err = f3.WriteString(formatLine(w, w.WordID))
		if err != nil {
			log.Println("ERRO", err)
		}
	}
	log.Println("Original: ", total_lines_original, "Após limpeza 1:", len(wordList))

}

func formatLine(w Word, cont int) string {
	return fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		w.Code, w.Name, w.QtdGenre, w.Pars, w.Part,
		w.Email, w.Age, w.Gender, w.Reg, w.Sem, w.Org, w.Course, w.Language, w.Phone, w.CPF, w.ParsRead, w.DTBegin,
		w.HRBegin, w.ParID, w.SentID, cont, w.RawWord, w.Word,
		w.Resp, w.TBegin, w.TDig, w.TTot, w.TPar, w.TTest)
}
