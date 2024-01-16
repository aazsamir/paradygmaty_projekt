package main

import (
	"time"

	"golang.org/x/exp/slog"
)

type SlowDatabase struct {
	SimpleRepository
}

const DELAY = 1 * time.Second

func (db *SlowDatabase) GetUser(id int) (*User, error) {
	slog.Info("SDB::GetUser", "id", id)
	time.Sleep(DELAY)

	user, err := db.SimpleRepository.GetUser(id)

	return user, err
}

func (db *SlowDatabase) GetUsers() ([]User, error) {
	slog.Info("SDB::GetUsers")
	time.Sleep(DELAY)

	users, err := db.SimpleRepository.GetUsers()

	return users, err
}

func (db *SlowDatabase) GetUsersWithTags() ([]User, error) {
	slog.Info("SDB::GetUsersWithTags")
	time.Sleep(DELAY)

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

func (db *SlowDatabase) GetUserTags(user *User) ([]Tag, error) {
	slog.Info("SDB::GetUserTags", "user", user)
	time.Sleep(DELAY)

	tags, err := db.SimpleRepository.GetUserTags(user)

	if err != nil {
		return nil, err
	}

	return tags, nil
}
