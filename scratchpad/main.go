package main

import (
	"fmt"
	"os"

	"github.com/lach-killerpythons/alexandrite/BLUE"
	"github.com/lach-killerpythons/alexandrite/JADE"
)

func main() {
	cwd, _ := os.Getwd()
	fmt.Println(cwd)
	jade := JADE.JADE_FILE{"db.json", cwd}

	DB, err := BLUE.DB_JadeConnect("local", jade)
	if err != nil {
		fmt.Println(err)
	}
	DB.SELECT_ALL("blue_test")
}
