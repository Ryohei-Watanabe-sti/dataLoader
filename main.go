package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func getCSV() ([][]string, [][]string) {
	f, err1 := os.Open("sample1.csv")
	if err1 != nil {
		log.Fatal(err1)
	}
	r1 := csv.NewReader(f)

	g, err2 := os.Open("sample2.csv")
	if err2 != nil {
		log.Fatal(err2)
	}
	r2 := csv.NewReader(g)

	records1, err := r1.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	records2, err := r2.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records1, records2
}

func main() {
	csv1, csv2 := getCSV()
	fmt.Println(csv1, csv2)
}
