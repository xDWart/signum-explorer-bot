package models

import (
	"gorm.io/gorm"
)

type NetworkInfo struct {
	gorm.Model
	AverageCommitment float64
	NetworkDifficulty float64
}
