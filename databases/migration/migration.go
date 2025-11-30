package main

import (
	"github.com/Bannawat101/project-shop-api/config"
	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/entities"
	"gorm.io/gorm"
)

func main() {
	conf := config.ConfigGetting()                     //โหลดค่าคอนฟิกจากไฟล์ yml
	db := databases.NewPostgresDatabase(conf.Database) //สร้างการเชื่อมต่อฐานข้อมูล Postgres

	tx := db.Connect().Begin() //เริ่มต้น transaction เพื่อให้การมิเกรชันมีความปลอดภัยมากขึ้น

	playerMigration(tx) //เรียกใช้ฟังก์ชันมิเกรชันสำหรับแต่ละตาราง
	adminMigration(tx)
	itemMigration(tx)
	playerCoinMigration(tx)
	inventoryMigration(tx)
	purchaseHistoryMigration(tx)

	if err := tx.Commit().Error; err != nil { //ตรวจสอบข้อผิดพลาดในการคอมมิต transaction
		tx.Rollback()
		panic("Migration failed: " + err.Error())
	} else {
		println("Migration completed successfully!")
	}
}

func playerMigration(tx *gorm.DB) { //สร้างตาราง Player
	tx.Migrator().CreateTable(&entities.Player{}) //ใช้ GORM สร้างตารางจาก struct Player ในแพ็กเกจ entities
}
func adminMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.Admin{})
}
func itemMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.Item{})
}

func playerCoinMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.PlayerCoin{})
}
func inventoryMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.Inventory{})
}
func purchaseHistoryMigration(tx *gorm.DB) {
	tx.Migrator().CreateTable(&entities.PurchaseHistory{})
}
