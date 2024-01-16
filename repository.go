package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SimpleRepository interface {
	SaveUser(user *User) error
	SaveTag(tag *Tag) error
	DeleteUser(user *User) error
	DeleteTag(tag *Tag) error
	GetUser(id int) (*User, error)
	GetUsers() ([]User, error)
	GetUsersWithTags() ([]User, error)
	GetUserTags(user *User) ([]Tag, error)
	Close() error
}

type Database struct {
	db *sql.DB
}

type User struct {
	ID   int
	Name string
	Tags []Tag
}

type Tag struct {
	ID     int
	UserID int
	Name   string
}

func NewSQLiteRepository() (*Database, error) {
	const file string = "db.db"
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

func Migrate(db *Database) error {
	_, err := db.db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")

	if err != nil {
		return err
	}

	_, err = db.db.Exec("CREATE TABLE IF NOT EXISTS tags (id INTEGER PRIMARY KEY, user_id INTEGER, name TEXT)")

	if err != nil {
		return err
	}

	_, err = db.db.Exec("INSERT INTO users (name) VALUES ('John')")

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) SaveUser(user *User) error {
	_, err := db.db.Exec("INSERT INTO users (name) VALUES (?)", user.Name)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) SaveTag(tag *Tag) error {
	_, err := db.db.Exec("INSERT INTO tags (user_id, name) VALUES (?, ?)", tag.UserID, tag.Name)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteUser(user *User) error {
	_, err := db.db.Exec("DELETE FROM users WHERE id = ?", user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteTag(tag *Tag) error {
	_, err := db.db.Exec("DELETE FROM tags WHERE id = ?", tag.ID)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUser(id int) (*User, error) {
	row := db.db.QueryRow("SELECT id, name FROM users WHERE id = ?", id)

	var user User

	err := row.Scan(&user.ID, &user.Name)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *Database) GetUsers() ([]User, error) {
	rows, err := db.db.Query("SELECT id, name FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.Name)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (db *Database) GetUsersWithTags() ([]User, error) {
	users, err := db.GetUsers()

	if err != nil {
		return nil, err
	}

	for _, user := range users {
		tags, err := db.GetUserTags(&user)

		if err != nil {
			return nil, err
		}

		user.Tags = tags
	}

	return users, nil
}

func (db *Database) GetUserTags(user *User) ([]Tag, error) {
	rows, err := db.db.Query("SELECT id, user_id, name FROM tags WHERE user_id = ?", user.ID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var tag Tag

		err := rows.Scan(&tag.ID, &tag.UserID, &tag.Name)

		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
