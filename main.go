package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
    "fmt"
    "encoding/json"
)

type User1 struct {
    gorm.Model
	Name        string  `json:"name"`
    Image_url   string  `json:"image_url"`
    Email       string  `json:"email"`
    
}

func main() {
    dsn := "test:12345678@tcp(127.0.0.1:3306)/oh_test?character_set_database=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
		panic("failed to connect database")
	}
    up01(dsn, db)

    db.Debug().Create(&User1{Name: "aaa", Image_url: "image.png", Email: "ry-watanabe@sios.com"})
    
    var data1 User1

    db.Debug().First(&data1, 1)
    fmt.Println(data1)
    read_json, err := json.Marshal(data1)
    fmt.Println(string(read_json))

    db.Debug().Model(&User1{}).Where("id = ?", 1).Update("Name", "あああ")

    // db.Debug().Delete(&data1, 1)

    down01(dsn, db);

}

func up01(dsn string, db *gorm.DB) {
    db.AutoMigrate(User1{})
    fmt.Println("User1テーブルを作成")
}
func down01(dsn string, db *gorm.DB) {
    db.Migrator().DropTable(User1{})
    fmt.Println("User1テーブルを削除")
}
