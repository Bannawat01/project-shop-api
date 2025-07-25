package main

import (
	"log"
	"github.com/Bannawat101/project-shop-api/config"
	"github.com/Bannawat101/project-shop-api/databases"
	"github.com/Bannawat101/project-shop-api/server"
)

func main() {
	log.Println("Loading configuration...")
	conf := config.ConfigGetting()
	log.Println("Configuration loaded successfully")
	
	log.Println("Connecting to database...")
	db := databases.NewPostgresDatabase(conf.Database)
	log.Printf("Connected to database %s", conf.Database.DBName)
	
	log.Println("Creating server...")
	server := server.NewEchoServer(conf, db.ConnectionGetting())
	log.Println("Server created, starting...")

	server.Start()
}
