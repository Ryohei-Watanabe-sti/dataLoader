package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
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

	if db.Migrator().HasTable(tName) == true {
		panic("failed to create " + tName)
	}

	createTableQuery := `create table ` + tName + ` (id int); `
	db.Exec(createTableQuery)

	var addColumnQuery string
	var cName string
	for i := 0; i < len(columns1); i++ {
		cName = columns1[i]
		addColumnQuery = `alter table ` + tName + ` add column csv1_` + cName + ` varchar(255);`
		db.Exec(addColumnQuery)
	}
	for i := 0; i < len(columns2); i++ {
		cName = columns2[i]
		addColumnQuery = `alter table ` + tName + ` add column csv2_` + cName + ` varchar(255);`
		db.Exec(addColumnQuery)
	}

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
