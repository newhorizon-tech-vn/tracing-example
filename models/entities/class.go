package entities

type ClassStudent struct {
	BaseV2
	ID      int    `gorm:"column:id;primary key"`
	Uuid    string `gorm:"column:uuid;type:char;size:36;not null"`
	ClassId int    `gorm:"column:school_class_id;not null"`
	UserId  int    `gorm:"column:user_UserId;not null"`
	Class   *Class `json:"classes,omitempty" gorm:"references:school_class_id;foreignKey:id;"`
}

func (*ClassStudent) TableName() string {
	return "class_student"
}

type ClassTeacher struct {
	BaseV2
	ID      int    `gorm:"column:id;primary key"`
	Uuid    string `gorm:"column:uuid;type:char;size:36;not null"`
	ClassId int    `gorm:"column:school_class_id;not null"`
	UserId  int    `gorm:"column:user_UserId;not null"`
	Class   *Class `json:"classes,omitempty" gorm:"references:school_class_id;foreignKey:id;"`
}

func (*ClassTeacher) TableName() string {
	return "class_teacher"
}

type Class struct {
	BaseV2    `swaggerignore:"true"`
	Uuid      string `gorm:"column:uuid;type:char;size:36;not null"`
	ClassId   int    `gorm:"column:id;primary key"`
	ClassName string `gorm:"column:name;size:100;not null"`
	SchoolId  int    `gorm:"column:school_SchoolId;not null"`
	GradeId   int    `gorm:"column:grade_id"`
}
