package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `gorm:"primarykey" json:"id"`
	UserName string `gorm:"unique;type:varchar(100)" json:"user_name"`
	Password string `gorm:"type:varchar(100)" json:"password"`
	Role     string `gorm:"type:varchar(20)" json:"role"`
}
