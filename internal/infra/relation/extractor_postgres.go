package relation

import (

	// import postgresql connector
	_ "github.com/lib/pq"

	"github.com/xo/dburl"
	"makeit.imfr.cgi.com/lino/pkg/relation"
)

// NewPostgresExtractorFactory creates a new postgres extractor factory.
func NewPostgresExtractorFactory() *PostgresExtractorFactory {
	return &PostgresExtractorFactory{}
}

// PostgresExtractorFactory exposes methods to create new Postgres extractors.
type PostgresExtractorFactory struct{}

// New return a Postgres extractor
func (e *PostgresExtractorFactory) New(url string) relation.Extractor {
	return NewPostgresExtractor(url)
}

// PostgresExtractor provides relation extraction logic from Postgres database.
type PostgresExtractor struct {
	url string
}

// NewPostgresExtractor creates a new postgres extractor.
func NewPostgresExtractor(url string) *PostgresExtractor {
	return &PostgresExtractor{
		url: url,
	}
}

// Extract relations from the database.
func (e *PostgresExtractor) Extract() ([]relation.Relation, *relation.Error) {
	db, err := dburl.Open(e.url)
	if err != nil {
		return nil, &relation.Error{Description: err.Error()}
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, &relation.Error{Description: err.Error()}
	}

	SQL := `
SELECT
    tc.constraint_name,
    tc.table_schema,
    tc.table_name,
    kcu.column_name,
    ccu.table_schema AS foreign_table_schema,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM
    information_schema.table_constraints AS tc
    JOIN information_schema.key_column_usage AS kcu
      ON tc.constraint_name = kcu.constraint_name
      AND tc.table_schema = kcu.table_schema
    JOIN information_schema.constraint_column_usage AS ccu
      ON ccu.constraint_name = tc.constraint_name
      AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
            `

	rows, err := db.Query(SQL)
	if err != nil {
		return nil, &relation.Error{Description: err.Error()}
	}

	relations := []relation.Relation{}

	var (
		relationName string
		sourceSchema string
		sourceTable  string
		sourceColumn string
		targetSchema string
		targetTable  string
		targetColumn string
	)

	for rows.Next() {
		err := rows.Scan(&relationName, &sourceSchema, &sourceTable, &sourceColumn, &targetSchema, &targetTable, &targetColumn)
		if err != nil {
			return nil, &relation.Error{Description: err.Error()}
		}
		relation := relation.Relation{
			Name: relationName,
			Parent: relation.Table{
				Name: sourceSchema + "." + sourceTable,
				Keys: []string{sourceColumn},
			},
			Child: relation.Table{
				Name: targetSchema + "." + targetTable,
				Keys: []string{targetColumn},
			},
		}
		relations = append(relations, relation)
	}
	err = rows.Err()
	if err != nil {
		return nil, &relation.Error{Description: err.Error()}
	}

	return relations, nil
}