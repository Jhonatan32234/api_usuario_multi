package core

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// Configura tu conexión a MySQL (reemplaza con tus credenciales)
    // URL de conexión usando variables de entorno
    dsn := "root:root@tcp(mysql:3306)/prueba?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Ejecutar cada consulta por separado
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		);
	`).Error
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS devices (
			id INT AUTO_INCREMENT PRIMARY KEY,
			id_esp32 VARCHAR(100) UNIQUE NOT NULL
		);
	`).Error
	if err != nil {
		log.Fatal("Error creating devices table:", err)
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS userdevices (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT,
			device_id INT,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
		);
	`).Error
	if err != nil {
		log.Fatal("Error creating userdevices table:", err)
	}

	return db
}
