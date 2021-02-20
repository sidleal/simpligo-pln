package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// ------------------------------------
var participantsToKeep = map[string]int{
	"AGATA JESSICA AVELAR DE OLIVEIRA":                 1,
	"Alex Verrua":                                      1,
	"Alexandre Brito Gomes":                            1,
	"Alexandre Lima Palles Rocha":                      1,
	"Alexandre Willian Biazão":                         1,
	"Alícia Mayumi Wauke":                              1,
	"alisson da silva margraf":                         1,
	"Alyne Bittencourt de Macedo Neves":                1,
	"Amanda Carvalho de Castro Santos":                 1,
	"Amanda Sayeg":                                     1,
	"Ana Beatriz Refundini Castellani":                 1,
	"Ana Carolina Chiciuc Vital":                       1,
	"Ana Carolina Estevão Cruz":                        1,
	"ana carolina leal pflanzer":                       1,
	"Ana Carolina Rodrigues":                           1,
	"Ana Caroline Lopes Rocha":                         1,
	"Ana Caroline Medeiros Brito":                      1,
	"ANA CLARA DOS SANTOS PEREIRA":                     1,
	"Ana Luisa Araujo de Souza":                        1,
	"Ana Paula da Silva Nascimento":                    1,
	"Ana Paula Martins de Meneses":                     1,
	"Ana Rachel Pereira Ribeiro":                       1,
	"André Zambroni Riedel":                            1,
	"Angela Hikari Miyagi":                             1,
	"Angela Maria Sariolli":                            1,
	"Antonio Rodrigues Rigolino":                       1,
	"Ariadne Augusto":                                  1,
	"ariane valeska mello de oliveira":                 1,
	"Ariel da Silva Pires":                             1,
	"Ariel Nascimento":                                 1,
	"Bárbara Sampaio Fernandes Tribino":                1,
	"Beatriz Torres Klem ":                             1,
	"Bernardo Marques Costa":                           1,
	"Bernardo Tampasco Cerri Costa":                    1,
	"Breno Cunha Queiroz":                              1,
	"Bruno Henrique Antunes":                           1,
	"BRUNO JOSE MACHADO DE CAMARGO":                    1,
	"Bruno Vinícius Melo da Cunha":                     1,
	"Caio Gomes Pariz":                                 1,
	"Cairo Alencar Zancanella":                         1,
	"Camila Coimbra Gomes Santos":                      1,
	"Camila Wielmowicki Uchoa":                         1,
	"Carlos Eduardo Barizon":                           1,
	"Carlos Eduardo Nishi":                             1,
	"Carlos Henrique de Barros Franzin":                1,
	"Carlos Leoni Rodrigues Siqueira Junior":           1,
	"carolina coentro da cunha":                        1,
	"Carolina Silva Couto":                             1,
	"Christian Bernard":                                1,
	"Cibele Andreciolli de Matos ":                     1,
	"ciro grossi falsarella":                           1,
	"Claudio Faria Marques":                            1,
	"Cristian koch winter":                             1,
	"Daniel da Silva Souto":                            1,
	"Daniele de Souza Aride":                           1,
	"Danielle Vitoria Santana Dias":                    1,
	"Danilo Bressiani Solek":                           1,
	"Davi Machado Calderari":                           1,
	"David Hallisson Rodrigues Souza da Costa":         1,
	"David Silvio Nucitelli Saquette":                  1,
	"Diógenes Silva Pedro":                             1,
	"Diovânia Maria Sabino da Fonseca Melhorance":      1,
	"Edson Alberto Pereira Rosa Junior":                1,
	"Eduardo Amaral":                                   1,
	"Eduardo Antonio da Cruz Esmagnoto":                1,
	"Eduardo Garcia de Gáspari Valdejão":               1,
	"Eduardo Higa":                                     1,
	"Erica Aparecida Novaes da Silva":                  1,
	"FABIO DA FONSECA MOREIRA":                         1,
	"Felipe Fahd Jensen Abud":                          1,
	"Felipe Monsores Franco":                           1,
	"Fernanda Alves Avelino":                           1,
	"Fernanda Venites Buzinaro":                        1,
	"Fernanda Wassano Daher":                           1,
	"Fernanda Yumi Ribeiro Mori":                       1,
	"Fernando Sales":                                   1,
	"Flávia Correia Lima Huber Costa":                  1,
	"Francisco Fábio Sales Melo":                       1,
	"Gabriel Alves Kuabara":                            1,
	"GABRIEL COZER LEAL":                               1,
	"Gabriel do Carmo Acioli de Oliveira":              1,
	"Gabriel Filipe Silva Santos":                      1,
	"Gabriel Freitas Ximenes de Vasconcelos":           1,
	"Gabriel Gomes da Costa":                           1,
	"Gabriel Luiz Scariot":                             1,
	"Gabriel Miashiro":                                 1,
	"Gabriel Miehe Machado":                            1,
	"Gabriel Pantano Signorini":                        1,
	"Gabriel Ribeiro Evangelista":                      1,
	"Gabriel Santos Martinelli":                        1,
	"Gabriel Silva Malta":                              1,
	"Gabriel Tomazini Marani":                          1,
	"Gabriela Cruz":                                    1,
	"Gabriela Giglio Gangoni":                          1,
	"Gabriela Rodrigues do Prado":                      1,
	"Gabriela Satie Faria Nishimi":                     1,
	"Giovana Meloni Craveiro":                          1,
	"Giovanna Pereira Campos":                          1,
	"Giovanni Barazetti Genero":                        1,
	"Giovanni Fantucci Soave":                          1,
	"GIULIA FERREIRA CHAGAS":                           1,
	"Guilherme Lima da Silva Vieira":                   1,
	"Guilherme Machado Rios":                           1,
	"Guilherme Pacheco":                                1,
	"Gustavo Moreira":                                  1,
	"Heitor Lélis Dias de Freitas":                     1,
	"Helena Lima Settecerze":                           1,
	"Hellen Aparecida Thiago Otaviano":                 1,
	"Henrique Alberto Mendonça":                        1,
	"Henrique Hiram Libutti Núñez":                     1,
	"Igor Mateus Queiroz Gato":                         1,
	"Igor Patrick Da Silva Molina":                     1,
	"indira costa cordeiro":                            1,
	"Ingrid Lima Pereira Peres":                        1,
	"Ingrid Santos de Oliveira":                        1,
	"Isabela De Freitas Zerbini":                       1,
	"Isabela Maria Freitas de Souza":                   1,
	"ISABELLE PEREIRA DA SILVA":                        1,
	"Jamilson Souza Soares":                            1,
	"Jean Domingos Lopes":                              1,
	"Jessica Siqueira Alvarenga ":                      1,
	"João Antônio Misson Milhorim":                     1,
	"João Carlos dos Santos ":                          1,
	"Joao Francisco Caprioli Barbosa Camargo de Pinho": 1,
	"João Guilherme Jarochinski Marinho":               1,
	"João Lucas Rodrigues Constantino":                 1,
	"João Marcos Cardoso da Silva":                     1,
	"JOAO PAULO DUTRA COSTA WADY":                      1,
	"João Pedro Cunha Guska":                           1,
	"João Pedro Fritsch Meinerz":                       1,
	"João Victor Sene Araújo":                          1,
	"João Vitor Lima Martins":                          1,
	"Johanna Kirchner":                                 1,
	"Joris Bianca da Silva":                            1,
	"José Vitor Montanger Ribeiro da Silva":            1,
	"Julia Correa Paulino":                             1,
	"Júlia Da Silva Lima":                              1,
	"Júlia Moser e Bittencourt":                        1,
	"Julia Yumi Ito":                                   1,
	"Julianne Cortez Tavares":                          1,
	"Karinna Adad de Miranda":                          1,
	"Karol Kniaczewski":                                1,
	"Karoline Fraga de Freitas":                        1,
	"Katarina Duarte Fernandes":                        1,
	"Kátia Selene de Melo":                             1,
	"KETHELEEN DOS SANTOS ANTUNES":                     1,
	"Laís Belmudo Gregorio":                            1,
	"LAIS ROLDAN AZEVEDO":                              1,
	"Leandro Moreira De Sousa":                         1,
	"Leonardo Ramos de Oliveira":                       1,
	"Leonardo Secco Alves":                             1,
	"Lethicia Roberta Barros Gonçalves":                1,
	"Liandra Ayumi Federizi Yoshida":                   1,
	"Lidia Gianne Souza da Rocha":                      1,
	"Lidia Morcelli Duarte":                            1,
	"Ligia Thomaz Vieira Leite":                        1,
	"Lívia Padovan Ricardo":                            1,
	"Lívia Peres Vieira":                               1,
	"Lourenço de Sallles Roselino":                     1,
	"lucas da silva alves":                             1,
	"Lucas Ferreira de Almeida":                        1,
	"Lucas Massao Fukusawa Dagnone":                    1,
	"Lucas Meneguetti Gaias":                           1,
	"Lucas Pilla Pimentel":                             1,
	"Lucas Roberto Valério Gonçalves":                  1,
	"Lucas Samuel Zanchet":                             1,
	"Luiz Gustavo Comarella":                           1,
	"Luiza Machado ":                                   1,
	"Maike de Oliveira Sala":                           1,
	"Maisa Mara Krupek":                                1,
	"Marcela Severo Mote da Silva":                     1,
	"Marcelle Belchior Cruz Ramos":                     1,
	"Marcio Lima Inácio":                               1,
	"Marcos Patrício Nogueira Filho":                   1,
	"Marcos Robert Lima Carneiro":                      1,
	"Marcus Vinícius Teixeira Huziwara":                1,
	"Maria Clara Castro da Silva":                      1,
	"Maria Eduarda De Bastiani":                        1,
	"Maria Eduarda Oliveira Avellar":                   1,
	"Maria Fernanda Pissolato":                         1,
	"Maria Haddock Lobo":                               1,
	"MARIA LUIZA GONÇALVES DE LIMA":                    1,
	"Maria Tainá Bernardino de Brito":                  1,
	"Mariana Fernandes Fonseca":                        1,
	"Marina Maia Reis":                                 1,
	"Marina Maximiano Ferreira de Souza":               1,
	"Martin Avila Buitron":                             1,
	"Mateus da Silva Páscoa":                           1,
	"Mateus Israel Silva":                              1,
	"Matheus Bermudes Viana":                           1,
	"Matheus Bueno Pereira":                            1,
	"Matheus Fernando Weis Cavalher":                   1,
	"Matheus Henrique Oliveira dos Santos":             1,
	"Matheus Ventura de Sousa":                         1,
	"Matheus Wabiszezewicz Baldacim":                   1,
	"Maynara Natalia Scoparo":                          1,
	"Melissa Motoki Nogueira":                          1,
	"Miguel de Mattos Gardini":                         1,
	"MIKAELA ALVES DE SOUZA":                           1,
	"Naiade Barros Gonçalves Marques":                  1,
	"Naiara Pereira Botezine":                          1,
	"Naími Moreira Nobre Leite":                        1,
	"Naioby Kelen Rosa":                                1,
	"Natan Henrique Sanches":                           1,
	"Natassia":                                         1,
	"Nathalia Anderson Passeri":                        1,
	"Nathália Fernanda Sczesny":                        1,
	"Nathalia Werneck ":                                1,
	"Nayara de Almeida Mello":                          1,
	"NEI CARLOS DOS SANTOS ROCHA":                      1,
	"Nicole Bertoncin Hermam":                          1,
	"Nicole Duailibi Barbosa":                          1,
	"Nicole Hafner Fagundes":                           1,
	"Nikolas Davies de Freitas":                        1,
	"Nivia de Souza Costa":                             1,
	"Noemi Mitsunaka":                                  1,
	"Paolla Tufolo Busto":                              1,
	"Patrick Galvão Neris":                             1,
	"Paula de Almeida Ribas Soares":                    1,
	"Paula Mendonça Dias":                              1,
	"Paulo Teixeira Vinhosa":                           1,
	"Paulo Victor Nogueira Rodrigues":                  1,
	"Pedro Afonso Perez Chagas":                        1,
	"Pedro Augusto Ribeiro Gomes":                      1,
	"Pedro Henrique dos Santos":                        1,
	"Pedro Henrique Raymundi":                          1,
	"Pedro Henrique Santana Costa":                     1,
	"Pedro Kenzo Muramatsu Carmo":                      1,
	"Pedro Liduino do Nascimento":                      1,
	"Pedro Lucas Linhares Paiva":                       1,
	"Petala Matutino Santos":                           1,
	"Queila Soares Candeias":                           1,
	"Rafael Corona":                                    1,
	"Rafael de Matos Pires":                            1,
	"Rafael Eizo Watabe":                               1,
	"Rafael Heidy Hirano":                              1,
	"Rafael Simões de Paula":                           1,
	"Rafaelli de Miranda Pereira":                      1,
	"Raíssa Torres Barreira":                           1,
	"Raphael Freires Pessoa":                           1,
	"Raquel Bauer Gomes da Silva":                      1,
	"Rayani de Lima Navega":                            1,
	"Roberta Benedetti Zanotto":                        1,
	"Robson Casciano Barros":                           1,
	"Rodrigo Augusto da Conceição Magalotti":           1,
	"Rodrigo Lopes Assaf":                              1,
	"Rosanne Pauzeiro Pousada":                         1,
	"Samuel Damiani Frigotto":                          1,
	"Samuel José Inacio":                               1,
	"Samuel Staskoviack Iglikovski":                    1,
	"Sandro Mariano Silva Filho":                       1,
	"Sara Melo de Queiroz":                             1,
	"Sofhia de Souza Gonçalves":                        1,
	"Stephanie Ferreira de Souza":                      1,
	"Taise Andrade dos Santos":                         1,
	"Talita de Oliveira Vargas Yamada":                 1,
	"Talita Miranda da Costa Mathias":                  1,
	"Teodoro Gomes de Moraes Colombo":                  1,
	"Thaisy Moraes Costa ":                             1,
	"Thales Nascimento Buzan":                          1,
	"Thayna Skerratt Pereira":                          1,
	"Thiago M Marchesan":                               1,
	"Thiago Onofre Pinheiro Rosa":                      1,
	"Thiago Takechi Ohno Bezerra":                      1,
	"Tiago Augusto Kirsch Andreis":                     1,
	"Tomás Gonçalves Grieser":                          1,
	"Valdinho Júnior Machado dos Santos":               1,
	"Vanessa Rodrigues da Silva":                       1,
	"Veronica Ruside higa":                             1,
	"Victor Cologni Seles":                             1,
	"Victor Henrique de Sa Silva":                      1,
	"Victor Paulo Cruz Lutes":                          1,
	"VICTORIA DE FATIMA DO CARMO ARRAES":               1,
	"Vinicius Bertuol":                                 1,
	"Vinícius de Oliveira Cardoso":                     1,
	"VINICIUS FELIPE FERRARI":                          1,
	"Vinicius Fernandes":                               1,
	"Vinicius Kamiya Svierk":                           1,
	"Vitor Azevedo Abou Mourad":                        1,
	"Vitor Felicio Salema":                             1,
	"Vitória Caetano de Oliveira":                      1,
	"Washington Luís":                                  1,
	"Wictor Dalbosco Silva":                            1,
	"Yann Amado Nunes Costa":                           1,
	"Yasmin Borges Tosta":                              1,
	"Yasmin dos Anjos de Deus Cardoso":                 1,
}

func main_clean() {

	// list := map[string]Word{}

	rawFiles := []string{
		"cloze_baJwn24BmRI7xu8F5xxd_data15.csv", //puc
		"cloze_bqJxn24BmRI7xu8Fqhz6_data15.csv", //usp
		"cloze_bKJwn24BmRI7xu8FSBzU_data15.csv", //ufc
		"cloze_UgXZ-28BaYrNtDuxSkMf_data15.csv", //utfpr
		"cloze_ogW4nXABaYrNtDuxrk04_data15.csv", //ufabc
		"cloze_YwXxiHABaYrNtDux7Uh8_data15.csv", //uerj
	}

	path := "/home/sidleal/sid/usp/cloze_exps/"

	exportDate := "2020_05_05_clean1"
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
			word.Word = cols[21]
			word.Resp = cols[22]
			word.TBegin, _ = strconv.Atoi(cols[23])
			word.TDig, _ = strconv.Atoi(cols[24])
			word.TTot, _ = strconv.Atoi(cols[25])
			word.TPar, _ = strconv.Atoi(cols[26])
			word.TTest, _ = strconv.Atoi(cols[27])

			//descarta participantes:
			discard := map[string]int{
				"55.019.616-x": 1, // não entendeu a tarefa - UTFPR
				"362726343":    1, // teste UFABC
				"098512262":    1, // não entendeu a tarefa  -  UERJ
			}
			if _, found := discard[word.Reg]; found {
				continue
			}

			if _, found := participantsToKeep[word.Part]; !found {
				continue
			}

			//ajuste... ufabc e uerj não tem o primeiro paragrafo
			if strings.Index(outFiles[k], "ufabc") > 0 || strings.Index(outFiles[k], "uerj") > 0 {
				word.ParID++
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
			fmt.Println(outFiles[k], "\t", key, "\t", val, "\t", len(val))
		}
	}

}
