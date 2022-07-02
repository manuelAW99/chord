package chord

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

type DataBasePlatform interface {
	GetByName([]byte) ([]byte, error)
	//GetByFun(string) ([][]byte, error)
	GetAll() ([][]byte, error)
	Set([]byte, []byte) error
	Update([]byte, []byte) error
	Delete([]byte) error
}

// Base de Datos para la plataforma
type DataBasePl struct {
	fileName string
}

// Cada informacion es un par key-data
type RowData struct {
	Key  []byte
	Data []byte
}

func NewDataBase(fileName string) *DataBasePl {
	db := &DataBasePl{
		fileName: fileName,
	}

	file, err := os.OpenFile(db.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer file.Close()

	return db
}

func (db *DataBasePl) readAll() ([]RowData, error) {
	file, err := os.Open(db.fileName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return []RowData{}, nil
	}

	var rows []RowData
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&rows); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return rows, nil
}

func (db *DataBasePl) writeAll(rows []RowData) error {
	file, err := os.Create(db.fileName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err := enc.Encode(rows); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (db *DataBasePl) GetByName(key []byte) ([]byte, error) {
	rows, err := db.readAll()
	if err != nil {
		return make([]byte, 0), err
	}

	for _, elem := range rows {
		if bytes.Equal(elem.Key, key) {
			return elem.Data, nil
		}
	}

	return make([]byte, 0), StorageError{message: "There is no agent with that name"}
}

/*
func (db *DataBasePl) GetByFun(fun string) ([][]byte, error) {
	rows, err := db.readAll()
	if err != nil {
		return make([][]byte, 0), err
	}

	data := make([]string, 0)
	for _, elem := range rows {
		if SearchString(elem.Data, fun) != -1 {
			data = append(data, elem.Data)
		}
	}
	if len(data) > 0 {
		return data, StorageError{message: "There is no agents with that function"}
	}
	return data, nil
}
*/

func (db *DataBasePl) GetAll() ([][]byte, error) {
	rows, err := db.readAll()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, StorageError{message: "No agents"}
	}

	var data [][]byte = make([][]byte, len(rows))

	for i, elem := range rows {
		data[i] = elem.Data
	}

	return data, nil
}

func (db *DataBasePl) Set(vKey []byte, vData []byte) error {
	rows, err := db.readAll()
	if err != nil {
		return err
	}

	for _, elem := range rows {
		if bytes.Equal(elem.Key, vKey) {
			return StorageError{message: "There is another agent with that name"}
		}
	}

	var newRows []RowData = append(rows, RowData{Key: vKey, Data: vData})
	err = db.writeAll(newRows)
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBasePl) Update(vKey []byte, vData []byte) error {
	rows, err := db.readAll()
	if err != nil {
		return err
	}

	for i, elem := range rows {
		if bytes.Equal(elem.Key, vKey) {

			newRows := append(rows[:i], rows[i+1:]...)
			newRows = append(newRows, RowData{Key: vKey, Data: vData})
			err = db.writeAll(newRows)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return StorageError{"There is no agent with that name"}
}

func (db *DataBasePl) Delete(vKey []byte) error {
	rows, err := db.readAll()
	if err != nil {
		return err
	}

	for i, elem := range rows {
		if bytes.Equal(elem.Key, vKey) {
			newRows := append(rows[:i], rows[i+1:]...)
			err = db.writeAll(newRows)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return StorageError{"There is no agent with that name"}
}

// Errors
type StorageError struct {
	message string
}

func (se StorageError) Error() string {
	return se.message
}
