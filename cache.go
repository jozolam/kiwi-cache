package main

import (
	"database/sql"
	"fmt"
	"kiwi/cache/cache_impl"
	"kiwi/cache/fetcher"
	"time"
)

type locationCache struct {
	cache map[int]string
}

func main() {
	countryCache := cache_impl.InitCache(1000, &fetcher.CountryFetcher{}, "countryCache")
	sqliteDatabase, _ := sql.Open("sqlite3", "currencies.db") // Open the created SQLite File
	defer sqliteDatabase.Close()
	currencyCache := cache_impl.InitCache(1000, &fetcher.CurrencyFetcher{DB: sqliteDatabase}, "currencyCache")

	for i := 0; i < 100; i++ {
		go func() {
			//index := i
			v1, e := countryCache.Get(2208)
			if e != nil {
				fmt.Println(e)
			}
			fmt.Println(" country value is ", v1)

			v2, e1 := currencyCache.Get(979)
			if e1 != nil {
				fmt.Println(e1)
			}
			fmt.Println(" currency value is ", v2)
		}()
	}

	time.Sleep(2000 * time.Millisecond)
}
