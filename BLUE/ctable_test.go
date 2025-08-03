package BLUE

import (
	"testing"
)

// test connectivitys
func TestCreate(t *testing.T) {
	inputA := []interface{}{
		"James",
		42,
	}

	DB, err := DB_Connect("local")
	if err != nil {
		t.Errorf("db not connect")
	}
	inputB := []string{"name", "age"}

	inputC := []interface{}{
		"Penrith",
		"Green",
	}

	inputD := []string{"address", "color"}

	DB.CREATE_TABLE_v1("testPeople", inputA, inputB, PrimaryKey{"id"})
	ss := DB.CREATE_TABLE_v1("testAddress", inputC, inputD, PrimaryKey{"id"}, ForeignKey{"f_key", "INT", "testPeople", "id"})
	if ss != nil {
		t.Error(ss)
	}
}

