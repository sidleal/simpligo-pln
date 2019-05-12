package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/sidleal/simpligo-pln/tools/senter"
)

type MsgWSRanker struct {
	Authorization string `json:"auth"`
	Content       string `json:"content"`
	Options       string `json:"options"`
	RawResult     string `json:"raw_result"`
}

func RankerWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	go wsEcho(conn)

}

func wsEcho(conn *websocket.Conn) {
	for {
		m := MsgWSRanker{}

		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
			conn.Close()
			return
		}
		fmt.Printf("Got message: %#v\n", m)

		m.Authorization = ""

		content := m.Content
		options := m.Options

		if options == "unique" {

			resp, err := http.Post("http://"+mainServerIP+":8008/ranker", "text", bytes.NewReader([]byte(content)))
			if err != nil {
				m.RawResult = "Error: " + err.Error()
			} else {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					m.RawResult = "Error reading response: " + err.Error()
				} else {

					log.Println(string(body))

					m.RawResult = string(body)
				}
			}

		} else {
			parsed := senter.ParseText(content)
			m.RawResult = ""

			for _, p := range parsed.Paragraphs {
				for _, s := range p.Sentences {

					resp, err := http.Post("http://"+mainServerIP+":8008/ranker", "text", bytes.NewReader([]byte(s.Text)))
					if err != nil {
						m.RawResult = "Error: " + err.Error()
					} else {
						defer resp.Body.Close()
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							m.RawResult = "Error reading response: " + err.Error()
						} else {

							log.Println(string(body))
							resLines := strings.Split(string(body), "\n")
							result := fmt.Sprintf("%v --> %v\n\n", resLines[1][17:], resLines[5])

							m.RawResult += result
						}
					}

				}
			}

		}

		if err = conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
	}
}
