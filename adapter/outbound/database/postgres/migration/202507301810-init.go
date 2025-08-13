package migration

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/robertobff/nexpos/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var M_202507301810 *gormigrate.Migration = func() *gormigrate.Migration {
	type BaseTimestamps struct {
		CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime"`
		UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`
		DeletedAt *time.Time `gorm:"column:deleted_at;index"`
	}

	type BaseID struct {
		ID *string `gorm:"type:uuid;primaryKey"`
	}

	type Base struct {
		BaseID
		BaseTimestamps
	}

	type Image struct {
		Base
		Name        *string `gorm:"column:name"`
		Data        *[]byte `gorm:"column:data;type:bytea"`
		ContentType *string `gorm:"column:content_type"`
	}

	type User struct {
		Base
		Name        *string    `gorm:"column:name"`
		Email       *string    `gorm:"column:email"`
		Username    *string    `gorm:"column:username;unique"`
		BirthDate   *time.Time `gorm:"column:birth_date"`
		ImageID     *string    `gorm:"column:image_id;type:uuid"`
		Cpf         *string    `gorm:"column:cpf;unique"`
		PhoneNumber *string    `gorm:"column:phone_number"`
		ExternalID  *string    `gorm:"column:external_id;unique"`
	}

	type Item struct {
		Base
		Name        *string  `gorm:"column:name"`
		Description *string  `gorm:"column:description"`
		ImageID     *string  `gorm:"column:image_id;type:uuid"`
		Price       *float64 `gorm:"column:price"`
		CategoryID  *string  `gorm:"column:category_id;type:uuid"`
	}

	type PaidStatus int
	type UserOrders struct {
		Base
		UserID *string     `gorm:"column:user_id;type:uuid"`
		Item   []*Item     `gorm:"many2many:user_orders_items; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		Date   *time.Time  `gorm:"column:date"`
		Total  *float64    `gorm:"column:total"`
		Status *PaidStatus `gorm:"column:status"`
	}

	type Country struct {
		Base
		Name       *string `gorm:"column:name"`
		Identifier *string `gorm:"column:identifier"`
	}

	type State struct {
		Base
		Name       *string `gorm:"column:name"`
		Identifier *string `gorm:"column:identifier"`
		CountryID  *string `gorm:"column:country_id;type:uuid"`
	}

	type City struct {
		Base
		Name     *string `gorm:"column:name"`
		StateID  *string `gorm:"column:state_id;type:uuid"`
		StreetID *string `gorm:"column:street_id;type:uuid"`
	}

	type District struct {
		Base
		Name   *string `gorm:"column:name"`
		CityID *string `gorm:"column:city_id;type:uuid"`
	}

	type Street struct {
		Base
		Name       *string `gorm:"column:name"`
		ZipCode    *string `gorm:"column:zip_code"`
		Number     *string `gorm:"column:number"`
		DistrictID *string `gorm:"column:district_id;type:uuid"`
	}

	type UserAddress struct {
		Base
		UserID *string `gorm:"column:user_id;type:uuid"`
		CityID *string `gorm:"column:city_id;type:uuid"`
	}

	type UserOrdersItem struct {
		Base
		ItemID       *string `gorm:"column:item_id;type:uuid"`
		UserOrdersID *string `gorm:"column:user_id;type:uuid"`
	}

	type Category struct {
		Base
		Name        *string `gorm:"column:name"`
		Description *string `gorm:"column:description"`
		ImageID     *string `gorm:"column:image_id;type:uuid"`
	}

	type Discount struct {
		Base
		CategoryID *string    `gorm:"column:category_id;type:uuid"`
		ItemID     *string    `gorm:"column:item_id;type:uuid"`
		Date       *time.Time `gorm:"column:date"`
		Value      *float64   `gorm:"column:value"`
	}

	return &gormigrate.Migration{
		ID: "202507301810-init",
		Migrate: func(db *gorm.DB) error {
			return db.Transaction(
				func(tx *gorm.DB) error {
					if err := tx.AutoMigrate(
						&Image{},
						&User{},
						&Item{},
						&UserOrders{},
						&UserOrdersItem{},
						&Country{},
						&City{},
						&District{},
						&Street{},
						&UserAddress{},
						&Category{},
						&Discount{},
					); err != nil {
						return err
					}
					{
						if err := tx.Create(&Country{
							Base: Base{
								BaseID: BaseID{
									ID: utils.PString(uuid.NewV4().String()),
								},
								BaseTimestamps: BaseTimestamps{
									CreatedAt: utils.PTime(time.Now()),
								},
							},
							Name:       utils.PString("Brasil"),
							Identifier: utils.PString("BR"),
						}).Error; err != nil {
							return err
						}
					}
					return nil
				},
			)
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable(
				&Image{},
				&User{},
				&Item{},
				&UserOrders{},
				&UserOrdersItem{},
				&Country{},
				&City{},
				&District{},
				&Street{},
				&UserAddress{},
				&Category{},
				&Discount{},
			)
		},
	}
}()
