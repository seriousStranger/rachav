package database

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	DB_FILE_PERM = 0o600
	DB_DIR_PERM  = 0o750
)

func Save(userList map[string]string) error {
	jsonUserList, err := json.Marshal(userList)
	if err != nil {
		return fmt.Errorf("can't marshal user list: %w", err)
	}

	err = os.WriteFile(`./db/users.json`, jsonUserList, DB_FILE_PERM)
	if err != nil {
		return fmt.Errorf("can't write db: %w", err)
	}

	return nil
}

func Load() (map[string]string, error) {
	jsonUserList, err := os.ReadFile(`./db/users.json`)
	if err != nil {
		return nil, fmt.Errorf("can't read db: %w", err)
	}

	var UserList map[string]string

	err = json.Unmarshal(jsonUserList, &UserList)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal db: %w", err)
	}

	return UserList, nil
}

func CreateDbIfNotExist() error {
	err := os.MkdirAll("./db", DB_DIR_PERM)
	if err != nil {
		return fmt.Errorf("can't create db directory: %w", err)
	}

	_, err = os.Stat("./db/users.json")
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("can't check db file: %w", err)
	}

	emptyDb := make(map[string]string)

	err = Save(emptyDb)
	if err != nil {
		return fmt.Errorf("cant save empty database: %w", err)
	}

	return nil
}
