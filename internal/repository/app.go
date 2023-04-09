package repository

import "gorm.io/gorm"

type App interface{}

type app struct {
	db *gorm.DB
}

func NewApp(db *gorm.DB) App {
	return &app{
		db: db,
	}
}
