package models

type DocumentInterface interface {
	Value() (interface{}, error)
	Scan(src interface{}) error
}
