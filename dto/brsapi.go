package dto

type BrsResponse struct {
	Gold           []BrsItem `json:"gold"`
	Currency       []BrsItem `json:"currency"`
	Cryptocurrency []BrsItem `json:"cryptocurrency"`
}

type BrsItem struct {
	Symbol string      `json:"symbol"`
	Name   string      `json:"name"`  // در خروجی API شما "name" است
	Price  interface{} `json:"price"` // قیمت ممکن است عدد یا رشته باشد
	Unit   string      `json:"unit"`
	Date   string      `json:"date"`
	Time   string      `json:"time"`
}
