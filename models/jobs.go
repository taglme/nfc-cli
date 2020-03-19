package models

type GenericJobParams struct {
	Cmd       Command
	AdapterId string
	Repeat    int
	Expire    int
	Auth      []byte
	Export    bool
	JobName   string
}

type NewJob struct {
	JobName string `json:"job_name" binding:"required"`
}
