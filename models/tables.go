package models

import (
	"time"
)

type PlanType string

type Privacy string

type Status string

const (
	Free       PlanType = "free"
	Premium    PlanType = "premium"
	Enterprise PlanType = "enterprise"
)

const (
	Unchecked Status = "unchecked"
	Checked   Status = "checked"
	Verifed   Status = "verifed"
	Banned    Status = "banned"
)

type User struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"unique"`
	Email     string `gorm:"unique;not null;size:100"`
	Name      string `gorm:"size:255"`
	Password  string `gorm:"not null"`
	Plan      string `gorm:"type:varchar(20);default:'free'"`
	Status    Status `gorm:"type:varchar(20);default:'unchecked'"`
	Salt      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Scans struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	Title     string    `gorm:"not null;size:100"`
	UserId    uint      `gorm:"not null"`
	Slug      string    `gorm:"unique;not null;size:100"`
	Urls      string    `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Exports struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	UserId    uint      `gorm:"not null"`
	FileName  string    `gorm:"unique;not null;size:255"`
	Slug      string    `gorm:"unique;not null;size:100"`
	ScanId    uint      `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "hs_users"
}

func (Scans) TableName() string {
	return "hs_scans"
}

func (Exports) TableName() string {
	return "hs_exports"
}
