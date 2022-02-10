package fetcher

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type CurrencyFetcher struct {
	DB *sql.DB
}

func (c *CurrencyFetcher) FetchAll() (map[int]string, error) {
	values := make(map[int]string)
	row, err := c.DB.Query("SELECT id, code FROM ISO4217;")
	if err != nil {
		log.Fatal(err)
	}
	for row.Next() {
		var id int
		var code string
		e := row.Scan(&id, &code)
		if e != nil {
			return values, e
		}
		values[id] = code
	}

	return values, nil
}

func (c *CurrencyFetcher) Fetch(id int) (string, error) {
	stmt, err := c.DB.Prepare("SELECT code FROM ISO4217 WHERE id = ? ;")
	if err != nil {
		log.Fatal(err)
	}

	var code string
	e := stmt.QueryRow(id).Scan(&code)
	if e != nil {
		return "", e
	}

	return code, nil
}
