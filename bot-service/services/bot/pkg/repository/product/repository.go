package product

import (
	domainProduct "bot/pkg/type/product"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func New(Db *gorm.DB) *Repository {
	return &Repository{
		Db: Db,
	}
}

type Product struct {
	Name        string `json:"name"`
	StoreID     string `json:"storeId,omitempty" gorm:"column:store_id"`
	UserID      string `json:"userId,omitempty"`
	UUID        string `json:"uuid"`
	ParentUUID  string `json:"parentUuid"  gorm:"column:parent_uuid"`
	Group       bool   `json:"group"`
	Image       string `json:"image"`
	MeasureName string `json:"measureName,omitempty"  gorm:"column:measure"`
	Price       uint64 // Цена в копейках
}

type Tabler interface {
	TableName() string
}

func (Product) TableName() string {
	return "products"

}

func (g *Repository) ImportMenu(items []*Product) error {
	err := g.Db.Transaction(func(tx *gorm.DB) error {

		g.Db.Exec("DELETE FROM products")

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

func (g *Repository) GetMenuItemByParent(parent string) ([]*domainProduct.Product, error) {
	var items []Product
	var products = []*domainProduct.Product{}

	err := g.Db.Where("parent_uuid = ?", parent).Model(&Product{}).Find(&items).Error
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		products = append(products, ToDomain(&item))

	}

	return products, err
}

func (g *Repository) GetMenu(id string) (*domainProduct.Product, error) {
	var menuItem = domainProduct.Product{}
	err := g.Db.Where("uuid = ?", id).Find(&menuItem).Error
	return &menuItem, err
}

func ToDomain(r *Product) *domainProduct.Product {
	return &domainProduct.Product{
		Name: r.Name,
		//StoreID:     ,
		UUID:        r.UUID,
		ParentUUID:  r.ParentUUID,
		Group:       r.Group,
		Image:       r.Image,
		MeasureName: r.MeasureName,
		Price:       r.Price,
	}
}
