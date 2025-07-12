package BLUE

//package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	JADE "github.com/lach-killerpythons/alexandrite/JADE"

	_ "github.com/lib/pq"
)

var (
	//KeyDB *redis.Client
	CTX = context.Background()
)

type DB struct {
	p        *sql.DB // pointer to DB
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Table    myTable
}

type myTable struct {
	Name string   // table name
	Cols []string // col names
}

type PrimaryKey struct {
	Name string
}

type ForeignKey struct {
	Name              string
	Type              string
	ReferenceTable    string
	ReferenceVariable string
}

type myTable2 struct {
	Name   string   // table name
	Cols   []string // col names
	MyPKEY PrimaryKey
}

// need this type to use psql sequences
type id_sequence struct {
	name string
}

// deprecate
func xDB(name string) DB {
	var emptyDB DB
	dbc, err := JADE.GET_DB_creds(name)
	if err != nil {
		fmt.Println("JADE credentials error:", err)
		return emptyDB
	}
	fmt.Println(dbc.Name, dbc.User)
	port_int, err := strconv.Atoi(dbc.Port)
	if err != nil {
		fmt.Println(err)
		port_int = 5432
	}
	newDB := DB{
		Host:     dbc.Host,
		Port:     port_int,
		User:     dbc.User,
		Password: dbc.PW,
		Dbname:   dbc.Name,
	}
	return connectDB(newDB)

}

func DB_Connect(name string) (DB, error) {
	var emptyDB DB
	dbc, err := JADE.GET_DB_creds(name)
	if err != nil {
		fmt.Println("JADE credentials error:", err)
		return emptyDB, err
	}
	//fmt.Println(dbc.Name, dbc.User)
	port_int, err := strconv.Atoi(dbc.Port)
	if err != nil {
		fmt.Println(err)
		port_int = 5432
	}
	db := DB{
		Host:     dbc.Host,
		Port:     port_int,
		User:     dbc.User,
		Password: dbc.PW,
		Dbname:   dbc.Name,
	}
	var errOpen error
	// connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Dbname)
	//psql parser
	db.p, errOpen = sql.Open("postgres", psqlInfo)
	if errOpen != nil {
		fmt.Println(psqlInfo)
		return emptyDB, errOpen
	}
	errPing := db.p.Ping()

	if errPing != nil {
		return emptyDB, errPing
	}
	// error is nil if no errors
	return db, nil

}

// for ping_test
func Test_ConnectDB(name string) bool {
	dbc, err := JADE.GET_DB_creds(name)
	if err != nil {
		fmt.Println(err)
		return false
	}
	port_int, err := strconv.Atoi(dbc.Port)
	if err != nil {
		fmt.Println(err)
		port_int = 5432
	}
	db := DB{
		Host:     dbc.Host,
		Port:     port_int,
		User:     dbc.User,
		Password: dbc.PW,
		Dbname:   dbc.Name,
	}
	var errOpen error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Dbname)

	db.p, errOpen = sql.Open("postgres", psqlInfo)
	if errOpen != nil {
		fmt.Printf("failed to open connection \n %s \n", psqlInfo)
		return false // ❌
	}
	errPing := db.p.Ping()

	if errPing != nil {
		fmt.Printf("failed ping %s \n", psqlInfo)
		return false // ❌
	}
	return true // ✅

}

// for query test
func (DBx DB) Test_SELECT_ALL(tableName string, desiredOutput []string) bool {
	//var point *sql.DB
	db_pointer := DBx.p
	myQuery := fmt.Sprintf(`SELECT * FROM %s;`, tableName)
	if db_pointer != nil {
		rows, err := db_pointer.Query(myQuery)
		if err != nil {
			fmt.Println("🐸")
			fmt.Println(err)
			return false
		}
		// get column names
		cols, err := rows.Columns()
		// empty interface array to hold values by name
		values := make([]interface{}, len(cols))
		// empty interface array for pointers to values
		valuePtrs := make([]interface{}, len(cols))

		for rows.Next() {
			for i := range values {

				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				fmt.Println("row scan failed: %w", err)
				return false
			}

			var line []string
			for _, val := range values {
				var str string
				if val != nil {
					str = fmt.Sprintf("%v", val)
				} else {
					str = "NULL"
				}
				line = append(line, str)
			}
			//fmt.Println(line)
			for i := range line {
				if line[i] != desiredOutput[i] {
					return false
				}
			}

			break // break after first line
		}
	}
	return true
}

// deprecate
func connectDB(db DB) DB {
	var errOpen error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Dbname)
	db.p, errOpen = sql.Open("postgres", psqlInfo)
	if errOpen != nil {
		fmt.Println(psqlInfo)
		panic(errOpen)
	}
	errPing := db.p.Ping()

	if errPing != nil {
		fmt.Printf("connection failed ❌ \n %s", errPing)
	} else {
		fmt.Println("connection success ✅")
	}
	return db
}

// deprecate
func (DBx *DB) SetTable(name string, col_names []string) {
	DBx.Table.Name = name
	DBx.Table.Cols = col_names
}

// gets the col names if exists
func (DBx *DB) OpenTable(tableName string) error {
	// check if table exists
	var err error
	if !TableExists(*DBx, tableName) { // derefencer
		err = fmt.Errorf("open table failed : !TableExists ")
		return err
	}
	DBx.Table.Name = tableName
	DBx.Table.Cols, err = DescribeTable(*DBx, tableName)
	if err != nil {
		err = fmt.Errorf("open table failed : DescribeTable error")
		return err // DescribeTable failed
	}
	return nil
}

func get_SQL_TYPE(input interface{}) (string, error) {
	var typeStr string

	type SQL_TYPE string
	//var pkey PrimaryKey

	const (
		SQL_TEXT       SQL_TYPE = "TEXT"
		SQL_INT        SQL_TYPE = "INT"
		SQL_FLOAT      SQL_TYPE = "DOUBLE PRECISION"
		SQL_BYTE_ARRAY SQL_TYPE = "BYTEA"
		SQL_TEXT_ARRAY SQL_TYPE = "TEXT[]"
	)

	switch input.(type) {
	case int:
		typeStr = string(SQL_INT)
	case string:
		typeStr = string(SQL_TEXT)
	case float64:
		typeStr = string(SQL_FLOAT)
	case []byte:
		typeStr = string(SQL_BYTE_ARRAY)
	case []string:
		typeStr = string(SQL_TEXT_ARRAY)
	default:
		return "", fmt.Errorf("type not found")
	}

	return typeStr, nil
}

// infer text or int
func (DBx *DB) CREATE_TABLE_v1(tableName string, cols []interface{}, colNames []string, args ...any) error {
	var err error = nil
	var result sql.Result

	// create the SQL body string
	colStr := func() string {
		var colBody []string // BODY OF SQL CREATE TABLE
		var typeStr string   // SQL_TYPE -> String

		// insert primary key
		if primary_key, ok := args[0].(PrimaryKey); ok { // type check arg0
			p_str := primary_key.Name + " SERIAL PRIMARY KEY" // autoincrement
			colBody = append(colBody, p_str)
		} else {
			fmt.Println("no primary key specified!")
		}

		// insert foreign key
		for _, arg := range args[1:] { // skip first key
			if f_key, ok := arg.(ForeignKey); ok {
				// department_id INTEGER,
				f_str_1 := fmt.Sprintf("%s %s,\n", f_key.Name, f_key.Type)
				// FOREIGN KEY (department_id) REFERENCES departments(id)
				f_str_2 := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", f_key.Name, f_key.ReferenceTable, f_key.ReferenceVariable) // department_id INTEGER,
				f_str := f_str_1 + f_str_2
				colBody = append(colBody, f_str)
			}
		}

		// insert other variables with SQL types
		for i, val := range cols {
			typeStr, err = get_SQL_TYPE(val)
			if err != nil {
				panic(err) // panic if no SQL type
			}
			line := colNames[i] + " " + typeStr
			colBody = append(colBody, line)
		}

		// join it all together
		return strings.Join(colBody, ",")
	}

	// SQL query
	createTableSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s
	);`, tableName, colStr())

	fmt.Println(createTableSQL)

	// execute CREATE TABLE
	result, err = DBx.p.Exec(createTableSQL)
	fmt.Println(result)
	if err != nil {
		return err
	}
	// update *DB
	DBx.SetTable(tableName, colNames)
	INSERT_err := DBx.INSERT_WILD(cols)
	if INSERT_err != nil {
		return INSERT_err
	}
	return err
}

// print every line in a sqlDB for a table in DB
func (DBx DB) SELECT_ALL(tableName string) {
	//var point *sql.DB
	db_pointer := DBx.p
	myQuery := fmt.Sprintf(`SELECT * FROM %s;`, tableName)
	fmt.Println(myQuery)
	if db_pointer != nil {
		rows, err := db_pointer.Query(myQuery)
		if err != nil {
			fmt.Println("🐸")
			fmt.Println(err)
			return
		}
		// get column names
		cols, err := rows.Columns()
		// empty interface array to hold values by name
		values := make([]interface{}, len(cols))
		// empty interface array for pointers to values
		valuePtrs := make([]interface{}, len(cols))

		for rows.Next() {
			for i := range values {

				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				fmt.Println("row scan failed: %w", err)
				return
			}

			var line []string
			for _, val := range values {
				var str string
				if val != nil {
					str = fmt.Sprintf("%v", val)
				} else {
					str = "NULL"
				}
				line = append(line, str)
			}
			fmt.Println(line)
		}
	}
}

func (DBx DB) QueryToJSON(query string) ([]byte, error) {
	db := DBx.p
	if db == nil {
		e := errors.New("*sql.DB == nil")
		fmt.Println(e)
		return nil, e
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []map[string]interface{}{}

	for rows.Next() {
		// Create a slice of interface{}s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// handle []byte to string conversion
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		results = append(results, rowMap)
	}

	return json.MarshalIndent(results, "", "  ")
}

func (DBx DB) INSERT_WILD(vals []interface{}) error {

	// concat any types for SQL query (type inference) 🔎
	concat2 := func(input []interface{}) string {
		output := ""
		for _, data := range input {
			s := ""
			switch v := data.(type) {
			case float64:
				s = strconv.FormatFloat(v, 'f', -1, 64)
			case int:
				s = strconv.Itoa(v)
			case []string:
				s = fmt.Sprintf("'{%s}'", strings.Join(v, ","))
			// case float64:
			// 	s = strconv.FormatFloat(v, f, 6, 64)
			case id_sequence:
				s = fmt.Sprintf("nextval(%s)", v.name)
			default: //default is string
				s = fmt.Sprintf(`'%s'`, v)

			}
			output += s + ","
		}
		return output[:len(output)-1]
	}

	// if pointer is not empty
	if DBx.p != nil {
		tableName := DBx.Table.Name
		colNameStr := strings.Join(DBx.Table.Cols, ",")
		db := DBx.p
		if len(vals) != len(DBx.Table.Cols) {
			return fmt.Errorf("invalid number of args v column names")
		}
		query := fmt.Sprintf(`insert into %s (%s) VALUES (%s);`, tableName, colNameStr, concat2(vals))

		//test purposes
		fmt.Println(query)

		i, err := db.Exec(query)
		if err != nil {
			fmt.Println(err)
			return err
		}
		j, _ := i.RowsAffected()
		fmt.Println("INSERT 0", j, "~", tableName)
	}
	return nil

}

// return the col names
func DescribeTable(myDB DB, tableName string) ([]string, error) {
	var output []string
	// `/d golden` -> list of col names
	// query the information_schema.columns
	dt_query := fmt.Sprintf(`SELECT column_name FROM information_schema.columns
	 WHERE table_schema = 'public' AND table_name = '%s';`, tableName)

	result, err := myDB.p.QueryContext(CTX, dt_query)
	if err != nil {
		fmt.Println(err) // not exist / fail
		return output, err
	}
	//output, err = result.Columns()
	for result.Next() {
		var col string
		if err := result.Scan(&col); err != nil {
			return []string{}, err
		}
		output = append(output, col)
	}

	return output, nil

}

// does table exist
func TableExists(db DB, tableName string) bool {
	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = $1
		);`

	err := db.p.QueryRowContext(CTX, query, tableName).Scan(&exists)
	if err != nil {
		fmt.Println(err)
	}

	if exists {
		fmt.Printf("%s exists ✅ \n", tableName)
	} else {
		fmt.Printf("%s does NOT exist. \n", tableName)
	}
	return exists
}

// func main() {

// }
