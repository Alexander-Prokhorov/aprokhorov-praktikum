package files

import (
	"encoding/json"
	"os"

	"aprokhorov-praktikum/cmd/server/storage"
)

func SaveData(fileName string, s storage.Storage) error {

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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

func LoadData(fileName string, s storage.Storage) error {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err == nil {
		defer file.Close()
		if err := json.NewDecoder(file).Decode(s); err != nil {
			return err
		}
	}
	return nil
}
