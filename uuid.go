package filetrove

import "github.com/google/uuid"

// CreateUUID returns a UUID v4 as a string
func CreateUUID() (string, error) {
	newuuid, err := uuid.NewRandom()

	return newuuid.String(), err
}
