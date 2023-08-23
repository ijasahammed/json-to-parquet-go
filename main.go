package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	datatype "parquet-go/type"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
)

func main() {

	inputFileName := ""
	outputFileName := "output.parquet"

	if len(os.Args) > 1 {
		fmt.Println(os.Args, len(os.Args))
		inputFileName = os.Args[1]
		if len(os.Args) > 2 {
			outputFileName = os.Args[2]
		}
	} else {
		log.Println("Input file name missing")
		return
	}

	fw, err := local.NewLocalFileWriter(outputFileName)
	if err != nil {
		log.Println("Can't create file", err)
		return
	}

	file, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var data []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatal(err)
	}

	metadata := datatype.BuildMedadata(data[0])

	fmt.Println(metadata)

	pw, err := writer.NewJSONWriter(metadata, fw, 4)
	if err != nil {
		log.Println("Can't create json writer : ", err)
		return
	}
	for _, v := range data {
		json, _ := json.Marshal(v)
		if err = pw.Write(string(json)); err != nil {
			log.Println("Write error", err)
			break
		}
	}
	if err = pw.WriteStop(); err != nil {
		log.Println("WriteStop error", err)
	}
	log.Println("Write Finished")
	fw.Close()

}
