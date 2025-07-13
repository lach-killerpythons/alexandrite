package JADE

// API to hold all the DB credentials
import (
	"encoding/json"
	"fmt"
	"os"
)

const dataFile = "db.json"
const file_loc = "/home/lach/Babylon/GO/Alexandrite/JADE"

type JADE_FILE struct {
	Name     string // should be .json
	Location string
}

// see example JSON file
type DB_creds struct {
	Name string
	Port string
	User string
	Host string
	PW   string
}

// hardcoded
func GET_DB_creds(dataObj string) (DB_creds, error) {
	fmt.Println("Warning! Using hardcoded JadeFile ")
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

func OPEN_DB_creds(myFile JADE_FILE, dataObj string) (DB_creds, error) { // hardcoded
	//var output []string

	fmt.Println("Warning! Using hardcoded JadeFile ")

	Name := ""
	Port := ""
	User := ""
	Host := ""
	PW := ""

	var db_output DB_creds

	// make sure correct file location
	dir_err := os.Chdir(myFile.Location)
	if dir_err != nil {
		fmt.Println("error opening file location")
	}

	jsonData, err := os.ReadFile(myFile.Name)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("{\"error\": \"File %s not found in %s\"}", myFile.Name, myFile.Location)
		} else {
			fmt.Printf("{\"error\": \"Internal Server Error\"}")
		}
		return db_output, fmt.Errorf("json file error")
	}
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
