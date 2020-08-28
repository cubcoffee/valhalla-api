package dao

import (
	"log"
)

// Credential model for the database Credentials
type Credential struct {
	ID   int64
	Hash string
	Salt string
}

// AddCredential insert new credential to the database
func AddCredential(credential Credential) (Credential, error) {
	db, err := InitDb()
	if err != nil {
		log.Print(err)
		return credential, err
	}
	defer db.Close()
	err = db.Create(&credential).Error
	if err != nil {
		log.Printf("Error adding credentials: %v", err)
		return credential, err
	}
	return credential, err
}

// DeleteCredential remove credential from the database
func DeleteCredential(credential Credential) error {
	db, err := InitDb()
	if err != nil {
		log.Print(err)
		return err
	}
	defer db.Close()
	return db.Delete(&credential).Error
}

// FindCredentialByID find a credential by its ID
func FindCredentialByID(credentialID int64) (Credential, error) {
	credential := Credential{}
	db, err := InitDb()
	defer db.Close()
	if err != nil {
		log.Print(err)
		return credential, err
	}
	db.First(&credential, credentialID)
	return credential, db.Error
}
