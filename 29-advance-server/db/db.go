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

func (db *Database) broker() {
	for request := range db.request {
		switch r := request.(type) {
		case getListRequest:
			r.response <- db.data
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

func NewDatabase(defaults *map[string]Dollars) *Database {
	db := Database{}
	db.data = make(map[string]Dollars)
	if defaults != nil {
		db.data = *defaults
	}
	db.request = make(chan Request)
	go db.broker()
	return &db
}
