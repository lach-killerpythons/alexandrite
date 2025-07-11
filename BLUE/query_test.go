package BLUE

import (
	"fmt"
	"testing"
)

// run a single test
// `go test -run ^TestDt$`

const testTable = "blue_test"

func InsertBlueTest1(DBx DB) error {

	test2 := []interface{}{
		1,
		"yellow",
	}

	col_names := []string{"id", "placeholder"}

	DBx.SetTable(testTable, col_names)
	fmt.Println(DBx.Table.Name, "jeez")
	INSERT_err := DBx.INSERT_WILD(test2)
	if INSERT_err != nil {
		return INSERT_err
	}
	return nil

}

func InsertBlueTest2(DBx DB) error {

	test69 := []interface{}{
		69,
		"GREEN",
	}
	INSERT_err := DBx.INSERT_WILD(test69)
	if INSERT_err != nil {
		return INSERT_err
	}

	return nil
}

func TestInsert(t *testing.T) {
	DBx, err := DB_Connect("local")
	if err != nil {
		t.Errorf("DB_Connect error")
	}
	insert_err := DBx.OpenTable(testTable)
	if insert_err != nil {
		t.Errorf("error opening table %s", insert_err)
	}
	insert_err = InsertBlueTest2(DBx)
	if insert_err != nil {
		t.Errorf("InsertBlueTest2 Error %s", insert_err)
	}

}

func TestQuery(t *testing.T) {
	// query each DB for table blue_test
	// select first row
	// insert
	testDBs := []string{"pi", "mac", "local"}
	desiredOutput := []string{"1", "Test entry"} // first row of blue_test
	for i, db := range testDBs {
		DBx := xDB(db)
		// query test
		if !DBx.Test_SELECT_ALL(testTable, desiredOutput) {
			t.Errorf("failed test query %s (%s)", db, testTable)
		} else {
			fmt.Printf("Query (%d) %s - PASS \n", i+1, db)
		}
	}

}
