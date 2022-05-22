/*
Copyright © 2022 JEREMY PHUA <jeremyphuachengtoon@gmail.com>
*/
package db

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

var infoBucket = []byte("info")
var db *bolt.DB

type Information struct {
	Site     string // Website or Application
	UserInfo User   // Username and Password
}

type User struct {
	Username string
	Password string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(infoBucket)
		return err
	})
}

// Add username/password
func AddCredentials(site string, username string, password string) error {
	var info Information
	info.Site = site
	info.UserInfo.Username = username
	info.UserInfo.Password = password
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(infoBucket)

		encoded, err := json.Marshal(info.UserInfo)
		if err != nil {
			return err
		}

		return b.Put([]byte(info.Site), encoded)
	})

	return err
}

// List all username/password
func AllCredentials() ([]Information, error) {
	var infos []Information
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(infoBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var user User
			json.Unmarshal([]byte(v), &user)
			infos = append(infos, Information{
				Site: string(k),
				UserInfo: User{
					Username: user.Username,
					Password: user.Password,
				},
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return infos, nil
}

// Delete username/password
func DeleteCredentials(key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(infoBucket)
		return b.Delete([]byte(key))
	})
}
