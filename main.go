package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Csv1 string `json:"csv1"`
	Csv2 string `json:"csv2"`
	Dsn  string `json:"dsn"`
}

type Table1 struct {
	ID        uint `gorm:"primaryKey"`
	Name_ja   string
	Name_en   string
	CreatedAt time.Time
	UpdatedAt time.Time
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

func createTable(tName string, columns1 []string, columns2 []string) {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	createTableQuery := `create table if not exists ` + tName + ` (
		id int,
		name_ja varchar(10),
		name_en varchar(10)
	); `
	db.Exec(createTableQuery)
}

func readOption() string {
	var text = flag.String("t", "tablename", "help message for t")
	flag.Parse()
	return *text
}

func main() {
	csv1, csv2 := getCSV()
	var tableName = readOption()
	createTable(tableName, csv1[0], csv2[0])

}
