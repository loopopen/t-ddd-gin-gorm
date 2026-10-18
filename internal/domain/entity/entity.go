package entity

type ID interface {
	uint | uint64 | string
}

type Entity[Tid ID] interface {
	Id() Tid

	Equals(other Entity[Tid]) bool

	Validate() error
}

type EntityConstraint[Tid ID, T any] interface {
	Entity[Tid]
	~*T
}
