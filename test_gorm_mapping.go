package main

import (
	"fmt"
	"github.com/BangNopall/paskihub-be/domain/dto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	db, _ := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=paskihub port=5432 sslmode=disable"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})
    if db == nil {
       fmt.Println("Failed to init DB dummy")
       return
    }
    stmt := &gorm.Statement{DB: db}
    stmt.Parse(&dto.ScoreboardItem{})
    for _, field := range stmt.Schema.Fields {
        fmt.Printf("Field: %s, Column: %s\n", field.Name, field.DBName)
    }
}
