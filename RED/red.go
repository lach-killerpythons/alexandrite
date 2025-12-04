package RED

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedDB struct {
	p         *redis.Client
	addess    string
	pw        string
	db_number int
}

type R_TYPE string

const (
	R_LIST R_TYPE = "LIST"
	R_HSET R_TYPE = "HSET"
	R_ZSET R_TYPE = "ZSET"
	R_SET  R_TYPE = "SET"
)

func (rdb RedDB) GetP() *redis.Client {
	return rdb.p
}

var (
	CTX = context.Background()
)

func NewRedDB(address string, pw string, db_number int) (RedDB, error) {
	default_port := ":6379"
	rdb := RedDB{
		addess:    address + default_port,
		pw:        pw,
		db_number: db_number,
	}

	new_redis_client := redis.NewClient(&redis.Options{
		Addr:     rdb.addess,    // host:port of the redis server
		Password: rdb.pw,        // no password set
		DB:       rdb.db_number, // use default DB
	})
	// ping client
	_, err := new_redis_client.Ping(CTX).Result()
	if err != nil {
		return rdb, err
	}
	rdb.p = new_redis_client
	return rdb, nil
}

func (rdb RedDB) GetAllKeys() []byte {
	result, err := rdb.p.Keys(CTX, "*").Result()
	if err != nil {
		panic(err)
	}

	var sets []string
	var strings []string
	var lists []string
	var zset []string
	var hashs []string

	for _, val := range result {
		myType, err := rdb.p.Type(CTX, val).Result()
		if err != nil {
			panic(err)
		}
		//fmt.Println(val, myType)
		switch myType {
		case "set":
			sets = append(sets, val)
		case "string":
			strings = append(strings, val)
		case "list":
			lists = append(lists, val)
		case "zset":
			strings = append(zset, val)
		case "hash":
			lists = append(hashs, val)
		default:
			fmt.Println("no support currently for:", myType)
		}
	}

	type GroupedItems struct {
		Sets    []string `json:"sets"`
		Strings []string `json:"strings"`
		Lists   []string `json:"lists"`
		Zset    []string `json:"zset"`
		Hash    []string `json:"hash"`
	}

	GroupedKeys := GroupedItems{lists, strings, sets, zset, hashs}

	jsonData, err := json.MarshalIndent(GroupedKeys, "", "  ")
	if err != nil {
		// Handle any errors that occur during marshaling.
		fmt.Println("Error marshaling JSON:", err)
	}

	// dont need to dupe this
	//fmt.Println("Generated JSON:")
	//fmt.Println(string(jsonData))

	return jsonData
}

func (rdbX RedDB) SetKey(keyname string, val interface{}) {
	rdb := rdbX.p
	switch v := val.(type) {
	case int:
		rdb.Set(CTX, keyname+"_str", v, 0)
	case string:
		rdb.Set(CTX, keyname+"_int", v, 0)
	default:
		fmt.Println("passed wrong type to setKey!", v)
	}
}

// func (rdb RedDB)
func (rdbX RedDB) AddToSet(setname string, inputs interface{}) {
	rdb := rdbX.p

	// ADD []string or string
	// f, ok := inputs.([]string)
	// if ok && len(f) == 0 {
	// 	return
	// }

	result, err := rdb.SAdd(CTX, setname, inputs).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("AddToSet", result)
}

func (rdbX RedDB) Txt2List(keyname string, targetFile string) {
	rdb := rdbX.p
	file, err := os.Open(targetFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		rdb.LPush(CTX, keyname, line)

		fmt.Println("added to '", keyname, "' ~", line)
	}
}

// sometimes I want to export a RDB list out to a text file to use somewhere else
func (rdbX RedDB) List2Text(listname string, outputFile string) error {
	input, err := rdbX.BetterListGet(listname)
	if err != nil {
		fmt.Println("RedDB Error")
		return err
	}
	// Open the file for writing. Create it if it doesn't exist, append if it does.
	// 0644 sets the file permissions.
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed opening file: %s \n", err)
		return err
	}
	defer file.Close() // Ensure the file is closed when the function exits

	writer := bufio.NewWriter(file)

	for _, line := range input {
		fmt.Println(line)
		_, err := writer.WriteString(line + "\n") // Write the line and a newline character
		if err != nil {
			fmt.Printf("failed writing to file: %s\n", err)
		}
	}

	// Flush the buffer to ensure all data is written to the file
	err = writer.Flush()
	if err != nil {
		fmt.Printf("failed flushing buffer: %s\n", err)
	}

	fmt.Printf("Successfully exported RDB List '%s' to %s\n", listname, outputFile)

	return err
}

// get a random list item
func (rdbX RedDB) List_RandItem(keyname string) string {
	rdb := rdbX.p
	listLen, err := rdb.LLen(CTX, keyname).Result() //
	if err != nil {
		return err.Error()
	}
	random_i := rand.Intn(int(listLen)) // generate random index in range
	result, err := rdb.Do(CTX, "LRANGE", keyname, random_i, random_i).Result()
	if err != nil {
		return err.Error()
	}
	getItem := result.([]interface{})[0].(string)
	return getItem
}

// add to a list (Lpush)
func (rdbX RedDB) List_Add(keyname string, val string) string { // 1 is success
	rdb := rdbX.p
	err := rdb.LPush(CTX, keyname, val).Err()
	if err != nil {
		fmt.Println("something failed", err)
		return "0"
	}
	fmt.Println(`"`, val, `"`, "-added to list-")
	return "1"
}

// del from a list (LRem)
func (rdbX RedDB) List_DelStr(key string, targetList string) (string, int) {
	rdb := rdbX.p
	result, err := rdb.Do(CTX, "LREM", targetList, 0, key).Result()
	i_result := int(result.(int64))
	fmt.Println("n removed:", i_result)
	if err != nil {
		fmt.Println("err:", err)
		return key, 0
	}
	return key, i_result

}

// return the json byte string
func (rdbX RedDB) List2JSON(listKey string) []byte {
	rdb := rdbX.p
	var output []byte
	values, err := rdb.LRange(CTX, listKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("Failed to read Redis list: %v \n", err)
		return output
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		fmt.Printf("Failed to parse JSON: %v \n", err)
		return output
	}
	return jsonData
}

// returns the string array too
func (rdbX RedDB) List2JSON_alpha(listKey string) ([]string, []byte) {
	rdb := rdbX.p
	var output []byte
	var strArr []string
	values, err := rdb.LRange(CTX, listKey, 0, -1).Result()
	if err != nil {
		fmt.Printf("Failed to read Redis list: %v \n", err)
		return strArr, output
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		fmt.Printf("Failed to parse JSON: %v \n", err)
		return strArr, output
	}
	return values, jsonData
}

func (rdbX RedDB) List2_1wordset(listName string, setName string) {
	rdb := rdbX.p
	values, err := rdb.LRange(CTX, listName, 0, -1).Result()
	if err != nil {
		fmt.Printf("Failed to read Redis list: %v \n", err)
		return
	}
	for _, s := range values {
		str_split := strings.Split(s, " ")

		rdbX.AddToSet(setName, str_split)
	}

}

// Run any function on an entire redis list
// see red_test.go for example functions

// best example function is
// printItems := func(items ...any) {
//     for _, item := range items {
//         fmt.Println(item)
//     }
// }

// this is a bit of a retarded way of forcing any type
// this should always return a []string - so this isn't really needed
func (rdbX RedDB) ListDo(listName string, f func(...any)) {
	result, err := rdbX.GetP().Do(CTX, "LRANGE", listName, 0, -1).Result()
	if err != nil {
		fmt.Printf("Redis LRANGE error: %v", err)
		return
	}

	items, ok := result.([]interface{})
	if !ok {
		fmt.Printf("Unexpected result type: %T", result)
		return
	}

	f(items...)

}

// use the LRange result directly and get the []string
func (rdbX RedDB) BetterListGet(listName string) ([]string, error) {
	return rdbX.GetP().LRange(CTX, listName, 0, -1).Result()
}

// AnyDo --> get any redis key type and do something

// (pass key name) -> get type

// switch for type -> get all

func (rdbX RedDB) AddToHSet(hsetName, keyname string, val interface{}) {
	rdb := rdbX.p

	var valStr string
	switch v := val.(type) {
	case int:
		valStr = strconv.Itoa(v)
	case string:
		valStr = v
	default:
		fmt.Println("passed wrong type to setKey!", v)
		return
	}
	r, err := rdb.HSet(CTX, hsetName, keyname, valStr).Result()
	if err != nil {
		fmt.Println("failed to add to Hset", hsetName, r)
		fmt.Println(err)
		return
	}
	fmt.Println("hset add", r)

}

// 💡 wildtype -> is H v Z

func (rdbX RedDB) CummulativeHSET(hsetName string, input string) {
	rdb := rdbX.p
	_, err := rdb.HIncrBy(CTX, hsetName, input, 1).Result()
	if err != nil {
		panic(err)
	}

}

func (rdbX RedDB) CummulativeZSET(hsetName string, input string) {
	rdb := rdbX.p
	_, err := rdb.ZIncrBy(CTX, hsetName, 1, input).Result()
	if err != nil {
		panic(err)
	}

}

// check if list exists
// 💡 implement wildtype - IsX (key string, rt REDIS_TYPE)
func (rdbX RedDB) IsList(listKey string) bool {
	rdb := rdbX.p
	keyType, err := rdb.Type(CTX, listKey).Result()
	if err != nil {
		fmt.Println(err)
	}
	switch keyType {
	case "list":
		fmt.Println(listKey, "-> IsList ✅")
		return true
	default:
		fmt.Println(listKey, "-> IsList ❌")
		return false
	}

}

// 💡 Feature ideas
// faux string enum for redis types
//
