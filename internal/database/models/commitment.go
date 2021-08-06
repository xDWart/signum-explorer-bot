package models

import (
	"gorm.io/gorm"
)

type AverageCommitment struct {
	gorm.Model
	AverageCommitment float64
}
