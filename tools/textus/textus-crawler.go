package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// ------------------------------------

var path = "/home/sidleal/sid/usp/textus"

func mainx() {

	log.Println(path)

	token := "1F7F9ACD43D13040387BCCF5BBA23DB2"
	// id := 5571
	for id := 3981; id <= 5571; id++ {
		time.Sleep(2 * time.Second)
		downloadTextusFile(id, token, "texto")
		time.Sleep(1 * time.Second)
		downloadTextusFile(id, token, "digital")
	}

}

func downloadTextusFile(id int, token string, fileType string) {
	cookie := http.Cookie{}
	cookie.Name = "JSESSIONID"
	cookie.Value = token

	url := fmt.Sprintf("https://www.convenios.grupogbd.com/redacoes/Redacao.download?id=%v&tipo=%s", id, fileType)
	log.Println(url)

	req, _ := http.NewRequest("GET", url, nil)

	req.AddCookie(&cookie)
	log.Println(req.Header)

	timeout := time.Duration(300 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(fmt.Sprintf("Error: %v", err))
	}

	log.Println(resp.Header)
	content := resp.Header["Content-Disposition"]
	if len(content) == 0 {
		log.Println("Sem id: ", id)
		return
	}
	contentTokens := strings.Split(content[0], `"`)
	fileName := contentTokens[1]

	defer resp.Body.Close()

	f, err := os.Create(fmt.Sprintf("%s/%04d-%s", path, id, fileName))
	if err != nil {
		log.Println("ERRO", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Error: %v", err))
	}

}

// func main_testing() {

// 	log.Println(path)

// 	id := 2

// 	timeout := time.Duration(300 * time.Second)
// 	client := http.Client{
// 		Timeout: timeout,
// 	}

// 	resp, err := client.PostForm("https://www.convenios.grupogbd.com/redacoes/efetuaLogin",
// 		url.Values{"email": {"sandra@icmc.usp.br"}, "senha": {"1997"}})
// 	if err != nil {
// 		log.Println(fmt.Sprintf("Error: %v", err))
// 	}

// 	log.Println("login header", resp.Header)
// 	log.Println("cookies", resp.Cookies())
// 	log.Println("req header", resp.Request.Header)
// 	log.Println("req cook", resp.Request.Cookies())

// 	cookies := resp.Cookies()

// 	for _, cookie := range cookies {
// 		fmt.Println("Found a cookie named:", cookie.Name)
// 	}

// 	rdrBody := io.Reader(resp.Body)
// 	// rdrBody = charmap.ISO8859_1.NewDecoder().Reader(rdrBody)

// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(rdrBody)
// 	if err != nil {
// 		log.Println(fmt.Sprintf("Error: %v", err))
// 	}

// 	bodyStr := string(body)
// 	log.Println(bodyStr)

// 	// callTextus(2, client, cookies)

// 	url := fmt.Sprintf("https://www.convenios.grupogbd.com/redacoes/Redacao.download?id=%v&tipo=texto", id)
// 	log.Println(url)

// 	req, _ := http.NewRequest("GET", url, nil)
// 	for i := range cookies {
// 		req.AddCookie(cookies[i])
// 	}
// 	cookie := http.Cookie{}
// 	cookie.Name = "JSESSIONID"
// 	cookie.Value = "7FBC6F8DA1AFC890BE79B79D7FFE301A"

// 	req.AddCookie(&cookie)
// 	log.Println(req.Header)

// 	resp, err = client.Do(req)
// 	if err != nil {
// 		log.Println(fmt.Sprintf("Error: %v", err))
// 	}

// 	log.Println(resp.Header)

// 	defer resp.Body.Close()

// 	f, err := os.Create(fmt.Sprintf("%s/%v-arquivo.doc", path, id))
// 	if err != nil {
// 		log.Println("ERRO", err)
// 	}
// 	defer f.Close()

// 	_, err = io.Copy(f, resp.Body)
// 	if err != nil {
// 		log.Println(fmt.Sprintf("Error: %v", err))
// 	}

// }

// func callTextus(id int, client http.Client, cookies []*http.Cookie) {

// 	time.Sleep(1 * time.Second)

// 	url := fmt.Sprintf("https://www.convenios.grupogbd.com/redacoes/Redacao.download?id=%v&tipo=texto", id)
// 	log.Println(url)

// 	req, _ := http.NewRequest("GET", url, nil)
// 	for i := range cookies {
// 		req.AddCookie(cookies[i])
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Println(fmt.Sprintf("Error: %v", err))
// 	}

// 	log.Println(resp.Header)

// 	defer resp.Body.Close()

// 	f, err := os.Create(fmt.Sprintf("%s/%v-arquivo.doc", path, id))
// 	if err != nil {
// 		log.Println("ERRO", err)
// 	}
// 	defer f.Close()

// 	_, err = io.Copy(f, resp.Body)
// 	if err != nil {
// 		log.Println(fmt.Sprintf("Error: %v", err))
// 	}

// }
