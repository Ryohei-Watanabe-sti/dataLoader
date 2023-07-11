package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Csv1 string `json:"csv1"`
	Csv2 string `json:"csv2"`
	Dsn  string `json:"dsn"`
}

type Table1 struct {
	gorm.Model
	ID      int `json:"id"`
	Name_ja int `json:"name_ja"`
	Name_en int `json:"name_en"`
}

func loadConfig() (*Config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal("loadConfig os.Open err:", err)
		return nil, err
	}
	defer f.Close()

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	return &cfg, err
}

func getCSV() ([][]string, [][]string) {
	cnf, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(cnf.Csv1)
	if err != nil {
		log.Fatal(err)
	}
	r1 := csv.NewReader(f)

	g, err := os.Open(cnf.Csv2)
	if err != nil {
		log.Fatal(err)
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

func createTable() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	createTableQuery := `create table if not exists 
	table1(
		id int,
		name_ja varchar(10),
		name_en varchar(10)
	)
	CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci; `
	db.Exec(createTableQuery)
}

func main() {
	// csv1, csv2 := getCSV()
	// fmt.Println(csv1)
	// fmt.Println(csv2)
	createTable()
}
