package main

import (
	"fmt"

	"github.com/lach-killerpythons/alexandrite/BLUE"
	"github.com/lach-killerpythons/alexandrite/RED"
)

// eg: select website from sandstone where robots LIKE '%shopify%';
func BlueToRed_1dList(rdb RED.RedDB, db BLUE.DB, query string, keyname string) error {
	if rdb.GetP() == nil {
		return fmt.Errorf("RDB pointer error : not found")
	}
	if db.Table.Name == "" {
		return fmt.Errorf("psql must have a table selected")
	}

	//myQuery := fmt.Sprintf(`SELECT * FROM %s;`, tableName)
	rows, err := db.GetP().Query(query)
	if err != nil {
		return err
	}

	var nextLine string
	for rows.Next() {
		if err := rows.Scan(&nextLine); err != nil {
			fmt.Printf("Error scanning row: %v", err)
			continue // Skip to the next row if scan fails for this one
		}
		//fmt.Println(nextLine)
		rdb.List_Add(keyname, nextLine)
	}

	return nil
}

// new idea --> create API end point that is plug and play
// API baby

func main() {
	// jf := JADE.Open("db.json", "/JADE")
	// db, err := BLUE.DB_JadeConnect("local", jf)
	// if err != nil {
	// 	panic(err)
	// }
	// OpenTable_err := db.OpenTable("sandstone")
	// if OpenTable_err != nil {
	// 	fmt.Println(OpenTable_err)
	// }

	rdb, err := RED.NewRedDB("localhost", "", 0)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	rdb.List2Text("animals", "animals.txt")

	// animal_list, err := rdb.BetterListGet("animals")
	// for _, item := range animal_list {
	// 	fmt.Println(item)
	// }

	// add the URL list from psql to redis

	// myQuery := "select website from sandstone where robots LIKE '%shopify%';"
	// error_Blue2Red := BlueToRed_1dList(rdb, db, myQuery, "shopify_websites")
	// if error_Blue2Red != nil {
	// 	fmt.Println(error_Blue2Red)
	// }

	// rdb.List2_1wordset("shopify_websites", "web_set")

}
