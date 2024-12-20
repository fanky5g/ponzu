package datasource

import "io"

type DataSourceReaderFactory func(DataSource) (io.Reader, error)
