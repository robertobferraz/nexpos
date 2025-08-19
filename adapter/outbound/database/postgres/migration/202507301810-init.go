package migration

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
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
		Email       *string    `gorm:"column:email;unique"`
		Username    *string    `gorm:"column:username;unique"`
		BirthDate   *time.Time `gorm:"column:birth_date"`
		ImageID     *string    `gorm:"column:image_id;type:uuid"`
		PhoneNumber *string    `gorm:"column:phone_number"`
		ExternalID  *string    `gorm:"column:external_id;unique"`
		Role        *string    `gorm:"column:role"` // Added for user roles like 'customer', 'admin'
	}
	type Category struct {
		Base
		Name        *string `gorm:"column:name"`
		Description *string `gorm:"column:description"`
		ImageID     *string `gorm:"column:image_id;type:uuid"`
		ParentID    *string `gorm:"column:parent_id;type:uuid"` // Added for hierarchical categories
	}
	type Item struct {
		Base
		Name        *string  `gorm:"column:name"`
		Description *string  `gorm:"column:description"`
		ImageID     *string  `gorm:"column:image_id;type:uuid"`
		Price       *float64 `gorm:"column:price"`
		CategoryID  *string  `gorm:"column:category_id;type:uuid"`
		Stock       *int     `gorm:"column:stock"`      // Added for inventory management
		SKU         *string  `gorm:"column:sku;unique"` // Added for unique product identifier
	}
	type DiscountType int // 0: fixed amount, 1: percentage
	type Discount struct {
		Base
		Code       *string       `gorm:"column:code;unique"` // Added for coupon codes
		Type       *DiscountType `gorm:"column:type"`
		Value      *float64      `gorm:"column:value"`
		StartDate  *time.Time    `gorm:"column:start_date"`
		EndDate    *time.Time    `gorm:"column:end_date"`
		CategoryID *string       `gorm:"column:category_id;type:uuid"` // Optional: applies to category
		ItemID     *string       `gorm:"column:item_id;type:uuid"`     // Optional: applies to item
		MinAmount  *float64      `gorm:"column:min_amount"`            // Optional: minimum order amount to apply
		MaxUses    *int          `gorm:"column:max_uses"`              // Optional: maximum uses per discount
	}
	type Country struct {
		Base
		Name         *string `gorm:"column:name"`
		Iso2         *string `gorm:"column:iso2;index"`
		Iso3         *string `gorm:"column:iso3;index"`
		PhoneCode    *string `gorm:"column:phone_code"`
		Capital      *string `gorm:"column:capital"`
		CurrencyCode *string `gorm:"column:currency_code"`
		Emoji        *string `gorm:"column:emoji"`
		ExternalID   *int    `gorm:"column:external_id;unique"`
	}
	type State struct {
		Base
		Name       *string `gorm:"column:name"`
		Iso2       *string `gorm:"column:iso2;index"`
		CountryID  *string `gorm:"column:country_id;type:uuid"`
		ExternalID *int    `gorm:"column:external_id;unique"`
	}
	type City struct {
		Base
		Name       *string `gorm:"column:name"`
		StateID    *string `gorm:"column:state_id;type:uuid"`
		ExternalID *int    `gorm:"column:external_id;unique"`
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
		DistrictID *string `gorm:"column:district_id;type:uuid"`
	}
	type Address struct { // Renamed from UserAddress for clarity, as it can be billing/shipping
		Base
		UserID     *string `gorm:"column:user_id;type:uuid"`
		StreetID   *string `gorm:"column:street_id;type:uuid"`
		Number     *string `gorm:"column:number"`      // Moved from Street, as it's per address
		Complement *string `gorm:"column:complement"`  // Added for address details like apt/suite
		IsBilling  *bool   `gorm:"column:is_billing"`  // Added to distinguish billing/shipping
		IsShipping *bool   `gorm:"column:is_shipping"` // Added to distinguish billing/shipping
		IsDefault  *bool   `gorm:"column:is_default"`  // Added for default address
	}
	type OrderStatus int
	type Order struct {
		Base
		UserID          *string      `gorm:"column:user_id;type:uuid"`
		BillingAddress  *string      `gorm:"column:billing_address_id;type:uuid"`  // Added reference to Address
		ShippingAddress *string      `gorm:"column:shipping_address_id;type:uuid"` // Added reference to Address
		Date            *time.Time   `gorm:"column:date"`
		Total           *float64     `gorm:"column:total"`
		Subtotal        *float64     `gorm:"column:subtotal"`        // Added to separate from discounts/taxes
		DiscountAmount  *float64     `gorm:"column:discount_amount"` // Added for applied discounts
		ShippingCost    *float64     `gorm:"column:shipping_cost"`   // Added for shipping fees
		TaxAmount       *float64     `gorm:"column:tax_amount"`      // Added for taxes
		Status          *OrderStatus `gorm:"column:status"`
		PaymentMethod   *string      `gorm:"column:payment_method"`  // Added like 'credit_card', 'paypal'
		TrackingNumber  *string      `gorm:"column:tracking_number"` // Added for shipment tracking
	}
	type OrderItem struct { // Renamed from UserOrdersItem, added fields for quantity and snapshot pricing
		Base
		OrderID   *string  `gorm:"column:order_id;type:uuid"`
		ItemID    *string  `gorm:"column:item_id;type:uuid"`
		Quantity  *int     `gorm:"column:quantity"`
		UnitPrice *float64 `gorm:"column:unit_price"` // Snapshot of item price at order time
		Discount  *float64 `gorm:"column:discount"`   // Per-item discount if applied
	}
	type Payment struct { // Added for payment details, linked to order
		Base
		OrderID       *string    `gorm:"column:order_id;type:uuid"`
		Amount        *float64   `gorm:"column:amount"`
		Date          *time.Time `gorm:"column:date"`
		Status        *string    `gorm:"column:status"` // e.g., 'success', 'failed'
		TransactionID *string    `gorm:"column:transaction_id"`
	}
	type Cart struct { // Added for shopping cart functionality
		Base
		UserID *string `gorm:"column:user_id;type:uuid"`
	}
	type CartItem struct {
		Base
		CartID   *string `gorm:"column:cart_id;type:uuid"`
		ItemID   *string `gorm:"column:item_id;type:uuid"`
		Quantity *int    `gorm:"column:quantity"`
	}
	type Review struct { // Added for product reviews
		Base
		UserID  *string    `gorm:"column:user_id;type:uuid"`
		ItemID  *string    `gorm:"column:item_id;type:uuid"`
		Rating  *int       `gorm:"column:rating"` // 1-5 stars
		Comment *string    `gorm:"column:comment"`
		Date    *time.Time `gorm:"column:date"`
	}

	return &gormigrate.Migration{
		ID: "202507301810-init",
		Migrate: func(db *gorm.DB) error {
			return db.Transaction(
				func(tx *gorm.DB) error {
					if err := tx.AutoMigrate(
						&Image{},
						&User{},
						&Category{},
						&Item{},
						&Discount{},
						&Country{},
						&State{},
						&City{},
						&District{},
						&Street{},
						&Address{},
						&Order{},
						&OrderItem{},
						&Payment{},
						&Cart{},
						&CartItem{},
						&Review{},
					); err != nil {
						return err
					}

					return nil
				},
			)
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable(
				&Image{},
				&User{},
				&Category{},
				&Item{},
				&Discount{},
				&Country{},
				&State{},
				&City{},
				&District{},
				&Street{},
				&Address{},
				&Order{},
				&OrderItem{},
				&Payment{},
				&Cart{},
				&CartItem{},
				&Review{},
			)
		},
	}
}()
