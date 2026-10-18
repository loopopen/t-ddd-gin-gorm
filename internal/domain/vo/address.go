package vo

//go:generate go tool shoot new -getset -sign=@ -file=$GOFILE

// Address
// @getter
type Address struct {
	province string
	city     string
	street   string
	zipCode  string
}
