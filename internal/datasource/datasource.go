package datasource

type DataSource interface {
	GetNumberOfRows() (int, error)
	GetColumns() ([]string, error)
	LoadData(offset int) ([]interface{}, error)
}
