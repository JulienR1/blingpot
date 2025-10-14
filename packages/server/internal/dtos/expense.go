package dtos

type Expense struct {
	Id         int      `json:"id"`
	SpenderId  string   `json:"spenderId"`
	Label      string   `json:"label"`
	Amount     int      `json:"amount"`
	Timestamp  UnixTime `json:"timestamp"`
	AuthorId   string   `json:"authorId"`
	CategoryId int      `json:"categoryId"`
}

type ExpensesSummary struct {
	Total     int         `json:"total"`
	Subtotals map[int]int `json:"categories"`
}
