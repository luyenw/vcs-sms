package entity

type Role struct {
	ID     uint    `gorm:"primary_key" gorm:"column:id"`
	Name   string  `gorm:"column:name"`
	Scopes []Scope `gorm:"many2many:role_scope"`
}
