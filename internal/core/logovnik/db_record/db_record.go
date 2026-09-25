package db_record

import "github.com/google/uuid"

type Item struct {
	Id       uuid.UUID
	Program  string
	Desc     string
	Login    string
	Password string
}
