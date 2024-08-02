package entity

type User struct {
	ID       uint   `gorm:"primary_key" gorm:"column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	RoleID   uint   `gorm:"column:role_id"`
	Role     Role
}
