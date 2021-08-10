package models

import (
	"gorm.io/gorm"
)

type Price struct {
	gorm.Model
	SignaPrice float64
	BtcPrice   float64
}
