package models

type Country struct {
	BaseModel
	Name  string `gorm:"size:10;type:string;not null"`
	Citys *[]City
}
