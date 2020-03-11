package models

type GenericJobParams struct {
	Cmd       Command
	AdapterId string
	Repeat    int
	Expire    int
	Auth      []byte
}
