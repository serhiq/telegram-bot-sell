package entity

import (
	"gorm.io/gorm"
	"strconv"
)

type MenuItemDatabase struct {
	Name        string `json:"name"`
	StoreID     string `json:"storeId,omitempty" gorm:"column:store_id"`
	UserID      string `json:"userId,omitempty"`
	UUID        string `json:"uuid"`
	ParentUUID  string `json:"parentUuid"  gorm:"column:parent_uuid"`
	Group       bool   `json:"group"`
	Image       string `json:"image"`
	MeasureName string `json:"measureName,omitempty"  gorm:"column:measure"`
	Price       int
}

type Tabler interface {
	TableName() string
}

func (MenuItemDatabase) TableName() string {
	return "menu"

}
func (menuItem *MenuItemDatabase) PriceString() string {
	return strconv.FormatInt(int64(menuItem.Price), 10) + " руб"
}

func (g *GormDatabase) ImportMenu(items []*MenuItemDatabase) error {

	err := g.Db.Transaction(func(tx *gorm.DB) error {

		g.Db.Exec("DELETE FROM menu")

		for _, item := range items {
			if err := g.Db.Create(item).Error; err != nil {
				return err
			}

		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (g *GormDatabase) GetMenuItemByParent(parent string) ([]MenuItemDatabase, error) {
	var items []MenuItemDatabase
	err := g.Db.Where("parent_uuid = ?", parent).Model(&MenuItemDatabase{}).Find(&items).Error
	return items, err
}

func (g *GormDatabase) GetMenu(id string) (*MenuItemDatabase, error) {
	menuItem := new(MenuItemDatabase)
	err := g.Db.Where("uuid = ?", id).Find(menuItem).Error
	return menuItem, err
}
