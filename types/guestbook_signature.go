package types

import "time"

type GuestbookSignature struct {
	Id         int
	Name       string
	IsApproved bool
	CreatedAt  time.Time
}
