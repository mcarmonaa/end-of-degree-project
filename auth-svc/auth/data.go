package auth

type User struct {
	//gorm.Model
	ID       uint `gorm:"primary_key"`
	Mail     string
	Password string
	Salt     string
}
