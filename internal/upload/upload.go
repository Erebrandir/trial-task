package upload

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const DirectorySaved = "savedFiles"
const SizeBuffer = 1048576

func readFilePart(w http.ResponseWriter, r *http.Request) {
	file, heading, err := r.FormFile("file")

	if err != nil {
		return
	}
	defer file.Close()

	err = os.Mkdir(DirectorySaved, 0777)
	if err != nil {
		return
	}
	i := 0

	for {
		b := make([]byte, SizeBuffer)

		_, err := file.Read(b)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
			return
		}
		creatingFile(b, i, heading.Filename)
		i++
	}

	_, err = io.WriteString(w, "Файл сохранен")
	if err != nil {
		return
	}
}

func creatingFile(b []byte, i int, fileName string) {
	name := fileName + "-" + strconv.Itoa(i)

	err := ioutil.WriteFile(DirectorySaved+"/"+name, b, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("upload.html")

	err := t.Execute(w, "")
	if err != nil {
		return
	}

	readFilePart(w, r)
}
