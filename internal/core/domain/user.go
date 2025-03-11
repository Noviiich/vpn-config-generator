package domain

import "time"

type User struct {
	ID     string
	Device []Device
	CreateAt time.Time
}
