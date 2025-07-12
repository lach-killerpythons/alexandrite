package JADE

import (
	"encoding/json"
	"fmt"
	"os"
)

const dataFile = "db.json"
const file_loc = "/home/lach/Babylon/GO/Alexandrite/JADE"

// API to hold all the DB credentials

type DB_creds struct {
	Name string
	Port string
	User string
	Host string
	PW   string
}

func GET_DB_creds(dataObj string) (DB_creds, error) {
	//var output []string

	Name := ""
	Port := ""
	User := ""
	Host := ""
	PW := ""

	var db_output DB_creds

	// make sure correct file location
	w_dir, _ := os.Getwd()
	if w_dir != file_loc {
		dir_err := os.Chdir(file_loc)
		if dir_err != nil {
			fmt.Println("error opening file location")
		}
	}

	// open json file
	jsonData, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("{\"error\": \"File %s not found in %s\"}", dataFile, w_dir)
		} else {
			fmt.Printf("{\"error\": \"Internal Server Error\"}")
		}
		return db_output, fmt.Errorf("json file error")
	}
	// unmarshall into a **string** interface
	var dat map[string]interface{}
	m_error := json.Unmarshal(jsonData, &dat)
	if m_error != nil {
		return db_output, fmt.Errorf("json marshall error")
	}

	if jsonObj, ok := dat[dataObj].(map[string]interface{}); ok {
		// type assertion
		if dbName, ok := jsonObj["DB"].(string); ok {
			//fmt.Println(dbName)
			Name = dbName
		}
		if dbPort, ok := jsonObj["PORT"].(string); ok {
			//fmt.Println(dbPort)
			Port = dbPort
		}
		if dbUser, ok := jsonObj["USER"].(string); ok {
			//fmt.Println(dbUser)
			User = dbUser
		}
		if dbHost, ok := jsonObj["HOST"].(string); ok {
			//fmt.Println(dbHost)
			Host = dbHost
		}
		if dbPW, ok := jsonObj["PW"].(string); ok {
			PW = dbPW
		}
	} else {
		return db_output, fmt.Errorf("%s cannot mapped to parent jsonObj key", dataObj)
	}
	if Name == "" || Port == "" || User == "" || Host == "" || PW == "" {
		fmt.Println("returned blank credential")
		return db_output, fmt.Errorf("missing credential")
	}

	//return []string{Name, Port, User, Host}
	return DB_creds{Name, Port, User, Host, PW}, nil

}

func GET_DB_credentials(dataObj string) []string {
	var output []string
	var Name string
	var Port string
	var User string
	var Host string

	jsonData, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("{\"error\": \"File %s not found\"}", dataFile)
		} else {
			fmt.Printf("{\"error\": \"Internal Server Error\"}")
		}
		return output
	}

	var dat map[string]interface{}
	m_error := json.Unmarshal(jsonData, &dat)
	if m_error != nil {
		panic(m_error)
	}
	//dummy_data = dat
	// ss := dat["localhost"]
	// if dbName, ok := ss["local_Dbname"].(string); ok {
	// 	fmt.Println("Extracted local_Dbname:", dbName)
	// }

	//misingo := "another_object"

	if datObj, ok := dat[dataObj].(map[string]interface{}); ok {
		// type assertion
		if dbName, ok := datObj["DB"].(string); ok {
			//fmt.Println(dbName)
			Name = dbName
		}
		if dbPort, ok := datObj["PORT"].(string); ok {
			//fmt.Println(dbPort)
			Port = dbPort
		}
		if dbUser, ok := datObj["USER"].(string); ok {
			//fmt.Println(dbUser)
			User = dbUser
		}
		if dbHost, ok := datObj["HOST"].(string); ok {
			//fmt.Println(dbHost)
			Host = dbHost
		}
	}
	return []string{Name, Port, User, Host}

}
