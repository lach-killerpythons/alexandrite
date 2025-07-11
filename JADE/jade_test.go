package JADE

import (
	"fmt"
	"testing"
)

func TestJade(t *testing.T) {
	ss, err := GET_DB_creds("mac")
	if err != nil {
		t.Errorf("JADE GET_DB_creds ERROR: %s", err)
	}
	fmt.Println(ss)

	result := GET_DB_credentials("mac")
	expected := []string{"tiamat", "5432", "hit_target", "minime.local", "1991", "1991"}
	if result[0] != expected[0] {
		t.Errorf("get_db_cred = %s, want %s", result, expected)
	}
}
