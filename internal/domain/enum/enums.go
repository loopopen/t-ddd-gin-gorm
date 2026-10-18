package enum

//go:generate go tool shoot enum -json -text -sql -gorm -file=$GOFILE

type Gender int32

const (
	GenderUnknown Gender = iota
	GenderMale
	GenderFemale
)
