package entity

type User struct {
	ID       string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email    string `gorm:"type:text;unique;not null" json:"email"`
	Password string `gorm:"type:text;not null" json:"password"`
}

func (User) TableName() string {
	return "users"
}
