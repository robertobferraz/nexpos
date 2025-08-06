package dto

type GormQuery struct {
	Where      *[]GormWhere
	Preload    *[]GormPreload
	Debug      bool
	InnerJoins *[]GormInnerJoins
	Group      *GormGroup
	Model      *GormModel
	Select     *GormSelect
	Join       *[]GormJoin
	Order      *[]GormOrder
}

type GormOrder struct {
	Field interface{}
}

type GormWhere struct {
	Column    string
	Condition string
	Value     any
}

type GormPreload struct {
	Field string
}

type GormInnerJoins struct {
	Field string
	Where *[]GormWhere
}
type GormGroup struct {
	Field string
}

type GormModel struct {
	Interface interface{}
}
type GormSelect struct {
	Arg string
}

type GormJoin struct {
	Query string
}
