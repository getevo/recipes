package myapp

import (
	"github.com/getevo/evo-ng/lib/generic"
	"gorm.io/gorm"
)

type MyModel struct {
	Name generic.Value
	LastName string
	ID int `gorm:"primaryKey"`
	gorm.Model
}
