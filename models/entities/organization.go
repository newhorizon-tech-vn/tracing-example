package entities

type Organization struct {
	Base
	OrganizationId   int       `gorm:"column:OrganizationId;primaryKey"`
	OrganizationName string    `gorm:"column:OrganizationName;size:100;not null"`
	Schools          []*School `json:"schools,omitempty" gorm:"references:OrganizationId;foreignKey:OrganizationId;"`
}

func (o *Organization) GetSchools() []*School {
	if o == nil {
		return nil
	}
	return o.Schools
}

func (o *Organization) GetId() int {
	if o == nil {
		return 0
	}
	return o.OrganizationId
}

func (Organization) TableName() string {
	return "organizations"
}
