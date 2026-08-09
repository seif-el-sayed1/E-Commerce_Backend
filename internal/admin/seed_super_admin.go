package admin

import (
	"errors"
	"log"

	"seif-el-sayed1/E-Commerce_Backend.git/internal/constants"

	"gorm.io/gorm"
)

func SeedSuperAdmin(db *gorm.DB, superAdminEmail, superAdminPassword string) error {
	var admin Admin

	res := db.Where("is_super_admin = ?", true).First(&admin)

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	}

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		admin = Admin{
			FirstName:    "Super",
			LastName:     "Admin",
			Role:         constants.SuperAdminRole,
			Email:        superAdminEmail,
			Password:     superAdminPassword,
			IsSuperAdmin: true,
			IsVerified:   true,
		}

		newAdmin := db.Create(&admin)
		if newAdmin.Error != nil {
			log.Println(constants.Error(newAdmin.Error))
			return newAdmin.Error
		}

		log.Println(constants.Success("✅ Super admin seeded successfully"))
	}

	return nil
}
