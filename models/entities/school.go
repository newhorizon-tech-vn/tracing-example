package entities

type UserOrgSchool struct {
	Base
	UserId         int           `gorm:"column:UserId"`
	SchoolId       int           `gorm:"column:SchoolId"`
	OrganizationId int           `gorm:"column:OrganizationId"`
	School         *School       `gorm:"references:SchoolId;foreignKey:SchoolId;"`
	Organization   *Organization `gorm:"references:OrganizationId;foreignKey:OrganizationId;"`
}

func (u *UserOrgSchool) GetSchool() *School {
	if u == nil {
		return nil
	}
	return u.School
}

func (u *UserOrgSchool) GetOrganization() *Organization {
	if u == nil {
		return nil
	}
	return u.Organization
}

func (*UserOrgSchool) TableName() string {
	return "user_org_school"
}

type School struct {
	Base
	SchoolId           int      `json:"schoolId" example:"1" gorm:"column:SchoolId;primaryKey"`
	SchoolName         string   `json:"schoolName" example:"School 1" gorm:"column:SchoolName;size:100"`
	OrganizationId     int      `json:"organizationId" example:"1" gorm:"column:OrganizationId;not null"`
	SchoolContactName  string   `json:"schoolContactName" example:"Ms.Abc" gorm:"column:SchoolContactName;size:30"`
	SchoolContactPhone string   `json:"schoolContactPhone" example:"+84912345678" gorm:"column:SchoolContactPhone;size:20"`
	SchoolContactEmail string   `json:"schoolContactEmail" example:"abc@gmail.com" gorm:"column:SchoolContactEmail;size:50"`
	SchoolAddress      string   `json:"schoolAddress" example:"123 Example street" gorm:"column:SchoolAddressLine1;size:250"`
	CityId             int      `json:"cityId" example:"1" gorm:"column:CityId"`
	DistrictId         int      `json:"districtId" example:"1" gorm:"column:DistrictId"`
	CountryId          int      `json:"countryId" example:"1" gorm:"column:CountryId"`
	Classes            []*Class `json:"classes,omitempty" gorm:"references:SchoolId;foreignKey:school_SchoolId;"`
}

func (s *School) GetClasses() []*Class {
	if s == nil {
		return nil
	}
	return s.Classes
}

func (School) TableName() string {
	return "schools"
}
