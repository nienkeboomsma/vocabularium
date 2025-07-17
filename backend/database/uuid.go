package database

import "github.com/google/uuid"

func StringToUUID(s string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(s))
}
