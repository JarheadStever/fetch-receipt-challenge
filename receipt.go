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

type Validator func(value string) error

type Validators []struct {
	input    string
	function Validator
}

func validatePattern(pattern *regexp.Regexp, fieldName string) Validator {
	return func(value string) error {
		if !pattern.MatchString(value) {
			return fmt.Errorf("invalid %s: %s", fieldName, value)
		}
		return nil
	}
}

func validateDate(fieldName string) Validator {
	return func(value string) error {
		if _, err := time.Parse("2006-01-02", value); err != nil {
			return fmt.Errorf("invalid %s: %s", fieldName, value)
		}
		return nil
	}
}

func validateTime(fieldName string) Validator {
	return func(value string) error {
		if _, err := time.Parse("15:04", value); err != nil {
			return fmt.Errorf("invalid %s: %s", fieldName, value)
		}
		return nil
	}
}

func combineErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var sb strings.Builder
	for _, err := range errs {
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return errors.New(sb.String())
}

func checkValidators(validators Validators) []error {
	var validationErrors []error
	for _, v := range validators {
		if err := v.function(v.input); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}
	return validationErrors
}

var retailerRegex = regexp.MustCompile(`^[\w\s\-&]+$`)
var totalRegex = regexp.MustCompile(`^\d+\.\d{2}$`)

/*
Validate each part of a receipt by matching regexes and time/date values to the API spec
*/
func (r *Receipt) Validate() error {

	validators := Validators{
		{r.Retailer, validatePattern(retailerRegex, "retailer")},
		{r.PurchaseDate, validateDate("purchaseDate")},
		{r.PurchaseTime, validateTime("purchaseTime")},
		{r.Total, validatePattern(totalRegex, "total")},
	}

	validationErrors := checkValidators(validators)

	if len(r.Items) < 1 {
		validationErrors = append(validationErrors, errors.New("receipt needs at least one item"))
	}

	for _, item := range r.Items {
		if err := item.Validate(); err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("invalid item [%v]. Error: %s", item.ShortDescription, err))
		}
	}

	return combineErrors(validationErrors)
}

var descriptionRegex = regexp.MustCompile(`^[\w\s\-]+$`)
var priceRegex = regexp.MustCompile(`^\d+\.\d{2}$`)

/*
Validate an item with in a receipt by matching regexes to the API spec
*/
func (i *Item) Validate() error {

	validators := Validators{
		{i.ShortDescription, validatePattern(descriptionRegex, "shortDescription")},
		{i.Price, validatePattern(priceRegex, "price")},
	}

	return combineErrors(checkValidators(validators))
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
