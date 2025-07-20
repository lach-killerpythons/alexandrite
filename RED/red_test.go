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
	ss := rdb.GetAllKeys()
	fmt.Println(string(ss))

}

func Print(args ...any) {
	for i, v := range args {
		fmt.Printf("[%d] %v\n", i, v)
	}
}

// ✅ Function: Convert variadic ...any to []string
func ToStringSlice(values ...any) []string {
	strs := make([]string, len(values))
	for i, v := range values {
		strs[i] = fmt.Sprint(v)
	}
	return strs
}

// ✅ Functional polymorphism — you're passing different functions to change the behavior of ListDo.
func TestDo(t *testing.T) {

	fmt.Println("Print(...)")
	rdb, err := NewRedDB("localhost", "", 2)
	if err != nil {
		t.Errorf("failed %s", err)
	}
	rdb.ListDo("jacaru_products", Print)

	// pass additional function logic within the instance inside the ListDo(func(){...})
	rdb.ListDo("jacaru_products", func(args ...any) {
		stringList := ToStringSlice(args...)
		fmt.Printf("Got %d items:\n", len(stringList))
		for _, s := range stringList {
			fmt.Println("-", s)
		}
	})

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
