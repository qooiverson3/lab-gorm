package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"wesley.com/go-api-gorm/pkg/repository"
)

type User struct {
	Name string `gorm:"primary_key" json:"name"`
	Age  int8   `json:"age"`
}

func main() {
	db, err := repository.OpenGORM(
		"root",
    os.GetEnv("pass"),
		"127.0.0.1:3306",
		"iac",
	)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Table("user").AutoMigrate(&User{}); err != nil {
		log.Fatal(err)

	}

	user := []User{
		User{Name: "wesley", Age: 33},
		User{Name: "user2", Age: 34},
	}

	result, err := Insert(db, user)
	if err != nil {
		log.Fatal(err)
	}

	log.Info(fmt.Sprintf("impact amount: %v", result))
}

func Insert(db *gorm.DB, user []User) (int64, error) {
	tx := db.Begin()                         // 開始一個事務
	result := tx.Table("user").Create(&user) // 在事務中執行插入操作
	if result.Error != nil {
		tx.Rollback() // 發生錯誤時回滾事務
		return 0, result.Error
	}
	tx.Commit() // 提交事務

	return result.RowsAffected, nil
}
