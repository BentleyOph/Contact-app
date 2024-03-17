package models

import (
	"github.com/BentleyOph/htmx-go/pkg/config"
	"github.com/jinzhu/gorm" // gorm is a Go library for dealing with SQL databases
)

var db *gorm.DB // db is a pointer to a gorm.DB object

type Contact struct {
	gorm.Model
	FirstName string `gorm:""` // Name is a string ,the gorm metadata is empty to indicate that the field should be created with the default settings
	LastName  string // LastName is a string
	Email     string // Email is a string
	Phone     string // Phone is a string
}

func init() {
	config.Connect()           // Connect to the database
	db = config.GetDB()        // Get the database object
	db.AutoMigrate(&Contact{}) // Migrate the schema
}

func (b *Contact) CreateContact() *Contact {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func GetAllContacts(searchTerm string) []Contact {
	var Contacts []Contact
	if searchTerm != "" {
		db.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%").Find(&Contacts)
	} else {
		db.Find(&Contacts)
	}
	return Contacts
}

func GetContactById(id int64) (*Contact, *gorm.DB) {
	var contact Contact
	db.Where("ID = ?", id).Find(&contact)
	return &contact, db
}

func DeleteContact(id int64) {
	var contact Contact
	db.Where("ID = ?", id).Delete(&contact)
}
//
func UpdateContact(id int64, contact Contact) {
	db.Where("ID = ?", id).Save(&contact)
}


func GetContactsCount() int {
	var contacts []Contact
	var count int
	db.Find(&contacts).Count(&count)
	return count
}
