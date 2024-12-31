package mocks

import "github.com/stretchr/testify/mock"

type SearchClient struct {
	*mock.Mock
}

func (client *SearchClient) Update(id string, data interface{}) error {
	args := client.Mock.Called(id, data)
	return args.Error(0)
}

func (client *SearchClient) Delete(entityName, entityId string) error {
	args := client.Mock.Called(entityName, entityId)
	return args.Error(0)
}

func (client *SearchClient) Search(entityDefinition interface{}, query string, count, offset int) ([]interface{}, error) {
	args := client.Mock.Called(entityDefinition, query, count, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]interface{}), args.Error(1)
}

func (client *SearchClient) SearchWithPagination(entityDefinition interface{}, query string, count, offset int) ([]interface{}, int, error) {
	args := client.Mock.Called(entityDefinition, query, count, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}

	return args.Get(0).([]interface{}), args.Int(1), args.Error(2)
}
