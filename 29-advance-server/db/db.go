package db

import (
	"fmt"
)

type Dollars float32

func (d Dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type Request interface {
}

type getListRequest struct {
	response chan map[string]Dollars
}

type getPriceResponse struct {
	price Dollars
	err   error
}

type getPriceRequest struct {
	item     string
	response chan getPriceResponse
}

type createItemRequest struct {
	item  string
	price Dollars
}

type Database struct {
	data    map[string]Dollars
	request chan Request
}

// Uses monitor goroutine patttern to prevent data races
func (db *Database) broker() {
	for request := range db.request {
		switch r := request.(type) {
		case getListRequest:
			// need to return a deep copy of the database otherwise race conditions can occur with the reference
			c := make(map[string]Dollars)
			for k, v := range db.data {
				c[k] = v
			}
			r.response <- c
		case getPriceRequest:
			res := getPriceResponse{}
			price, ok := db.data[r.item]
			if !ok {
				res.err = fmt.Errorf("item: %s is not in database", r.item)
				r.response <- res
			}
			res.price = price
			r.response <- res
		case createItemRequest:
			db.data[r.item] = r.price
		default:
			panic("uknown type")
		}
	}
}

func (db *Database) GetAllItems() map[string]Dollars {
	listRequest := getListRequest{}
	listRequest.response = make(chan map[string]Dollars)
	db.request <- listRequest
	return <-listRequest.response
}

func (db *Database) GetPrice(item string) (Dollars, error) {
	getRequest := getPriceRequest{item: item, response: make(chan getPriceResponse)}
	db.request <- getRequest
	priceResponse := <-getRequest.response
	return priceResponse.price, priceResponse.err
}

func (db *Database) CreateItem(item string, price Dollars) {
	db.request <- createItemRequest{item: item, price: price}
}

func NewDatabase(defaults map[string]Dollars) *Database {
	db := Database{}
	db.data = make(map[string]Dollars)
	for k, v := range defaults {
		db.data[k] = v
	}
	db.request = make(chan Request)
	go db.broker()
	return &db
}
