package entity

type RegisteredEmail struct {
	Email string `gorm:"primary_key" gorm:"column:email" gorm:"unique" gorm:"not null" gorm:"type:varchar(255)"`
}
