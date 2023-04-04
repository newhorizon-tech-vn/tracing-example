package entities

import "encoding/json"

type User struct {
	UserId         int    `json:"userId,omitempty" gorm:"column:UserId;primaryKey"`
	RoleId         int    `json:"roleId" gorm:"column:RoleId"`
	UserFirstName  string `json:"userFirstName" gorm:"column:UserFirstName;size:100"`
	UserMiddleName string `json:"userMiddleName" gorm:"column:UserMiddleName;size:100"`
	UserLastName   string `json:"userLastName" gorm:"column:UserLastName;size:100"`
	UserStatus     string `json:"userStatus" gorm:"column:UserStatus;size:100"`
}

func (p *User) GetRoleId() int {
	if p == nil {
		return 0
	}
	return p.RoleId
}

func (*User) TableName() string {
	return "users"
}

func (u *User) ToJSON() string {
	result, _ := json.Marshal(u)
	return string(result)
}

func (u *User) FromJSON(s string) error {
	if u == nil {
		u = &User{}
	}
	return json.Unmarshal([]byte(s), u)
}
