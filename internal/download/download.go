package download

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"trial_task/internal/upload"
)

const directoryDownload = "downloadFiles"

func checkFiles(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func checkPartFile(partFile []byte) []byte {
	size := upload.SizeBuffer
	for i := 0; i < size; i++ {
		if partFile[i] == 0 {
			size = i
		}
	}

	if size == upload.SizeBuffer {
		return partFile
	} else {
		tempPartFile := make([]byte, size)
		for j := 0; j < size; j++ {
			tempPartFile[j] = partFile[j]
		}
		return tempPartFile
	}
}

func handDownload(w http.ResponseWriter, r *http.Request, fileName string) {
	var completeFile bytes.Buffer
	i := 0

	for {
		partFile, err := os.ReadFile(upload.DirectorySaved + "/" + fileName + "-" + strconv.Itoa(i))
		if err != nil {
			if strings.Contains(err.Error(), "нет такого файла или каталога") {
				break
			}
			log.Fatal(err)
		}
		partFile = checkPartFile(partFile)
		completeFile.Write(partFile)
		i++
	}

	if checkFiles(directoryDownload + "/" + fileName) {
		_, err := io.WriteString(w, "Ошибка: файл уже существует.")
		if err != nil {
			return
		}
		return
	}

	fileName = r.FormValue("text")
	if len(fileName) == 0 {
		return
	}

	err := os.Mkdir(directoryDownload, 0777)
	if err != nil {
		return
	}

	file, err := os.Create(directoryDownload + "/" + fileName)
	if err != nil {
		fmt.Println("Не удалось создать файл:", err)
		os.Exit(1)
	}

	defer file.Close()
	_, err = file.WriteString(completeFile.String())
	if err != nil {
		return
	}

	_, err = io.WriteString(w, "Файл собран и сохранен.")
	if err != nil {
		return
	}
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("download.html")

	err := t.Execute(w, "")
	if err != nil {
		return
	}

	fileName := r.FormValue("text")
	if len(fileName) == 0 {
		return
	}

	tempFile, err := os.Open(upload.DirectorySaved + "/" + string(fileName))
	if err != nil {
		log.Print("Файл не найден")
		_, err := io.WriteString(w, "Файл не найден")
		if err != nil {
			return
		}
		w.WriteHeader(404)
		return
	}

	err = tempFile.Close()
	if err != nil {
		return
	}
	handDownload(w, r, fileName)
}
