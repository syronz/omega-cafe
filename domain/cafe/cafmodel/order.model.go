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
	OrderTable = "caf_orders"
)

type Order struct {
	types.GormCol
	Customer    string      `json:"customer,omitempty"`
	Phone       string      `json:"phone,omitempty"`
	Total       int         `json:"total,omitempty"`
	Discount    int         `json:"discount,omitempty"`
	Description string      `json:"description"`
	Foods       []OrderFood `sql:"-" json:"foods"`
}

// Validate check the type of fields
func (p *Order) Validate(act coract.Action) (err error) {

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