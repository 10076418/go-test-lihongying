package main

type UserType struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type EandPType struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NandEType struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type NameType struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

const (
	OK           = 1
	DB_ERROR     = 2
	FORMAT_ERROR = 3
)
