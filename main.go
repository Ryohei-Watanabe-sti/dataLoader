package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Csv1 string `json:"csv1"`
	Csv2 string `json:"csv2"`
	Dsn  string `json:"dsn"`
}

func readOption() string { //
	text1 := flag.String("t", "", "help message for t")
	text2 := flag.String("tablename", "", "help message for tablename")
	flag.Parse()
	var text string

	if *text1 == "" && *text2 == "" {
		panic("Enter the table name using \"-t\" or \"-tablename\".")
	} else if *text1 != "" && *text2 != "" {
		panic("Enter the table name one time.")
	} else if *text1 != "" {
		text = *text1
	} else if *text2 != "" {
		text = *text2
	}

	return text
}

func loggingSettings(filename string) {
	logfile, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(multiLogFile)
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

	log.Println(cnf.Csv1 + " and " + cnf.Csv2 + " was read successfully")

	return records1, records2
}

func createTable(tName string, columns1 [][]string, columns2 [][]string) {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	log.Println("successed to connext in dataLoader")

	if db.Migrator().HasTable(tName) == true {
		panic("failed to create " + tName + ", already exists.")
	}

	if len(columns1) != len(columns2) {
		panic("failed to create " + tName + " caused by CSV error.")
	}

	createTableQuery := `create table ` + tName + ` (id int); `
	db.Exec(createTableQuery)
	log.Println("successed to create table \"" + tName + "\" in dataLoader")

	var addColumnQuery string
	var cName string
	for i := 0; i < len(columns1[0]); i++ {
		cName = columns1[0][i]
		addColumnQuery = `alter table ` + tName + ` add column csv1_` + cName + ` varchar(255);`
		db.Exec(addColumnQuery)
	}
	for i := 0; i < len(columns2[0]); i++ {
		cName = columns2[0][i]
		addColumnQuery = `alter table ` + tName + ` add column csv2_` + cName + ` varchar(255);`
		db.Exec(addColumnQuery)
	}

	log.Println("start to insert records")

	var insertArr []string
	var insertQuery string
	var iString string
	all := float64(len(columns1) - 1)
	var done float64
	var done_ratio float64
	for i := 1; i < len(columns1); i++ {
		insertArr = append(columns1[i], columns2[i]...)
		iString = strconv.Itoa(i)
		insertQuery = "insert into " + tName + " values (" + iString + ","

		for j := 0; j < len(insertArr); j++ {
			insertQuery += "'" + insertArr[j] + "', "
		}
		insertQuery = insertQuery[:len(insertQuery)-2]
		insertQuery += ");"
		db.Exec(insertQuery)
		//進捗を表示
		done = float64(i)
		done_ratio = math.Floor(done / all * 100)
		log.Printf("%g%% done\n", done_ratio)
	}
	log.Println("end to insert records successfully")

}

func main() {

	loggingSettings("dataLoader.log")
	log.Println("dataLoader start")

	csv1, csv2 := getCSV()
	var tableName = readOption()
	createTable(tableName, csv1, csv2)
}
