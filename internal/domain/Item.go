package domain

// Item - объект мерча
type Item struct {
	Id   uint64 `json:"id"`   // Уникальный идентификатор мерча
	Name string `json:"name"` // Название предмета
	Cost int    `json:"cost"` // Стоимость
}

// CreateItem создает новый объект Item
func CreateItem(name string, cost int) *Item {
	return &Item{
		Name: name,
		Cost: cost,
	}
}
