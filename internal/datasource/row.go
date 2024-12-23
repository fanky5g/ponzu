package datasource 

type Row interface {
	Columns() []string
}
