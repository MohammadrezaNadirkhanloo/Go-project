package models

type User struct {
	BaseModel
	Username     string `gorm:"type:string;size:20;not null;unique"`
	Firstname    string `gorm:"type:string;size:20;null;"`
	Lastname     string `gorm:"type:string;size:20;null;"`
	MobileNumber string `gorm:"type:string;size:12;null;unique;default:null"`
	Email        string `gorm:"type:string;size:64;null;unique;default:null"`
	Password     string `gorm:"type:string;size:64;not null"`
	Enabled      bool   `gorm:"default:true"`
}

