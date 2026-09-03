package migrations

import (
	"github.com/MohammadrezaNadirkhanloo/Go-project/config"
	"github.com/MohammadrezaNadirkhanloo/Go-project/constans"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/db"
	"github.com/MohammadrezaNadirkhanloo/Go-project/data/models"
	"github.com/MohammadrezaNadirkhanloo/Go-project/pkg/logging"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var logger = logging.NewLogger(config.GetConfig())

func Up_1() {
	database := db.GetDB()
	CreateTable(database)
	createDefaultInformation(database)
}

func CreateTable(database *gorm.DB) { // ساخت جداول
	tables := []interface{}{}

	country := models.Country{}
	city := models.City{}
	user := models.User{}
	role := models.Role{}
	userRole := models.UserRole{}

	tables = newFunction(database, country, tables)
	tables = newFunction(database, city, tables)
	tables = newFunction(database, user, tables)
	tables = newFunction(database, role, tables)
	tables = newFunction(database, userRole, tables)

	database.Migrator().CreateTable(tables...)
	logger.Info(logging.Postgres, logging.Migration, "tables Created", nil)
}

func newFunction(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func createDefaultInformation(database *gorm.DB) { // ساخت دیتا برای جدول
	/////////role///////
	adminRole := models.Role{Name: constans.AdminRoleName}
	createdDataRole(database, &adminRole)
	defaultRole := models.Role{Name: constans.DefaultRoleName}
	createdDataRole(database, &defaultRole)
	////////user admin///////
	u := models.User{Username: constans.DefaultUserName, Firstname: "develop", Lastname: "khanloo", Email: "khanl@gmail.co", MobileNumber: "09111111111"}
	pass := "123456"
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.Password = string(hashPassword)

	createdDataAdminUser(database, &u, adminRole.Id)
}

func createdDataRole(database *gorm.DB, r *models.Role) {
	exists := 0
	database.
		Model(&models.Role{}). // برو جدول سطح دسترسی
		Select("1").
		Where("name = ?", r.Name). // بررسی کن ببین همچین کلمه ای هست
		First(&exists)             // اگر هست ۱ بریز داخلش اگه نیست 0 میمونه
	if exists == 0 {
		database.Create(r)
	}
}
func createdDataAdminUser(database *gorm.DB, u *models.User, roleId int) {
	exists := 0
	database.
		Model(&models.User{}).
		Select("1").
		Where("username = ?", u.Username).
		First(&exists)
	if exists == 0 {
		database.Create(u)
		ur := models.UserRole{UserId: u.Id, RoleId: roleId}
		database.Create(&ur)
	}
}

func createCountry(database *gorm.DB) {
	count := 0
	database.
		Model(&models.Country{}).
		Select("count(*)").
		Find(&count)
	if count == 0 {
		database.Create(&models.Country{Name: "Iran", Cities: []models.City{
			{Name: "Tehran"},
			{Name: "Isfahan"},
			{Name: "Shiraz"},
			{Name: "Chalus"},
			{Name: "Ahwaz"},
		}})
		database.Create(&models.Country{Name: "USA", Cities: []models.City{
			{Name: "New York"},
			{Name: "Washington"},
		}})
		database.Create(&models.Country{Name: "Germany", Cities: []models.City{
			{Name: "Berlin"},
			{Name: "Munich"},
		}})
		database.Create(&models.Country{Name: "China", Cities: []models.City{
			{Name: "Beijing"},
			{Name: "Shanghai"},
		}})
		database.Create(&models.Country{Name: "Italy", Cities: []models.City{
			{Name: "Roma"},
			{Name: "Turin"},
		}})
		database.Create(&models.Country{Name: "France", Cities: []models.City{
			{Name: "Paris"},
			{Name: "Lyon"},
		}})
		database.Create(&models.Country{Name: "Japan", Cities: []models.City{
			{Name: "Tokyo"},
			{Name: "Kyoto"},
		}})
		database.Create(&models.Country{Name: "South Korea", Cities: []models.City{
			{Name: "Seoul"},
			{Name: "Ulsan"},
		}})
	}
}
