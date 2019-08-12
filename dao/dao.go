package dao

import "fmt"

var dsn string = "root:123456@tcp(localhost:3306)/digimon?charset=utf8&loc=Local"

func Init() {
	fmt.Println("dao connection init successful")
}
