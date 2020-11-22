package model

import "time"

type Image struct {
	Id        int
	Url       string
	MD5       string
	AddressId int
	CreatedAt time.Time
}

func (*Image) TableName() string {
	return "spd_image"
}
