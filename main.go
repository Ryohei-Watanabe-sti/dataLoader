package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func opencsv() {
	f, err := os.Open("sample1.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
	}

}

func main() {
	opencsv()
	//testtest
}
