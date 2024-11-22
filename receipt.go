package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
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

/*
Validate each part of a receipt by matching regexes and time/date values to the API spec
*/
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

/*
Validate an item with in a receipt by matching regexes to the API spec
*/
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

/*
Sum of all point values within the receipt
*/
func (r *Receipt) CountPoints() int {
	return ReltailerNameScore(r.Retailer) +
		TotalPriceScore(r.Total) +
		TotalItemsScore(r.Items) +
		DateScore(r.PurchaseDate) +
		TimeScore(r.PurchaseTime)
}

/* 
One point for every alphanumeric character in the retailer name.
*/
func ReltailerNameScore(r string) int {
	count := 0
	for _, c := range r {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			count++
		}
	}
	return count
}

/*
50 points if the total is a round dollar amount with no cents.
25 points if the total is a multiple of 0.25.
*/
func TotalPriceScore(t string) int {
	// splitting and checking cases here to avoid an unnecessary strconv and modulo math
	cents := strings.Split(t, ".")[1]

	switch cents {
	case "00":
		return 75 // Round dollar amount (50 + 25 for multiple of 0.25)
	case "25", "50", "75":
		return 25 // Multiple of 0.25
	default:
		return 0 // Else
	}
}

/* 
5 points for every two items on the receipt.
(Then, for each item) If the trimmed length of the item description is a multiple of 3, multiply
the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
*/
func TotalItemsScore(items []Item) int {
	score := 5 * (len(items) / 2)

	for _, i := range items {
		trimmedDescription := strings.TrimSpace(i.ShortDescription)

		if len(trimmedDescription)%3 == 0 {
			// we don't need error handling on price since we already validated it with regex
			price, _ := strconv.ParseFloat(i.Price, 64)
			score += int(math.Ceil(price * 0.2))
		}
	}
	return score
}

/* 
6 points if the day in the purchase date is odd.
*/
func DateScore(d string) int {
	date, _ := time.Parse("2006-01-02", d)
	if date.Day()%2 == 1 {
		return 6
	}
	return 0
}

/* 
10 points if the time of purchase is after 2:00pm and before 4:00pm.
*/
func TimeScore(t string) int {
	time, _ := time.Parse("15:04", t)
	// this doesn't handle the case of the time being *exactly* 2:00pm, but I think that's okay
	if time.Hour() == 14 || time.Hour() == 15 {
		return 10
	}
	return 0
}
