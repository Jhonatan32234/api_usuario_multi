package repositories

import (
	"apisuario/cmd/domain/entities"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserProductRepository interface {
	ConnectUserProduct(data entities.UserProduct) error
	VerifyUserDevice(data entities.UserVeryfication) (string, error)
}

type UserRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(db *gorm.DB) *UserRepositoryDB {
	return &UserRepositoryDB{db: db}
}
func (r *UserRepositoryDB) ConnectUserProduct(data entities.UserProduct) error {
	var existingRecord struct {
		Username string
		IdEsp32  string
	}
	result := r.db.Table("users").
		Select("username").
		Where("username = ?", data.Username).
		Scan(&existingRecord)
	if result.Error != nil {
		return fmt.Errorf("error al verificar usuario: %v", result.Error)
	}
	if result.RowsAffected > 0 {
		return fmt.Errorf("el usuario %s ya está registrado", data.Username)
	}

	result = r.db.Table("devices").
		Select("id_esp32").
		Where("id_esp32 = ?", data.IdEsp32).
		Scan(&existingRecord)
	if result.Error != nil {
		return fmt.Errorf("error al verificar dispositivo: %v", result.Error)
	}
	if result.RowsAffected > 0 {
		return fmt.Errorf("el dispositivo con ID %s ya está registrado", data.IdEsp32)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al encriptar la contraseña: %v", err)
	}

	user := struct {
		Username string
		Password string
	}{data.Username, string(hashedPassword)} // Guardar la contraseña encriptada
	if err := r.db.Table("users").Create(&user).Error; err != nil {
		return fmt.Errorf("error al crear el usuario: %v", err)
	}

	device := struct {
		IdEsp32 string
	}{data.IdEsp32}
	if err := r.db.Table("devices").Create(&device).Error; err != nil {
		return fmt.Errorf("error al crear el dispositivo: %v", err)
	}

	var userID, deviceID int
	r.db.Table("users").Where("username = ?", data.Username).Select("id").Scan(&userID)
	r.db.Table("devices").Where("id_esp32 = ?", data.IdEsp32).Select("id").Scan(&deviceID)

	userDevice := struct {
		UserID   int
		DeviceID int
	}{userID, deviceID}

	if err := r.db.Table("userdevices").Create(&userDevice).Error; err != nil {
		return fmt.Errorf("error al crear la relación usuario-dispositivo: %v", err)
	}

	return nil
}


func (r *UserRepositoryDB) VerifyUserDevice(data entities.UserVeryfication) (string, error) {
	var deviceID string
	var userData struct {
		ID       int
		Password string
	}

	log.Println("Verifying user and device...")

	result := r.db.Table("prueba.users").
		Where("username = ?", data.Username).
		Select("id, password").Scan(&userData)

	if result.Error != nil || result.RowsAffected == 0 {
		return "", fmt.Errorf("usuario o contraseña incorrectos")
	}

	err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(data.Password))
	if err != nil {
		return "", fmt.Errorf("usuario o contraseña incorrectos")
	}

	log.Print("User verified successfully")

	result = r.db.Table("prueba.userdevices").
		Joins("JOIN devices ON devices.id = userdevices.device_id").
		Where("userdevices.user_id = ?", userData.ID).
		Select("devices.id_esp32").Scan(&deviceID)

	if result.Error != nil || result.RowsAffected == 0 {
		return "", fmt.Errorf("no hay ningún dispositivo asociado al usuario")
	}

	return deviceID, nil
}