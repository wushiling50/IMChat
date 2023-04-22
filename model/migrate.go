package model

func migration() {
	err := DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&User{},
			&News{},
			&Group{},
			&Friend{},
		)

	if err != nil {
		panic("数据库迁移失败")
	}

}
