package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items        []Item `json:"items"`
    Total        string `json:"total"`
}

type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}

type ReceiptScore struct {
	Id 	   uuid.UUID `json:"id"`
	Points int       `json:"points"`
}

func (r *Receipt) Validate() error {

	retailerPattern := `^[\w\s\-&]+$`
	if !regexp.MustCompile(retailerPattern).MatchString(r.Retailer) {
		return fmt.Errorf("invalid retailer: %s", r.Retailer)
	}

	if _, err := time.Parse("2006-01-02", r.PurchaseDate); err != nil {
		return fmt.Errorf("invalid purchaseDate: %s", r.PurchaseDate)
	}

	if _, err := time.Parse("15:04", r.PurchaseTime); err != nil {
		return fmt.Errorf("invalid purchaseTime: %s", r.PurchaseTime)
	}

	totalPattern := `^\d+\.\d{2}$`
	if !regexp.MustCompile(totalPattern).MatchString(r.Total) {
		return fmt.Errorf("invalid total: %s", r.Total)
	}

	if len(r.Items) < 1 {
		return errors.New("receipt must have at least one item")
	}
	for _, item := range r.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("invalid item %v", item)
		}
	}

	return nil
}

func (i *Item) Validate() error {

	descriptionPattern := `^[\w\s\-]+$`
	if !regexp.MustCompile(descriptionPattern).MatchString(i.ShortDescription) {
		return fmt.Errorf("invalid shortDescription: %s", i.ShortDescription)
	}

	pricePattern := `^\d+\.\d{2}$`
	if !regexp.MustCompile(pricePattern).MatchString(i.Price) {
		return fmt.Errorf("invalid price: %s", i.Price)
	}

	return nil
}
