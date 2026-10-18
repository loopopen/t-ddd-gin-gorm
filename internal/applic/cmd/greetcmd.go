package cmd

//go:generate go tool shoot new -json -file=$GOFILE

type GreetCmd struct {
	Name string
}
