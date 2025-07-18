package RED

import (
	"fmt"
	"testing"
)

func TestRed(t *testing.T) {
	fmt.Println("testing")
	rdb, err := NewRedDB("localhost", "", 0)
	if err != nil {
		t.Errorf("failed %s", err)
	}
	ss := GetAllKeys(rdb.GetP())
	fmt.Println(string(ss))

}

// migrate these to tests
// func main() {

// 	piRedDB := NewRedDB("pi.local", "pi", 0)
// 	pi_rdbClient, err := piRedDB.New_Connection()
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		_ = GetAllKeys(pi_rdbClient)
// 	}

// 	lRedDB := NewRedDB("localhost", "", 0)
// 	localClient, err := lRedDB.New_Connection()
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		_ = GetAllKeys(localClient)
// 	}

// 	lRedDB.db_number = 1
// 	localClient, err = lRedDB.New_Connection()
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		_ = GetAllKeys(localClient)

// 	}

// }
