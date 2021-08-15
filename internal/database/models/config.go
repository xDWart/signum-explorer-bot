package models

import (
	"gorm.io/gorm"
)

type Config struct {
	gorm.Model
	Name   string `gorm:"uniqueIndex"`
	ValueS string
	ValueF float64
	ValueI int
}
