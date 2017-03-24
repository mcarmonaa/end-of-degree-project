package auth

import (
	"hash/fnv"
)

type User struct {
	ID       uint32 `gorm:"primary_key"`
	Mail     string
	Password string
	Salt     string
}

func getID(mail string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(mail))
	return h.Sum32()
}
