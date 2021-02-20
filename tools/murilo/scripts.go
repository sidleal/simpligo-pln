package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main_script() {

	log.Println("starting")

	fout, err := os.OpenFile("/home/sidleal/Downloads/murilo/tudo.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}
	defer fout.Close()

	path := "/home/sidleal/Downloads/murilo/analisando_curation_contagem/curation_v2"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	for _, f := range files {
		fileName := f.Name()
		log.Println("---------------------------------------------------------------")
		log.Println(fileName)

		files2, err := ioutil.ReadDir(path + "/" + fileName)
		if err != nil {
			log.Println("Erro", err)
		}

		for _, f2 := range files2 {
			fileName2 := f2.Name()
			log.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n--->", fileName2)

			files3, err := ioutil.ReadDir(path + "/" + fileName + "/" + fileName2)
			if err != nil {
				log.Println("Erro", err)
			}

			if fileName2 == "A042 - F.M._v4_clear.txt" ||
				fileName2 == "A089 – Y.S.P_v4_clear.txt" ||
				fileName2 == "A097 – K.V.S.S_v4_clear.txt" ||
				fileName2 == "A285_v4_clear.txt" {
				continue
			}
			i := 0
			for _, f3 := range files3 {
				fileName3 := f3.Name()
				log.Println("--------->", i, fileName3)

				buffer2 := readFile(path + "/" + fileName + "/" + fileName2 + "/" + fileName3)
				lines2 := strings.Split(buffer2, "\n")
				firstPos, adjustPos := 0, 0
				posIni, posEnd := 0, 0
				newBlock := false
				lastTag := "_"
				chunk := ""
				lastLength := 0
				for _, line := range lines2 {
					tokens := strings.Split(line, "\t")
					if tokens[0] == "" {
						continue
					}
					if strings.HasPrefix(tokens[0], "#Text=") {
						if lastTag != "_" {
							posEnd = posEnd + adjustPos
							runes := []rune(chunk)
							log.Println("=end->", "      posini:", posIni, "     posend:", posEnd, "   lenlst:", lastLength, "  -  first:", firstPos, "  =   adj:", adjustPos, "   len:", len([]rune(chunk)), "    tag:", lastTag, "  tokens:", tokens)
							// log.Println("---- end ----/", lastTag, posEnd, adjustPos)
							tagEnd := "</" + lastTag + ">"
							chunk = string(runes[0:posEnd]) + tagEnd + string(runes[posEnd:])
							// log.Println("aaaaaaaaaaaaaaaa", tagEnd, string(runes[posEnd:]))
							adjustPos += len(tagEnd)
							log.Println("--+--end+++--->", chunk)
							lastTag = "_" // reseta pra não fechar 2 vezes
						}
						lastLength = len([]rune(chunk))
						chunk += tokens[0][6:] + "  "
						if strings.Index(chunk, "\\r") > 0 {
							chunk = strings.ReplaceAll(chunk, "\\r", "")
						}
						log.Println("--+-->", chunk)
						newBlock = true
						continue
					}
					if strings.HasPrefix(tokens[0], "#") {
						continue
					}

					positions := strings.Split(tokens[1], "-")
					positionIni, _ := strconv.Atoi(positions[0])
					positionEnd, _ := strconv.Atoi(positions[1])

					// log.Println("xxxxxxxxxxxxxxxxxxxxxxxxx", positions, "first", firstPos, "adj", adjustPos, "len last", lastLength, "len", len([]rune(chunk)))

					if newBlock {
						// adjustPos--
						firstPos = positionIni
						if firstPos > 2 {
							adjustPos = lastLength - firstPos
						} else if firstPos > 0 {
							adjustPos = firstPos * -1
						}
						newBlock = false
					}

					tag := tokens[3]
					tag = strings.ReplaceAll(tag, "\\", "")

					if lastTag != tag && lastTag != "_" {
						runes := []rune(chunk)
						posEnd = posEnd + adjustPos
						log.Println("=end=>", "      posini:", posIni, "     posend:", posEnd, "   lenlst:", lastLength, "  -  first:", firstPos, "  =   adj:", adjustPos, "   len:", len([]rune(chunk)), "    tag:", tag, "  tokens:", tokens)
						// log.Println("---- end ----", lastTag, posEnd, positionEnd, adjustPos)
						tagEnd := "</" + lastTag + ">"
						chunk = string(runes[0:posEnd]) + tagEnd + string(runes[posEnd:])
						adjustPos += len(tagEnd)
						log.Println("--+--end+---->", chunk)
					}

					if lastTag != tag && tag != "_" {
						runes := []rune(chunk)
						posIni = positionIni + adjustPos
						log.Println("=ini=>", "      posini:", posIni, "     posend:", posEnd, "   lenlst:", lastLength, "  -  first:", firstPos, "  =   adj:", adjustPos, "   len:", len([]rune(chunk)), "    tag:", tag, "  tokens:", tokens)
						tagIni := "<" + tag + ">"
						chunk = string(runes[0:posIni]) + tagIni + string(runes[posIni:])
						adjustPos += len(tagIni)
						log.Println("--+--ini+---->", chunk)
					}

					posEnd = positionEnd
					lastTag = tag

					log.Println("=====>", "      posini:", posIni, "     posend:", posEnd, "   lenlst:", lastLength, "  -  first:", firstPos, "  =   adj:", adjustPos, "   len:", len([]rune(chunk)), "    tag:", tag, "  tokens:", tokens)

				}

				if lastTag != "_" {
					posEnd = posEnd + adjustPos
					runes := []rune(chunk)
					log.Println("=end+>", "      posini:", posIni, "     posend:", posEnd, "   lenlst:", lastLength, "  -  first:", firstPos, "  =   adj:", adjustPos, "   len:", len([]rune(chunk)), "    tag:", lastTag)
					// log.Println("end", lastTag, posEnd, adjustPos)
					tagEnd := "</" + lastTag + ">"
					chunk = string(runes[0:posEnd]) + tagEnd + string(runes[posEnd:])
				}

				fout.WriteString(fileName + "/" + fileName2 + "/" + fileName3 + "\n")
				fout.WriteString(chunk + "\n\n")
			}

		}

	}

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
