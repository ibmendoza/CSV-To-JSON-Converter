package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	path := flag.String("path", "./data.csv", "Path of the file")
	flag.Parse()
	fileBytes, fileNPath := ReadCSV(path)
	SaveFile(fileBytes, fileNPath)
	fmt.Println(strings.Repeat("=", 10), "Done", strings.Repeat("=", 10))
}

// ReadCSV to read the content of CSV File
func ReadCSV(path *string) ([]byte, string) {
	csvFile, err := os.Open(*path)

	if err != nil {
		log.Fatal("The file is not found || wrong root")
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	content, _ := reader.ReadAll()

	if len(content) < 1 {
		log.Fatal("Something wrong, the file maybe empty or length of the lines are not the same")
	}

	headersArr := make([]string, 0)
	for _, headE := range content[0] {
		headersArr = append(headersArr, headE)
	}

	//Remove the header row
	content = content[1:]

	var buffer bytes.Buffer
	var nuBuffer *bytes.Buffer
	var nuBytes []byte
	buffer.WriteString("[")
	//for i, d := range content {
	for i, d := range content {
		buffer.WriteString("{")

		for j, y := range d {
			if len(y) > 0 {

				_, fErr := strconv.ParseFloat(y, 32)
				_, bErr := strconv.ParseBool(y)

				buffer.WriteString(`"` + headersArr[j] + `":`)

				if fErr == nil {
					buffer.WriteString(y)
				} else if bErr == nil {
					buffer.WriteString(strings.ToLower(y))
				} else {
					if len(y) > 0 {
						buffer.WriteString((`"` + y + `"`))
					}
				}

				if j < len(d) {
					buffer.WriteString(",")
				}
			}

		}

		nuBytes = bytes.TrimSuffix(buffer.Bytes(), []byte(","))

		//nuBuffer.WriteTo(os.Stdout)

		if i < len(content)-1 {
			buffer.WriteString(",")
		}

		//end of object of the array
		nuBuffer = bytes.NewBuffer(nuBytes)
		nuBuffer.WriteString("}")
	}

	//delete the last comma
	nuBytes = bytes.TrimSuffix(buffer.Bytes(), []byte(","))
	nuBuffer = bytes.NewBuffer(nuBytes)
	nuBuffer.WriteString(`]`)

	rawMessage := json.RawMessage(nuBuffer.String())
	x, err := json.MarshalIndent(rawMessage, "", "  ")

	if err != nil {
		fmt.Println(err)
	}
	newFileName := filepath.Base(*path)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"
	r := filepath.Dir(*path)
	return x, filepath.Join(r, newFileName)
}

// SaveFile Will Save the file, magic right?
func SaveFile(myFile []byte, path string) {
	if err := ioutil.WriteFile(path, myFile, os.FileMode(0644)); err != nil {
		panic(err)
	}
}
