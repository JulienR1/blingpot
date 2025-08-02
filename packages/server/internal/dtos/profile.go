package dtos

type Profile struct {
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	Picture   NullString `json:"picture"`
}
