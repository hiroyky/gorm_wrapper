package gormwrapper

type Model interface {
	TableName() string
}

type OrderDirection string

const (
	Ascending  OrderDirection = "ASC"
	Descending OrderDirection = "DESC"
)
