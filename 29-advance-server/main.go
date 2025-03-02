package main

import (
	"fmt"
	"log"
	"main/db"
	"net/http"
	"strconv"
)

type ClothStore struct {
	db *db.Database
}

func (cs *ClothStore) list(w http.ResponseWriter, r *http.Request) {
	items := cs.db.GetAllItems()
	for item, price := range items {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (cs *ClothStore) price(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	items, ok := r.Form["item"]
	if !ok {
		log.Fatal("can't find item")
	}
	if len(items) != 1 {
		log.Fatalf("expect item slice to be size 1 but got: %d", len(items))
	}
	item := items[0]
	price, err := cs.db.GetPrice(item)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s\n", price)
}

func (cs *ClothStore) create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	items, ok := r.Form["item"]
	if !ok {
		log.Fatal("can't find item")
	}
	if len(items) != 1 {
		log.Fatalf("expect item slice to be size 1 but got: %d", len(items))
	}
	item := items[0]

	prices, ok := r.Form["price"]
	if !ok {
		log.Fatal("can't find price")
	}
	if len(prices) != 1 {
		log.Fatalf("expect price slice to be size 1 but got: %d", len(prices))
	}
	priceStr := prices[0]
	priceF64, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		log.Fatal(err)
	}
	price := db.Dollars(priceF64)

	cs.db.CreateItem(item, price)
}

func main() {
	// Each handle function is called on a different goroutine so need to make sure access
	// to the database is thread safe.
	defaults := map[string]db.Dollars{"shoes": 50, "socks": 5}
	db := db.NewDatabase(defaults)
	cs := ClothStore{db: db}

	http.HandleFunc("GET /list", cs.list)
	http.HandleFunc("GET /price", cs.price)
	http.HandleFunc("POST /price", cs.create)
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) // tells go to use global http handler
}
