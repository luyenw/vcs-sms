package entity

type User struct {
	ID       uint    `gorm:"primary_key" gorm:"column:id"`
	Username string  `gorm:"column:username"`
	Password string  `gorm:"column:password"`
	Scopes   []Scope `gorm:"many2many:permissions"`
}
