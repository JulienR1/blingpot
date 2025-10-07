package dtos

type Profile struct {
	Sub       string     `json:"sub"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Email     string     `json:"email"`
	Picture   NullString `json:"picture"`
}
