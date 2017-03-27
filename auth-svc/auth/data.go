package auth

// User represents an user and its credentials in the system.
type User struct {
	//gorm.Model
	ID       uint `gorm:"primary_key"`
	Mail     string
	Password string
	Salt     string
}
