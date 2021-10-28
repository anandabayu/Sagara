/*
 * File: Product.go
 * Project: models
 * File Created: Thursday, 28th October 2021 9:55:46 am
 * Author: Ananda Yudhistira (anandabayu12@gmail.com)
 * -----
 * Last Modified: Thursday, 28th October 2021 11:03:17 am
 * Modified By: Ananda Yudhistira (anandabayu12@gmail.com>)
 * -----
 * Copyright 2021 Ananda Yudhistira, -
 */

package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Image     string    `gorm:"type:text;not null;" json:"image"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Product) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Image = html.EscapeString(strings.TrimSpace(p.Image))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Product) Validate() error {

	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Content == "" {
		return errors.New("Required Content")
	}
	if p.Image == "" {
		return errors.New("Required Image")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (p *Product) SaveProduct(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Create(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	Products := []Product{}
	err = db.Debug().Model(&Product{}).Limit(100).Find(&Products).Error
	if err != nil {
		return &[]Product{}, err
	}
	if len(Products) > 0 {
		for i, _ := range Products {
			err := db.Debug().Model(&User{}).Where("id = ?", Products[i].AuthorID).Take(&Products[i].Author).Error
			if err != nil {
				return &[]Product{}, err
			}
		}
	}
	return &Products, nil
}

func (p *Product) FindProductByID(db *gorm.DB, pid uint64) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) UpdateAProduct(db *gorm.DB) (*Product, error) {

	var err error

	err = db.Debug().Model(&Product{}).Where("id = ?", p.ID).Updates(Product{Title: p.Title, Content: p.Content, Image: p.Image, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Product{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Product{}, err
		}
	}
	return p, nil
}

func (p *Product) DeleteAProduct(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Product{}).Where("id = ? and author_id = ?", pid, uid).Take(&Product{}).Delete(&Product{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Product not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
