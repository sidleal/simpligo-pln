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
	ParsRead string // Parágrafos Lidos
	DTBegin  string // Data Início
	HRBegin  string // Hora Início
	ParID    int    // Parágrafo
	SentID   int    // Sentença
	WordID   int    // Índice Palavra
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

func main() {

	// list := map[string]Word{}

	rawFiles := []string{
		"cloze_baJwn24BmRI7xu8F5xxd_data8.csv",
		"cloze_bqJxn24BmRI7xu8Fqhz6_data8.csv",
		"cloze_bKJwn24BmRI7xu8FSBzU_data8.csv",
		"cloze_UgXZ-28BaYrNtDuxSkMf_data8.csv",
		"cloze_ogW4nXABaYrNtDuxrk04_data8.csv",
		"cloze_YwXxiHABaYrNtDux7Uh8_data8.csv",
	}

	path := "/home/sidleal/sid/usp/cloze_exps/"

	exportDate := "2020_04_05"
	outFiles := []string{
		"cloze_puc_" + exportDate + ".csv",
		"cloze_usp_" + exportDate + ".csv",
		"cloze_ufc_" + exportDate + ".csv",
		"cloze_utfpr_" + exportDate + ".csv",
		"cloze_ufabc_" + exportDate + ".csv",
		"cloze_uerj_" + exportDate + ".csv",
	}

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

			// line = strings.ReplaceAll(line, "Atualmente, estou", "Atualmente estou")
			// line = strings.ReplaceAll(line, ", ", ",")
			// line = strings.ReplaceAll(line, " ,", ",")

			cols := strings.Split(line, ",")
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
			word.ParsRead = cols[10]
			word.DTBegin = cols[11]
			word.HRBegin = cols[12]
			word.ParID, _ = strconv.Atoi(cols[13])
			word.SentID, _ = strconv.Atoi(cols[14])
			word.WordID, _ = strconv.Atoi(cols[15])
			word.Word = cols[16]
			word.Resp = cols[17]
			word.TBegin, _ = strconv.Atoi(cols[18])
			word.TDig, _ = strconv.Atoi(cols[19])
			word.TTot, _ = strconv.Atoi(cols[20])
			word.TPar, _ = strconv.Atoi(cols[21])
			word.TTest, _ = strconv.Atoi(cols[22])

			//descarta participantes:
			discard := map[string]int{
				"55.019.616-x": 1, // não entendeu a tarefa - UTFPR
				"362726343":    1, // teste UFABC
				"098512262":    1, // não entendeu a tarefa  -  UERJ
			}
			if _, found := discard[word.Reg]; found {
				continue
			}

			//ajuste... ufabc e uerj não tem o primeiro paragrafo
			if strings.Index(outFiles[k], "ufabc") > 0 || strings.Index(outFiles[k], "uerj") > 0 {
				word.ParID++
			}

			//another cleaning
			word.Resp = strings.ReplaceAll(word.Resp, `"`, "")

			wordList = append(wordList, word)

			partKey := fmt.Sprintf("%s_%s", word.Part, word.Email)
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
			fmt.Println(outFiles[k], "\t", key, "\t", val, "\t", len(val))
		}
	}

}

func formatLine(w Word, cont int) string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
		w.Code, w.Name, w.QtdGenre, w.Pars, w.Part,
		w.Email, w.Age, w.Gender, w.Reg, w.Sem, w.ParsRead, w.DTBegin,
		w.HRBegin, w.ParID, w.SentID, cont, w.Word,
		w.Resp, w.TBegin, w.TDig, w.TTot, w.TPar, w.TTest)
}
