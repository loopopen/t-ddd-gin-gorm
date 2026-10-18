package event

type UserCreated struct {
	UserID uint
	Name   string
}

func (f *UserCreated) Topic() string {
	return "user.created"
}
