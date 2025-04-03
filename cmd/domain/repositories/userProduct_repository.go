package repositories

import (
	"apisuario/cmd/domain/entities"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserProductRepository interface {
    ConnectUserProduct(data entities.UserProduct) error
    VerifyUserAndDevice(data entities.UserVeryfication) (*entities.User, *entities.Device, error) // Nuevo método
    CreateUser(data entities.UserVeryfication) error
}

type UserRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(db *gorm.DB) *UserRepositoryDB {
	return &UserRepositoryDB{db: db}
}

func (r *UserRepositoryDB) CreateUser(data entities.UserVeryfication) error {
	// Verificar si el usuario ya existe
	var existingUser struct {
		Username string
	}
	result := r.db.Table("users").
		Select("username").
		Where("username = ?", data.Username).
		Scan(&existingUser)
	if result.Error != nil {
		return fmt.Errorf("error al verificar usuario: %v", result.Error)
	}
	if result.RowsAffected > 0 {
		return fmt.Errorf("el usuario %s ya está registrado", data.Username)
	}

	// Encriptar contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al encriptar la contraseña: %v", err)
	}

	// Establecer tipo por defecto si no viene
	tipo := "user"
	if data.Tipo != "" {
		tipo = strings.ToLower(data.Tipo)
		if tipo != "user" && tipo != "admin" {
			return fmt.Errorf("tipo de usuario inválido, debe ser 'user' o 'admin'")
		}
	}

	// Crear usuario
	user := struct {
		Username string
		Password string
		Tipo     string
	}{data.Username, string(hashedPassword), tipo}
	
	if err := r.db.Table("users").Create(&user).Error; err != nil {
		return fmt.Errorf("error al crear el usuario: %v", err)
	}

	return nil
}

func (r *UserRepositoryDB) ConnectUserProduct(data entities.UserProduct) error {
    // 1. Verificar si el dispositivo existe
    var device struct {
        ID int
    }
    err := r.db.Table("devices").Where("id_esp32 = ?", data.IdEsp32).First(&device).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("el dispositivo con ID %s no existe", data.IdEsp32)
        }
        return fmt.Errorf("error al buscar dispositivo: %v", err)
    }

    // 2. Verificar si el usuario existe
    var user struct {
        ID   int
        Tipo string
    }
    err = r.db.Table("users").Where("username = ?", data.Username).First(&user).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return fmt.Errorf("el usuario %s no existe", data.Username)
        }
        return fmt.Errorf("error al buscar usuario: %v", err)
    }

    // 3. Verificar contraseña si se proporcionó
    if data.Password != "" {
        var dbUser struct {
            Password string
        }
        err = r.db.Table("users").Where("username = ?", data.Username).Select("password").First(&dbUser).Error
        if err != nil {
            return fmt.Errorf("error al verificar credenciales: %v", err)
        }

        err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(data.Password))
        if err != nil {
            return fmt.Errorf("contraseña incorrecta")
        }
    }

    // 4. Verificar si ya existe una relación para este usuario con CUALQUIER dispositivo
    var existingUserRelation struct {
        DeviceID int
    }
    err = r.db.Table("userdevices").
        Where("user_id = ?", user.ID).
        Select("device_id").
        First(&existingUserRelation).Error
    
    if err == nil {
        // El usuario ya tiene un dispositivo asignado
        if existingUserRelation.DeviceID == device.ID {
            return fmt.Errorf("el usuario ya está relacionado con este dispositivo")
        }
        return fmt.Errorf("el usuario ya tiene un dispositivo asignado (ID: %d)", existingUserRelation.DeviceID)
    } else if err != gorm.ErrRecordNotFound {
        return fmt.Errorf("error al verificar relaciones existentes del usuario: %v", err)
    }

    // 5. Verificar si el dispositivo ya está asignado a otro usuario
    var existingDeviceRelation struct {
        UserID int
    }
    err = r.db.Table("userdevices").
        Where("device_id = ?", device.ID).
        Select("user_id").
        First(&existingDeviceRelation).Error
    
    if err == nil {
        return fmt.Errorf("el dispositivo ya está asignado a otro usuario (ID: %d)", existingDeviceRelation.UserID)
    } else if err != gorm.ErrRecordNotFound {
        return fmt.Errorf("error al verificar relaciones existentes del dispositivo: %v", err)
    }

    // 6. Crear la nueva relación (ambas verificaciones pasaron)
    userDevice := struct {
        UserID   int
        DeviceID int
    }{user.ID, device.ID}

    if err := r.db.Table("userdevices").Create(&userDevice).Error; err != nil {
        return fmt.Errorf("error al crear la relación usuario-dispositivo: %v", err)
    }

    return nil
}

func (r *UserRepositoryDB) VerifyUserAndDevice(data entities.UserVeryfication) (*entities.User, *entities.Device, error) {
	// 1. Verificar usuario
	var user entities.User
	if err := r.db.Table("users").
		Where("username = ?", data.Username).
		First(&user).Error; err != nil {
		
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, nil, fmt.Errorf("error al verificar usuario")
	}

	// 2. Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return nil, nil, fmt.Errorf("contraseña incorrecta")
	}

	// 3. Obtener dispositivo asociado
	var device entities.Device
	if err := r.db.Table("userdevices").
		Joins("INNER JOIN devices ON devices.id = userdevices.device_id").
		Where("userdevices.user_id = ?", user.ID).
		Select("devices.*").
		First(&device).Error; err != nil {
		
		if err == gorm.ErrRecordNotFound {
			return nil, nil, fmt.Errorf("no hay dispositivos asociados")
		}
		return nil, nil, fmt.Errorf("error al buscar dispositivo")
	}

	return &user, &device, nil
}