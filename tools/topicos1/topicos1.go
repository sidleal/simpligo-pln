package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type AnnotationData struct {
	Completions []Completion       `json:"completions"`
	ID          int                `json:"id"`
	TaskPath    string             `json:"task_path"`
	Text        AnnotationDataText `json:"data"`
}

type AnnotationDataText struct {
	Text string `json:"text"`
}

type Completion struct {
	Results []ResultItem `json:"result"`
}

type ResultItem struct {
	FromName string          `json:"from_name"`
	ToName   string          `json:"to_name"`
	ID       string          `json:"id"`
	Source   string          `json:"source"`
	Type     string          `json:"type"`
	Value    ResultItemValue `json:"value"`
}

type ResultItemValue struct {
	Start  int      `json:"start"`
	End    int      `json:"end"`
	Text   string   `json:"text"`
	Labels []string `json:"labels"`
}

const rootDir string = "/home/sidleal/sid/usp/TopicosPLN/trab1/res"

func main() {

	log.Println("starting")

	mapaGeral := map[string]map[string]int{}

	annDirs, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Println("Erro", err)
	}

	for _, d := range annDirs {
		annotator := d.Name()
		log.Println("Anotador: ", annotator)

		taskDirs, err := ioutil.ReadDir(rootDir + "/" + annotator)
		if err != nil {
			log.Println("Erro", err)
		}

		for _, d2 := range taskDirs {
			task := d2.Name()
			log.Println("Task: ", task)

			files, err := ioutil.ReadDir(rootDir + "/" + annotator + "/" + task)
			if err != nil {
				log.Println("Erro", err)
			}

			for _, f := range files {
				raw := readFile(rootDir + "/" + annotator + "/" + task + "/" + f.Name())
				// log.Println(raw)
				log.Println("Texto: ", task+"/"+f.Name())

				var data AnnotationData
				err := json.Unmarshal([]byte(raw), &data)
				if err != nil {
					log.Printf("Erro ao tratar json: %v", err)
				}

				log.Println("--------> ", len(data.Completions[0].Results), "Sentenças")
				if _, found := mapaGeral[task+"/"+f.Name()]; !found {
					mapaGeral[task+"/"+f.Name()] = map[string]int{}
				}
				mapaGeral[task+"/"+f.Name()][annotator] = len(data.Completions[0].Results)
				// for _, res := range data.Completions[0].Results {
				// 	log.Println("Início:", res.Value.Start, "Fim:", res.Value.End)
				// }

			}

		}
	}

	for k, v := range mapaGeral {
		//log.Println(k, v)
		log.Println(k, v["Marcio"], v["Sid"], v["Ana"], v["Denis"])
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
