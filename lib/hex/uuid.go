package hex

import "github.com/google/uuid"

func NewUUID() string {
	u := uuid.New()
	return u.String()
}
