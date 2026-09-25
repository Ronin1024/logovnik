package logovnik

import (
	"encoding/json"
	"go_logovnik/internal/core/logovnik/db_record"

	file_repository "go_logovnik/internal/core/logovnik/repository"
	"slices"
	"strings"

	"github.com/google/uuid"
)

// Расшифрованный item
type Item struct {
	Id       uuid.UUID `json:"id"`
	Program  string    `json:"program"`
	Desc     string    `json:"desc"`
	Login    string    `json:"login"`
	Password string    `json:"password"`
}

type NewItem struct {
	Program  string `json:"program"`
	Desc     string `json:"desc"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Repository interface {
	GetItems() (map[uuid.UUID]db_record.Item, error)
	AddItem(item db_record.Item) (*db_record.Item, error)
	DeleteItem(id uuid.UUID) error
	UpdateItem(item db_record.Item) (*db_record.Item, error)
	Close() error
}

type Logovnik struct {
	master_key string
	filename   string
	repo       Repository
}

func NewLogovnik(file string, master_key string) (*Logovnik, error) {
	r, err := file_repository.NewRepository(file, master_key)
	if err != nil {
		return nil, err
	}
	l := &Logovnik{filename: file, repo: r}
	return l, nil
}

func (l *Logovnik) Close() {
	if l.repo != nil {
		l.repo.Close()
	}
}

func (l *Logovnik) GetItems() (string, error) {
	r_items, err := l.repo.GetItems()
	if err != nil {
		return "", err
	}
	items := []Item{}
	for _, v := range r_items {
		items = append(items, Item{
			Id:       v.Id,
			Program:  v.Program,
			Desc:     v.Desc,
			Login:    v.Login,
			Password: v.Password,
		})
	}

	// Сортировка по возрастанию UUID
	slices.SortFunc(items, func(a, b Item) int {
		return strings.Compare(a.Program, b.Program)
	})

	jsonData, err := json.Marshal(items)
	if err != nil {
		return "{}", err
	}
	return string(jsonData), nil
}

func (l *Logovnik) DeleteItem(id string) error {
	nid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = l.repo.DeleteItem(nid)
	if err != nil {
		return err
	}
	return nil
}

func (l *Logovnik) AddItem(itemJson string) (string, error) {
	newItem := NewItem{}
	err := json.Unmarshal([]byte(itemJson), &newItem)
	if err != nil {
		return "{}", err
	}
	dbitem := db_record.Item{
		Id:       uuid.New(),
		Program:  newItem.Program,
		Desc:     newItem.Desc,
		Login:    newItem.Login,
		Password: newItem.Password,
	}
	ritem, err := l.repo.AddItem(dbitem)
	if err != nil {
		return "", err
	}

	item := Item{
		Id:       ritem.Id,
		Program:  ritem.Program,
		Desc:     ritem.Desc,
		Login:    ritem.Login,
		Password: ritem.Password,
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		return "{}", err
	}
	return string(jsonData), nil
}

func (l *Logovnik) UpdateItem(itemJson string) (string, error) {
	item := Item{}
	err := json.Unmarshal([]byte(itemJson), &item)
	if err != nil {
		return "{}", err
	}
	dbitem := db_record.Item{
		Id:       item.Id,
		Program:  item.Program,
		Desc:     item.Desc,
		Login:    item.Login,
		Password: item.Password,
	}
	ritem, err := l.repo.UpdateItem(dbitem)
	if err != nil {
		return "", err
	}

	item = Item{
		Id:       ritem.Id,
		Program:  ritem.Program,
		Desc:     ritem.Desc,
		Login:    ritem.Login,
		Password: ritem.Password,
	}
	jsonData, err := json.Marshal(item)
	if err != nil {
		return "{}", err
	}
	return string(jsonData), nil
}
