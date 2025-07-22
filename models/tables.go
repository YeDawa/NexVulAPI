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

type Users struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"unique"`
	Email     string    `gorm:"unique;not null;size:100"`
	Name      string    `gorm:"size:255"`
	Password  string    `gorm:"not null"`
	ApiKey    string    `gorm:"unique;not null;size:100"`
	Plan      PlanType  `gorm:"type:varchar(20);default:'free'"`
	Status    Status    `gorm:"type:varchar(20);default:'unchecked'"`
	Salt      string    `gorm:"unique;not null"`
	Timezone  string    `gorm:"default:'UTC'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Scans struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	UserId    uint      `gorm:"not null"`
	Slug      string    `gorm:"unique;not null;size:100"`
	Urls      string    `gorm:"not null;size:255"`
	Data      string    `gorm:"type:text"`
	Public    bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type CustomHeaders struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	UserId    uint      `gorm:"not null"`
	Header    string    `gorm:"not null;size:100"`
	Value     string    `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Reports struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	UserId    uint      `gorm:"not null"`
	FileName  string    `gorm:"unique;not null;size:255"`
	Slug      string    `gorm:"unique;not null;size:100"`
	ScanId    uint      `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type ScansXSS struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	UserId    uint      `gorm:"not null"`
	ScanId    uint      `gorm:"not null"`
	Url       string    `gorm:"not null;size:255"`
	Result    string    `gorm:"not null;size:255"`
	Parameter string    `gorm:"not null;size:100"`
	Payload   string    `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Payloads struct {
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null;size:100"`
	FileName  string    `gorm:"unique;not null;size:255"`
	Content   string    `gorm:"type:text"`
	Type      string    `gorm:"not null;size:50"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Profile struct {
	Id         uint   `gorm:"primaryKey,autoIncrement"`
	UserId     uint   `gorm:"not null"`
	Bio        string `json:"bio"`
	Location   string `json:"location"`
	Website    string `json:"website"`
	Contact    string `json:"contact"`
	PublicName string `json:"public_name"`
	Twitter    string `json:"twitter"`
	Github     string `json:"github"`
	Linkedin   string `json:"linkedin"`
}

type CustomWordlists struct {
	Id         uint      `gorm:"primaryKey;autoIncrement"`
	Slug       string    `gorm:"unique;not null;size:100"`
	UserId     uint      `gorm:"not null"`
	FileName   string    `gorm:"unique;not null;size:255"`
	TotalWords int       `gorm:"not null"`
	Type       string    `gorm:"not null;size:50"`
	Domain     string    `gorm:"size:255"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (Users) TableName() string {
	return "hs_users"
}

func (Scans) TableName() string {
	return "hs_scans"
}

func (Reports) TableName() string {
	return "hs_reports"
}

func (CustomHeaders) TableName() string {
	return "hs_custom_headers"
}

func (ScansXSS) TableName() string {
	return "hs_scans_xss"
}

func (Payloads) TableName() string {
	return "hs_payloads"
}

func (Profile) TableName() string {
	return "hs_profile"
}

func (CustomWordlists) TableName() string {
	return "hs_custom_wordlists"
}
