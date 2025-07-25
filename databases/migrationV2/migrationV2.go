package main

import (
	"github.com/Bannawat101/project-shop-api/config"
	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/entities"
	"gorm.io/gorm"
)

func main() {
	conf := config.ConfigGetting()
	db := databases.NewPostgresDatabase(conf.Database)

	tx := db.Connect().Begin()

	itemAdding(tx)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic("MigrationV2 failed: " + err.Error())
	} else {
		println("MigrationV2 completed successfully!")
	}
}

func itemAdding(tx *gorm.DB) {
	items := []entities.Item{
		{
			Name:        "Sword",
			Description: "A sword that can be used to fight enemies.",
			Price:       100,
			Picture:     "https://i.pinimg.com/736x/73/cc/79/73cc79391b764ec40a5c77052bb846b9.jpg",
		},
		{
			Name:        "Shield",
			Description: "A shield that can be used to block enemy attacks.",
			Price:       50,
			Picture:     "https://i.pinimg.com/736x/fe/83/27/fe832717d33f05c2dbd845809ce877b8.jpg",
		},
		{
			Name:        "Potion",
			Description: "A potion that can be used to heal wounds.",
			Price:       30,
			Picture:     "https://i.pinimg.com/564x/14/7e/7d/147e7d58fa2becce0045f3aadf1808b1.jpg",
		},
		{
			Name:        "Bow",
			Description: "A bow that can be used to shoot enemies from afar.",
			Price:       80,
			Picture:     "https://i.pinimg.com/564x/1f/91/72/1f9172f5bc27094c4e167e55f8cce2f2.jpg",
		},
		{
			Name:        "Arrow",
			Description: "An arrow that can be used with a bow to shoot enemies from afar.",
			Price:       10,
			Picture:     "https://i.pinimg.com/564x/3f/25/84/3f25842cb4a8ad53a19575cc3d25c844.jpg",
		},
	}
	tx.CreateInBatches(items, len(items))
}
