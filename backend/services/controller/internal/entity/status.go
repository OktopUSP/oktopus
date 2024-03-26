package entity

type Status uint8

const (
	Offline Status = iota
	Associating
	Online
)
