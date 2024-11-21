package main

type Receipt struct {

    Retailer string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items string `json:"items"`
    Total float32 `json:"total"`
}
