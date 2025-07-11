package BLUE

import (
	//"fmt"
	"fmt"
	"testing"
)

// test connectivitys
func TestDt(t *testing.T) {

	test_table := "blue_test"

	DBx := xDB("local")
	if !TableExists(DBx, test_table) {
		t.Errorf("%s does not exist!", test_table)
	}

	dt, err := DescribeTable(DBx, test_table)
	if err != nil {
		t.Errorf("/dt failed %s", err)

	}
	fmt.Println(test_table, "col_names:")
	fmt.Println(dt)

}
