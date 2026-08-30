package models

type Country struct {
	BaseModel
	Name  string `gorm:"size:10;type:string;not null"`
	Citys *[]City
}

type City struct {
	BaseModel
	Name      string `gorm:"size:10;type:string;not null"`
	CountryId int
	Country   Country `gorm:"foreignKey:CountryId"`
}
