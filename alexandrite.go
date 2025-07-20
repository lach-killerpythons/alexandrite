package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lach-killerpythons/alexandrite/BLUE"
	"github.com/lach-killerpythons/alexandrite/JADE"
	_ "github.com/lach-killerpythons/alexandrite/JADE"
	"github.com/lach-killerpythons/alexandrite/RED"
)

func URL_Get(baseURL string) (string, error) {
	// parsedURL, err := url.Parse(baseURL)
	// if err != nil {
	// 	return "", fmt.Errorf("invalid URL: %w", err)
	// }
	//robotsURL := fmt.Sprintf("%s://%s/robots.txt", parsedURL.Scheme, parsedURL.Host)
	// set timeout
	client := &http.Client{
		Timeout: 120 * time.Second, // Set timeout to 120 seconds
	}

	resp, err := client.Get(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch robots.txt: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil

}

// eg: select website from sandstone where robots LIKE '%shopify%';
func BlueToRed_1dList(rdb RED.RedDB, db BLUE.DB, query string, keyname string) error {
	if rdb.GetP() == nil {
		return fmt.Errorf("RDB pointer error : not found")
	}
	if db.Table.Name == "" {
		return fmt.Errorf("psql must have a table selected")
	}

	//myQuery := fmt.Sprintf(`SELECT * FROM %s;`, tableName)
	rows, err := db.GetP().Query(query)
	if err != nil {
		return err
	}

	var nextLine string
	for rows.Next() {
		if err := rows.Scan(&nextLine); err != nil {
			fmt.Printf("Error scanning row: %v", err)
			continue // Skip to the next row if scan fails for this one
		}
		//fmt.Println(nextLine)
		rdb.List_Add(keyname, nextLine)
	}

	return nil
}

func main() {
	//db, err := BLUE.DB_Connect("local")
	//cwd, _ := os.Getwd() // fix jade to get CWD
	//jf := JADE.JADE_FILE{"db.json", cwd + "/JADE"}
	jf := JADE.Open("db.json", "/JADE")
	db, err := BLUE.DB_JadeConnect("local", jf)
	if err != nil {
		panic(err)
	}
	OpenTable_err := db.OpenTable("sandstone")
	if OpenTable_err != nil {
		fmt.Println(OpenTable_err)
	}

	test_url := `https://kitandcradle.com.au/`
	sm, err := URL_Get(test_url + `sitemap.xml`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sm)

	// rdb, err := RED.NewRedDB("localhost", "", 2)
	// if err != nil {
	// 	panic(err)
	// }

	// add the URL list from psql to redis

	// myQuery := "select website from sandstone where robots LIKE '%shopify%';"
	// error_Blue2Red := BlueToRed_1dList(rdb, db, myQuery, "shopify_websites")
	// if error_Blue2Red != nil {
	// 	fmt.Println(error_Blue2Red)
	// }

	// rdb.List2_1wordset("shopify_websites", "web_set")

}
