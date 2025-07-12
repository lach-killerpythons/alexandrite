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

var (
	KeyDB *redis.Client
	CTX   = context.Background()
)

func NewRedDB(address string, pw string, db_number int) RedDB {
	default_port := ":6379"
	return RedDB{
		addess:    address + default_port,
		pw:        pw,
		db_number: db_number,
	}
}

func (rdb RedDB) New_Connection() (*redis.Client, error) {
	var placeholder *redis.Client
	new_redis_client := redis.NewClient(&redis.Options{
		Addr:     rdb.addess,    // host:port of the redis server
		Password: rdb.pw,        // no password set
		DB:       rdb.db_number, // use default DB
	})
	// ping client
	_, err := new_redis_client.Ping(CTX).Result()
	if err != nil {
		return placeholder, err
	}
	rdb.p = new_redis_client
	return new_redis_client, nil
}

// get all keys in a redis DB
func GetAllKeys(rdb *redis.Client) []byte {
	result, err := rdb.Keys(CTX, "*").Result()
	if err != nil {
		panic(err)
	}

	var sets []string
	var strings []string
	var lists []string
	var zset []string
	var hashs []string

	for _, val := range result {
		myType, err := rdb.Type(CTX, val).Result()
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

	fmt.Println("Generated JSON:")
	fmt.Println(string(jsonData))

	return jsonData
}

func SetKey(rdb *redis.Client, keyname string, val interface{}) {
	switch v := val.(type) {
	case int:
		rdb.Set(CTX, keyname+"_str", v, 0)
	case string:
		rdb.Set(CTX, keyname+"_int", v, 0)
	default:
		fmt.Println("passed wrong type to setKey!", v)
	}
}

func AddToSet(rdb *redis.Client, setname string, inputs interface{}) {
	f, ok := inputs.([]string)
	if ok && len(f) == 0 {
		return
	}
	result, err := rdb.SAdd(CTX, setname, inputs).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("AddToSet", result)
}

// LIST FUNCTIONS
//  _       _________ _______ _________ _______
// ( \      \__   __/(  ____ \\__   __/(  ____ \
// | (         ) (   | (    \/   ) (   | (    \/
// | |         | |   | (_____    | |   | (_____
// | |         | |   (_____  )   | |   (_____  )
// | |         | |         ) |   | |         ) |
// | (____/\___) (___/\____) |   | |   /\____) |
// (_______/\_______/\_______)   )_(   \_______)

// Txt2List, List_RandItem

// convert a textfile into a redis list
func Txt2List(rdb *redis.Client, keyname string, targetFile string) {
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

// get a random list item
func List_RandItem(rdb *redis.Client, keyname string) string {
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
func List_Add(rdb *redis.Client, keyname string, val string) string { // 1 is success
	err := rdb.LPush(CTX, keyname, val).Err()
	if err != nil {
		fmt.Println("something failed", err)
		return "0"
	}
	fmt.Println(`"`, val, `"`, "-added to list-")
	return "1"
}

// del from a list (LRem)
func List_DelStr(key string, targetList string, kdb *redis.Client) (string, int) {
	result, err := kdb.Do(CTX, "LREM", targetList, 0, key).Result()
	i_result := int(result.(int64))
	fmt.Println("n removed:", i_result)
	if err != nil {
		fmt.Println("err:", err)
		return key, 0
	}
	return key, i_result

}

// return the json byte string
func List2JSON(rdb *redis.Client, listKey string) []byte {
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
func List2JSON_alpha(rdb *redis.Client, listKey string) ([]string, []byte) {
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

func List2_1wordset(rdb *redis.Client, listName string, setName string) {
	values, err := rdb.LRange(CTX, listName, 0, -1).Result()
	if err != nil {
		fmt.Printf("Failed to read Redis list: %v \n", err)
		return
	}
	for _, s := range values {
		str_split := strings.Split(s, " ")
		AddToSet(rdb, setName, str_split)
	}

}

//  _   _  _____ _____ _____
// | | | |/  ___|  ___|_   _|
// | |_| |\ `--.| |__   | |
// |  _  | `--. \  __|  | |
// | | | |/\__/ / |___  | |
// \_| |_/\____/\____/  \_/
//
// Store object-like data : { key1:1value, key2:value2...}
// 🟥 HGETALL h_test             -> return all keys and values
// 🟥 HVALS h_test               -> return all vals
// 🟥 HGET h_test c1             -> return h_test[c1] (value)
// 🟥 HSET h_test field1 "Hello" -> add {field1:"Hello"}
// 🟥 HINCRBY h_test joon 1      -> add 1 HSET (one str at a time)
// 🟥 HMGET h_test field1 field2 -> get multiple field values

func AddToHSet(xbd *redis.Client, hsetName, keyname string, val interface{}) {

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
	r, err := xbd.HSet(CTX, hsetName, keyname, valStr).Result()
	if err != nil {
		fmt.Println("failed to add to Hset", hsetName, r)
		fmt.Println(err)
		return
	}
	fmt.Println("hset add", r)

}

func CummulativeHSET(xdb *redis.Client, hsetName string, input string) {
	_, err := xdb.HIncrBy(CTX, hsetName, input, 1).Result()
	if err != nil {
		panic(err)
	}

}

//  ______ _____ _____ _____
// |___  //  ___|  ___|_   _|
//    / / \ `--.| |__   | |
//   / /   `--. \  __|  | |
// ./ /___/\__/ / |___  | |
// \_____/\____/\____/  \_/
//
// ZSET --Store unique, ordered collections with scores (rankings)
// 🟥 ZRANGE key start stop [WITHSCORES]: Get members by index (rank).
// 🟥 ZRANGEBYSCORE key min max [WITHSCORES]: Get members by score range.
// 🟥 ZREM key member: Remove a member.
// 🟥 ZSCORE key member: Get the score of a specific member.
// 🟥 ZRANK key member: Get the rank (0-based index) of a member.
// 🟥 ZINCRBY key increment member: Increment the score of a member

func CummulativeZSET(xdb *redis.Client, hsetName string, input string) {
	_, err := xdb.ZIncrBy(CTX, hsetName, 1, input).Result()
	if err != nil {
		panic(err)
	}

}

// check if list exists
func IsList(rdb *redis.Client, listKey string) bool {
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
