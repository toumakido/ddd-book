package model

type CustomerID string

type Address struct {
	PostalCode string
	Prefecture string
	City       string
	Street     string
	Building   string
}
