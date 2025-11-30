package databases

import (
	"fmt"
	"log"
	"sync"

	"github.com/Bannawat101/project-shop-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct { //นี่คือ struct ที่เก็บการเชื่อมต่อฐานข้อมูล
	*gorm.DB //ฝัง gorm.DB เข้าไปใน struct นี้
}

var ( //ตรงนี้ใช้สำหรับ singleton pattern
	postgresDatabaseInstace *postgresDatabase //เก็บ instance ของฐานข้อมูล เพื่อให้มีแค่ instance เดียวตลอดการทำงาน
	once                    sync.Once         //ใช้ sync.Once เพื่อให้มั่นใจว่าโค้ดภายในจะถูกเรียกใช้แค่ครั้งเดียว
)

func NewPostgresDatabase(conf *config.Database) Database { //ฟังก์ชันนี้ใช้สำหรับสร้างการเชื่อมต่อฐานข้อมูล
	once.Do(func() { //ใช้ sync.Once เพื่อให้โค้ดภายในทำงานแค่ครั้งเดียว ไม่ว่าจะเรียกฟังก์ชันนี้กี่ครั้ง
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s search_path=%s",
			conf.Host,
			conf.User,
			conf.Password,
			conf.DBName,
			conf.Port,
			conf.SSLMode,
			conf.Schema,
		)

		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) //สร้างการเชื่อมต่อฐานข้อมูลด้วย GORM
		if err != nil {                                            //ตรวจสอบข้อผิดพลาดในการเชื่อมต่อ
			panic(err)
		}

		log.Printf("Connected to database %s", conf.DBName)

		postgresDatabaseInstace = &postgresDatabase{conn} //สร้าง instance ของ postgresDatabase และเก็บการเชื่อมต่อฐานข้อมูลไว้
	})

	return postgresDatabaseInstace
}

func (db *postgresDatabase) Connect() *gorm.DB {
	return postgresDatabaseInstace.DB
}
