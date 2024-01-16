package main

import (
	"golang.org/x/exp/slog"
)

type CuratorDatabase struct {
	SimpleRepository
	curator *Prolog
}

func (db *CuratorDatabase) GetUser(id int) (*User, error) {
	slog.Info("CDB::GetUser", "id", id)
	user, err := db.SimpleRepository.GetUser(id)

	return user, err
}

func (db *CuratorDatabase) GetUsers() ([]User, error) {
	slog.Info("CDB::GetUsers")
	users, err := db.SimpleRepository.GetUsers()

	return users, err
}

func (db *CuratorDatabase) GetUsersWithTags() ([]User, error) {
	slog.Info("CDB::GetUsersWithTags")
	users, err := db.GetUsers()

	if err != nil {
		return nil, err
	}

	for i := range users {
		user := &users[i]
		tags, err := db.GetUserTags(user)

		if err != nil {
			return nil, err
		}

		user.Tags = tags

		slog.Info("curator", "user", user, "result", db.curator.Query(user.ID), "tags", tags)

		if len(tags) > 0 && db.curator.Query(user.ID) {
			db.curator.save(user.ID)
		}

		slog.Info("user", "user", user)
	}

	slog.Info("users", "users", users)

	return users, nil
}

func (db *CuratorDatabase) GetUserTags(user *User) ([]Tag, error) {
	slog.Info("CDB::GetUserTags", "user", user)
	if !db.curator.Query(user.ID) {
		return nil, nil
	}

	tags, err := db.SimpleRepository.GetUserTags(user)

	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (db *CuratorDatabase) SaveTag(tag *Tag) error {
	slog.Info("CDB::SaveTag", "tag", tag)
	db.curator.save(tag.UserID)

	err := db.SimpleRepository.SaveTag(tag)

	if err != nil {
		return err
	}

	return nil
}

func (db *CuratorDatabase) DeleteUser(user *User) error {
	slog.Info("CDB::DeleteUser", "user", user)
	db.curator.delete(user.ID)

	err := db.SimpleRepository.DeleteUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (db *CuratorDatabase) DeleteTag(tag *Tag) error {
	slog.Info("CDB::DeleteTag", "tag", tag)
	db.curator.delete(tag.UserID)

	err := db.SimpleRepository.DeleteTag(tag)

	if err != nil {
		return err
	}

	return nil
}
