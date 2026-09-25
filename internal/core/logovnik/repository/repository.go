package file_repository

import (
	"encoding/gob"
	"fmt"
	"go_logovnik/internal/core/logovnik/db_record"
	"go_logovnik/internal/core/logovnik/lcrypt"
	"go_logovnik/internal/utils"
	"io"
	"os"

	"github.com/google/uuid"
)

type Item struct {
	Id       uuid.UUID
	Program  []byte
	Desc     []byte
	Login    []byte
	Password []byte
	Salt     []byte
}

type Repository struct {
	f          *os.File
	filename   string
	master_key string
	Items      map[uuid.UUID]db_record.Item
}

func NewRepository(file string, master_key string) (*Repository, error) {
	if !utils.FileExists(file) {
		err := os.WriteFile(file, []byte{}, 0644)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	r := &Repository{
		filename:   file,
		f:          f,
		Items:      make(map[uuid.UUID]db_record.Item),
		master_key: master_key,
	}
	err = r.ReadData(file)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Repository) Close() error {
	return r.f.Close()
}

func (r *Repository) ReadData(file string) error {
	if !utils.FileExists(file) {
		return fmt.Errorf("file not exists")
	}
	if _, err := r.f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	decoder := gob.NewDecoder(r.f)
	var decoded []Item
	if err := decoder.Decode(&decoded); err != nil {
		// Пустой файл — это не ошибка, просто нет данных
		if err == io.EOF {
			r.Items = make(map[uuid.UUID]db_record.Item)
			return nil
		}
		return err
	}

	for i := 0; i < len(decoded); i++ {
		key, _, err := lcrypt.DeriveKeyPBKDF2(r.master_key, decoded[i].Salt)
		if err != nil {
			return err
		}
		decodedPassword, err := lcrypt.Decrypt(string(decoded[i].Password), key)
		if err != nil {
			return err
		}
		r.Items[decoded[i].Id] = db_record.Item{
			Id:       decoded[i].Id,
			Program:  string(decoded[i].Program),
			Desc:     string(decoded[i].Desc),
			Login:    string(decoded[i].Login),
			Password: decodedPassword,
		}
	}
	return nil
}

func (r *Repository) GetItems() (map[uuid.UUID]db_record.Item, error) {
	return r.Items, nil
}

// Write перезаписывает файл текущим содержимым r.items.
func (r *Repository) Write() error {
	// 1. Обрезаем файл до нуля (стираем старое содержимое)
	if err := r.f.Truncate(0); err != nil {
		return err
	}

	// 2. Возвращаем курсор в начало
	if _, err := r.f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	items := []Item{}
	key, salt, err := lcrypt.DeriveKeyPBKDF2(r.master_key, nil)
	if err != nil {
		return err
	}
	for _, v := range r.Items {
		encryptedPassword, err := lcrypt.Encrypt(v.Password, key)
		if err != nil {
			return err
		}

		items = append(items, Item{
			Id:       v.Id,
			Program:  []byte(v.Program),
			Desc:     []byte(v.Desc),
			Login:    []byte(v.Login),
			Password: []byte(encryptedPassword),
			Salt:     salt,
		})
	}

	// 3. Кодируем весь слайс
	encoder := gob.NewEncoder(r.f)
	if err := encoder.Encode(items); err != nil {
		return err
	}

	// 4. Сбрасываем буферы ОС на диск (по желанию)
	return r.f.Sync()
}

func (r *Repository) AddItem(item db_record.Item) (*db_record.Item, error) {
	r.Items[item.Id] = item
	if err := r.Write(); err != nil {
		return nil, err
	}
	return &item, nil
}
func (r *Repository) DeleteItem(id uuid.UUID) error {
	if _, ok := r.Items[id]; !ok {
		return fmt.Errorf("not found")
	} else {
		delete(r.Items, id)
	}

	if err := r.Write(); err != nil {
		return err
	}
	return nil
}
func (r *Repository) UpdateItem(item db_record.Item) (*db_record.Item, error) {
	if _, ok := r.Items[item.Id]; !ok {
		return nil, fmt.Errorf("not found")
	}
	r.Items[item.Id] = item
	if err := r.Write(); err != nil {
		return nil, err
	}
	return &item, nil
}
