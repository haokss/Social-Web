package model

func migration() {
	DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&User{}).
		AutoMigrate(&Task{}).
		AutoMigrate(&TimingTask{}).
		AutoMigrate(&Image{})
	DB.Model(&Task{}).AddForeignKey("uid", "User(id)", "CASCADE", "CASCADE")
}
