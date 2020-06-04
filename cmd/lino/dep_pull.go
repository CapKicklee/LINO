package main

import (
	"io"
	"os"

	infra "makeit.imfr.cgi.com/lino/internal/infra/pull"
	domain "makeit.imfr.cgi.com/lino/pkg/pull"
)

func pullDataSourceFactory() map[string]domain.DataSourceFactory {
	return map[string]domain.DataSourceFactory{
		"postgres": infra.NewPostgresDataSourceFactory(logger),
	}
}

func pullRowExporterFactory() func(file io.Writer) domain.RowExporter {
	return func(file io.Writer) domain.RowExporter {
		return infra.NewJSONRowExporter(file)
	}
}

func traceListner(file *os.File) domain.TraceListener {
	return infra.NewJSONTraceListener(file)
}
