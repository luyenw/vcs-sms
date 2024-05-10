package entity

type Scope struct {
	ID   uint   `gorm:"primary_ky" gorm:"column:id"`
	Name string `gorm:"column:name" gorm:"unique"`
}
