package cafmodel

import (
	"omega/internal/core/coract"
	"omega/internal/core/corerr"
	"omega/internal/core/corterm"
	"omega/internal/types"
	"omega/pkg/dict"
	"omega/pkg/limberr"
)

const (
	OrderFoodTable = "caf_order_foods"
)

type OrderFood struct {
	ID          types.RowID `gorm:"primary_key" json:"id,omitempty" `
	OrderID     types.RowID `json:"order_id,omitempty"`
	FoodID      types.RowID `json:"food_id,omitempty"`
	Price       int         `json:"price,omitempty"`
	Qty         int         `json:"qty,omitempty"`
	Total       int         `json:"total,omitempty"`
	Description string      `json:"description"`
	Food        string      `sql:"-" json:"food" table:"caf_foods.name as food"`
}

// Validate check the type of fields
func (p *OrderFood) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		// if len(p.Name) > 255 {
		// 	err = limberr.AddInvalidParam(err, "name",
		// 		corerr.MaximumAcceptedCharacterForVisV,
		// 		dict.R(corterm.Name), 255)
		// }

		// if p.Price == 0 {
		// 	err = limberr.AddInvalidParam(err, "price",
		// 		corerr.VisRequired, dict.R("price"))
		// }

		if len(p.Description) > 255 {
			err = limberr.AddInvalidParam(err, "description",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Description), 255)
		}
	}

	return err
}
