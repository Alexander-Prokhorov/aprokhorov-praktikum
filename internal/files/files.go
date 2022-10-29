package files

import (
	"encoding/json"
	"os"

	"aprokhorov-praktikum/internal/storage"
)

const defaultPerm = 0o644

// Save data dump in local file.
func SaveData(fileName string, s storage.Storage) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, defaultPerm)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	err = enc.Encode(s)

	if err != nil {
		return err
	}

	return nil
}

// Load data-dump from local file.
func LoadData(fileName string, s storage.Storage) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY, defaultPerm)
	if err == nil {
		defer file.Close()

		if err := json.NewDecoder(file).Decode(s); err != nil {
			return err
		}
	}

	return nil
}