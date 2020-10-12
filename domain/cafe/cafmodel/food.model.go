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
	FoodTable = "caf_foods"
)

type Food struct {
	types.GormCol
	Name        string `gorm:"not null;unique" json:"name,omitempty"`
	Price       int    `json:"price,omitempty"`
	Color       string `json:"color,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description"`
}

// Validate check the type of fields
func (p *Food) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if len(p.Name) > 255 {
			err = limberr.AddInvalidParam(err, "name",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Name), 255)
		}

		if p.Price == 0 {
			err = limberr.AddInvalidParam(err, "price",
				corerr.VisRequired, dict.R("price"))
		}

		if len(p.Description) > 255 {
			err = limberr.AddInvalidParam(err, "description",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Description), 255)
		}
	}

	return err
}
