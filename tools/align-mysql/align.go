package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "railsanotador:simplific@ndoP@(172.18.0.4)/anotador")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query(
		`select id, title from productions where project_id in (1,2) and status <> 'REMOVE' order by id`,
	)

	for rows.Next() {
		var id int64
		var title string
		err = rows.Scan(&id, &title)
		if err != nil {
			log.Println("Scan failed:", err.Error())
		}

		log.Println("---------------------------------")
		log.Println(id, title)
		treatProduction(db, id)

	}
	rows.Close()

	log.Println("done.")
}

func treatProduction(db *sql.DB, prodId int64) {
	rows, err := db.Query(
		fmt.Sprintf(`select id, tipo from textos where production_id = %d`, prodId),
	)
	mapSimp := make(map[string]int64)

	var id int64
	var tipo string
	for rows.Next() {
		err = rows.Scan(&id, &tipo)
		if err != nil {
			log.Println("Scan failed:", err.Error())
		}
		mapSimp[tipo] = id
	}

	log.Println(mapSimp)

	align1 := alignText(db, mapSimp["ORIGINAL"], mapSimp["NATURAL"])
	align2 := alignText(db, mapSimp["NATURAL"], mapSimp["VIOLENTO"])

	f, err := os.OpenFile("/home/sidleal/usp/align-11.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERRO", err)
	}

	defer f.Close()

	for i, align := range align1 {
		if align[0] != align[1] && len(strings.Split(align[0], " ")) > 1 {
			n, err := f.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", prodId, "ORI->NAT", i, align[0], align[1]))
			if err != nil {
				log.Println("ERRO", err)
			}
			fmt.Printf("wrote %d bytes\n", n)
		}
		log.Println(i, align[0])
		log.Println(i, align[1])
	}

	for i, align := range align2 {
		if align[0] != align[1] && len(strings.Split(align[0], " ")) > 1 {
			n, err := f.WriteString(fmt.Sprintf("%v\t%v\t%v\t%v\t%v\n", prodId, "NAT->STR", i, align[0], align[1]))
			if err != nil {
				log.Println("ERRO", err)
			}
			fmt.Printf("wrote %d bytes\n", n)
		}
		log.Println(i, align[0])
		log.Println(i, align[1])
	}

	rows.Close()
}

func alignText(db *sql.DB, idFrom int64, idTo int64) [][]string {

	ret := [][]string{}

	rows, err := db.Query(
		fmt.Sprintf(`select sentenceA, sentenceB from alignments where texto_id = %d and textoB = %d`, idFrom, idTo),
	)

	for rows.Next() {
		var sentA int64
		var sentB int64
		err = rows.Scan(&sentA, &sentB)
		if err != nil {
			log.Println("Scan failed:", err.Error())
		}

		// log.Println("---", sentA, sentB)

		textSentA := getSentence(db, sentA)
		textSentB := getSentence(db, sentB)

		ret = append(ret, []string{textSentA, textSentB})
		// log.Println(textSentA)
		// log.Println(textSentB)

	}
	rows.Close()

	return ret
}

func getSentence(db *sql.DB, id int64) string {
	rows, err := db.Query(
		fmt.Sprintf(`select word from words where sentence_id = %d order by id`, id),
	)

	sent := ""
	for rows.Next() {
		var word string
		err = rows.Scan(&word)
		if err != nil {
			log.Println("Scan failed:", err.Error())
		}
		word = strings.TrimSpace(word)
		sent += strings.TrimSpace(word) + " "

		// log.Println("---", word)

	}
	rows.Close()

	sent = strings.Replace(sent, " $.", ".", -1)
	sent = strings.Replace(sent, " $,", ",", -1)
	sent = strings.Replace(sent, " $)", ")", -1)
	sent = strings.Replace(sent, " $:", ":", -1)
	sent = strings.Replace(sent, " $;", ";", -1)
	sent = strings.Replace(sent, " $?", "?", -1)
	sent = strings.Replace(sent, " $!", "!", -1)
	sent = strings.Replace(sent, " $%", "%", -1)
	sent = strings.Replace(sent, "$\"", "\"", -1)
	sent = strings.Replace(sent, "$'", "'", -1)
	sent = strings.Replace(sent, "$( ", "(", -1)
	sent = strings.Replace(sent, "$--", "--", -1)

	regEx := regexp.MustCompile(`"\s(.*?)\s*?"`)
	sent = regEx.ReplaceAllString(sent, `"$1"`)

	return sent
}

/*
 select t.tipo, f.description, count(*) qtde from textos t
 inner join sentences s on s.texto_id = t.id
 inner join words w on w.sentence_id = s.id
 inner join features f on f.word_id = w.id
 where t.production_id in (
     select id from productions where project_id in (1,2) and status <> 'REMOVE' order by id
 )
 and f.tipo = 'SINTATICO'
 group by t.tipo, f.description
 order by t.tipo, qtde desc


 create view sentences_size as (
 select t.tipo, t.id, count(*) qtde from textos t
 inner join sentences s on s.texto_id = t.id
 where t.production_id in (
     select id from productions where project_id in (1,2) and status <> 'REMOVE' order by id
 )
 group by t.tipo, t.id
)
 select tipo, sum(qtde), min(qtde), max(qtde), avg(qtde)
 from sentences_size
 group by tipo



  create view tokens_size as
( select t.tipo, s.id, count(*) qtde from textos t
 inner join sentences s on s.texto_id = t.id
 inner join words w on w.sentence_id = s.id
 where t.production_id in (
     select id from productions where project_id in (1,2) and status <> 'REMOVE' order by id
 )
 group by t.tipo, s.id
 )
 select tipo, sum(qtde), min(qtde), max(qtde), avg(qtde)
 from tokens_size
 group by tipo


*/
