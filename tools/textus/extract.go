package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func main() {

	log.Println("starting")

	// fout, err := os.OpenFile("/home/sidleal/sid/usp/textus/tudo.tsv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Println("ERRO", err)
	// }
	// defer fout.Close()

	path := "/home/sidleal/sid/usp/textus"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Erro", err)
	}

	for _, f := range files {
		fileName := f.Name()
		log.Println("---------------------------------------------------------------")
		log.Println(fileName)

		if !strings.HasSuffix(fileName, "doc") {
			continue
		}

		// if i > 10 {
		// 	break
		// }

		log.Println("/bin/bash", "-c", "cd "+path+"/out/ && libreoffice --headless --convert-to \"txt:Text (encoded):UTF8\" "+path+"/"+fileName)

		cmd := exec.Command("/bin/bash", "-c", "cd "+path+"/out/ && libreoffice --headless --convert-to \"txt:Text (encoded):UTF8\" "+path+"/"+fileName)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("cmd.Run() failed with " + err.Error())
		}
		fmt.Printf("combined out:\n%s\n", string(out))
	}

}
