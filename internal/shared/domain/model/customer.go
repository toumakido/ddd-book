package model

import "github.com/google/uuid"

type CustomerID string

type Address struct {
	PostalCode string
	Prefecture string
	City       string
	Street     string
	Building   string
}

func NewID() string {
	return uuid.New().String()
}
