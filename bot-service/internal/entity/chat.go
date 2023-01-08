package entity

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Gorm сам не может однозначно выбрать размер поля, потому всегда ставит по максимуму, что не эффективно
// стоит всегда указывать хотя бы размер поля, а иногда - точный тип с размером
type Chat struct {
	ChatId                int64  `gorm:"column:id;primary_key"`
	NameUser              string `gorm:"column:user_name"`	// зачем делать названия колонок в БД отличными от названия полей? Это лишь создаст путаницу
	PhoneUser             string `gorm:"column:user_phone"`
	ChatState             string `gorm:"column:chat_state"`	// не очевидны допустимые значения, поскольку в Go нет перечислений, то попросту создается новый тип для `string`
	Store                 string `gorm:"column:store"`
	CurrentFolderMenu     string `gorm:"column:last_product"`
	OrderStr              string `gorm:"column:order"`
	LastEditedMenuItemStr string `gorm:"column:last_edited_menu_item"`
}

func NewChat(chatId int64) *Chat {
	// старайтесь придерживаться единого стиля создания экземпляров структур: либо `new(Struct)`, либо `&Struct{}`. Последний вариант мне нравится больще
	p := new(Chat)
	p.ChatId = chatId
	p.ChatState = INPUT_NAME
	return p
}

func (c *Chat) HaveContact() bool {
	return len(c.NameUser) != 0 && len(c.PhoneUser) != 0
}
func (c *Chat) HaveUserName() bool {
	return len(c.NameUser) != 0
}
func (c *Chat) HaveUserPhone() bool {
	return len(c.PhoneUser) != 0
}

// можно сделать красивее через `Scanner`/`Valuer`: https://gorm.io/docs/data_types.html
// обратите внимание, что это механизм стандартной библиотеки, а не Gorm.
func (c *Chat) GetOrder() *OrderRequest {
	var o = &OrderRequest{}
	err := json.Unmarshal([]byte(c.OrderStr), o)
	if err != nil {
		log.Print("order: unmarshal error")
	}
	return o
}

func (c *Chat) NewOrder() {
	order := &OrderRequest{}
	c.OrderStr = order.ToJson()
}

const (
	// выбор магазина
	//SELECT_STORE     = "SELECT_STORE"
	INPUT_NAME          = "INPUT_NAME"
	INPUT_PHONE         = "INPUT_PHONE"
	STATE_PREPARE_ORDER = "PREPARE_ORDER"
)

///////////////////////////////////////////////////////
// это одна из попыток избежать циклических зависимостей? Есть же прекрасная структура `Database`
// но исходя из методов - это скорее `ChatsRepository`, чем `GormDatabase`
type GormDatabase struct {
	Db *gorm.DB
}

func CreateGorm(db *gorm.DB) *GormDatabase {
	return &GormDatabase{Db: db}
}

func (g *GormDatabase) InsertChat(chat *Chat) error {
	return g.Db.Create(chat).Error
}

func (r *GormDatabase) UpdateChat(chat *Chat) error {
	return r.Db.Save(chat).Error
}

func (r *GormDatabase) GetChat(id int64) (*Chat, error) {
	tag := new(Chat)
	err := r.Db.Where("id = ?", id).Find(tag).Error
	return tag, err
}

func (r *GormDatabase) GetOrCreateChat(id int64) (*Chat, error) {
	chat := new(Chat)
	err := r.Db.Where("id = ?", id).Take(chat).Error
	// спорно, ошибка может быть не только `gorm.ErrRecordNotFound`, но и какая-нибудь другая
	// надо явно проверять что за ошибка вернулась
	if err != nil {
		newChat := NewChat(id)
		r.InsertChat(newChat)
		return newChat, err

	}
	return chat, err
}

func (r *GormDatabase) DeleteChat(id string) error {
	result := r.Db.Select(clause.Associations).Unscoped().Delete(&Chat{}, id)
	return result.Error
}

func (c *Chat) SaveLaseEdited(item LastEditedMenuItem) {
	c.LastEditedMenuItemStr = item.ToJson()
}

func (c *Chat) GetLastEditedMenuItem() *LastEditedMenuItem {
	return FromJsonLastEditedMenuItem(c.LastEditedMenuItemStr)
}
func (c *Chat) IsCorrectName(name string) bool {
	return name != ""
}

func (c *Chat) IsCorrectPhone(phone string) bool {
	// очевидно, проверка сильно упрощенная)
	// и непонятно как она связана со структурой `Chat`, по сути это функция просится в отдельный модуль
	return len(phone) > 4
}

type LastEditedMenuItem struct {
	UuidMenuItem string
	MessageId    int
}

func (c *LastEditedMenuItem) ToJson() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		// ай-яй-яй так делать, никаких паник из глубин кода - из-за того что в одном чате что-то пошло не так Вы уронили всего бота
		panic(err)
	}
	return string(bytes)
}

func FromJsonLastEditedMenuItem(str string) *LastEditedMenuItem {
	var stu = &LastEditedMenuItem{}
	// а ошибку проверить?)
	json.Unmarshal([]byte(str), stu)
	return stu
}
