package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type PingResponse struct {
	Message string `json:"message"`
}

type Product struct {
	Variants []Variant `json:"variants"`
}

type Variant struct {
	ID                   int64            `json:"id,omitempty"`
	ProductID            int64            `json:"product_id,omitempty"`
	Title                string           `json:"title,omitempty"`
	Sku                  string           `json:"sku,omitempty"`
	Position             int              `json:"position,omitempty"`
	Grams                int              `json:"grams,omitempty"`
	InventoryPolicy      string           `json:"inventory_policy,omitempty"`
	Price                *decimal.Decimal `json:"price,omitempty"`
	CompareAtPrice       *decimal.Decimal `json:"compare_at_price,omitempty"`
	FulfillmentService   string           `json:"fulfillment_service,omitempty"`
	InventoryManagement  string           `json:"inventory_management,omitempty"`
	InventoryItemId      int64            `json:"inventory_item_id,omitempty"`
	Option1              string           `json:"option1,omitempty"`
	Option2              string           `json:"option2,omitempty"`
	Option3              string           `json:"option3,omitempty"`
	CreatedAt            *time.Time       `json:"created_at,omitempty"`
	UpdatedAt            *time.Time       `json:"updated_at,omitempty"`
	Taxable              bool             `json:"taxable,omitempty"`
	TaxCode              string           `json:"tax_code,omitempty"`
	Barcode              string           `json:"barcode,omitempty"`
	ImageID              int64            `json:"image_id,omitempty"`
	InventoryQuantity    int              `json:"inventory_quantity,omitempty"`
	Weight               *decimal.Decimal `json:"weight,omitempty"`
	WeightUnit           string           `json:"weight_unit,omitempty"`
	OldInventoryQuantity int              `json:"old_inventory_quantity,omitempty"`
	RequireShipping      bool             `json:"requires_shipping"`
	AdminGraphqlAPIID    string           `json:"admin_graphql_api_id,omitempty"`
	Metafields           []Metafield      `json:"metafields,omitempty"`
}

type Metafield struct {
	CreatedAt         *time.Time  `json:"created_at,omitempty"`
	Description       string      `json:"description,omitempty"`    //Description of the metafield.
	ID                int64       `json:"id,omitempty"`             //Assigned by Shopify, used for updating a metafield.
	Key               string      `json:"key,omitempty"`            //The unique identifier for a metafield within its namespace, 3-64 characters long.
	Namespace         string      `json:"namespace,omitempty"`      //The container for a group of metafields, 3-255 characters long.
	OwnerId           int64       `json:"owner_id,omitempty"`       //The unique ID of the resource the metafield is for, i.e.: an Order ID.
	OwnerResource     string      `json:"owner_resource,omitempty"` //The type of reserouce the metafield is for, i.e.: and Order.
	UpdatedAt         *time.Time  `json:"updated_at,omitempty"`     //
	Value             interface{} `json:"value,omitempty"`          //The data stored in the metafield. Always stored as a string, use Type field for actual data type.
	Type              string      `json:"type,omitempty"`           //One of Shopify's defined types, see metafieldType.
	AdminGraphqlAPIID string      `json:"admin_graphql_api_id,omitempty"`
}

type User struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"phone"`
	Email        string `json:"email"`
}

type CartItem struct {
	VariantID string `json:"variant_id"`
	Quantity  int    `json:"quantity"`
}

type ViewCartItem struct {
	VariantID    string `json:"variant_id"`
	VariantTitle string `json:"variant_title"`
}

type GetVariantResponse struct {
	Variant Variant `json:"variant"`
}