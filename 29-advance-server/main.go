package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range *db {
		fmt.Fprintf(w, "%v: %v\n", item, price)
	}
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := (*db)[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "can not find %v\n", item)
		return
	}
	fmt.Fprintf(w, "%v\n", price)
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "body is invalid\n")
	}
	text := string(data)
	var item string
	var price dollars
	_, err = fmt.Sscanf(text, "%s%v", &item, &price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "body is invalid\n")
	}
	(*db)[item] = price
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	_, ok := (*db)[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "can not find %v\n", item)
		return
	}
	data, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "body is invalid\n")
	}
	text := string(data)
	var newPrice dollars
	_, err = fmt.Sscanf(text, "%v", &newPrice)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "body is invalid2\n")
	}
	(*db)[item] = newPrice
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	_, ok := (*db)[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "can not find %v\n", item)
		return
	}
	delete(*db, item)
}

func main() {
	// there are a lot of things wrong with this example
	// Each handle function is called on a different goroutine so need to make sure access
	// to the database is thread safe. Also it is good practise to seperate the handle interface
	// from the database logic which this example does not do.
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("GET /list", db.list)
	http.HandleFunc("GET /price", db.price)
	http.HandleFunc("POST /price", db.create)
	http.HandleFunc("MODIFY /price", db.update)
	http.HandleFunc("DELETE /price", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) // tells go to use global http handler
}
