package BLUE

import (
	"fmt"
	"testing"
)

// test connectivitys
func TestBlue(t *testing.T) {
	test_DB := "local"
	result := Test_ConnectDB(test_DB)
	if !result {
		t.Errorf("Ping (1) %s - fail", test_DB)
	} else {
		fmt.Printf("Ping (1) %s - pass \n", test_DB)
	}

	test_DB = "mac"
	result = Test_ConnectDB(test_DB)
	if !result {
		t.Errorf("Ping (2) %s - fail", test_DB)
	} else {
		fmt.Printf("Ping (2) %s - pass \n", test_DB)
	}

	test_DB = "pi"
	result = Test_ConnectDB(test_DB)
	if !result {
		t.Errorf("Ping (3) %s - fail", test_DB)
	} else {
		fmt.Printf("Ping (3) %s - pass \n", test_DB)
	}

	// test_DB = "laptop"
	// result = Test_ConnectDB(test_DB)
	// if !result {
	// 	t.Errorf("Test (4) %s - fail", test_DB)
	// } else {
	// 	fmt.Printf("Test (4) %s - pass \n", test_DB)
	// }

}
