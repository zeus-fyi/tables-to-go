package table_formatting

import (
	"fmt"

	"github.com/zeus-fyi/tables-to-go/pkg/database"
	"github.com/zeus-fyi/tables-to-go/pkg/settings"
)

func ProcessTables(db database.Database, settings *settings.Settings, tables ...*database.Table) (map[string]string, error) {
	tableContent := make(map[string]string)
	for _, table := range tables {
		if settings.Verbose {
			fmt.Printf("> processing table %q\r\n", table.Name)
		}

		if err := db.GetColumnsOfTable(table); err != nil {
			if !settings.Force {
				return tableContent, fmt.Errorf("could not get columns of table %q: %w", table.Name, err)
			}
			fmt.Printf("could not get columns of table %q: %v\n", table.Name, err)
			continue
		}

		if settings.Verbose {
			fmt.Printf("\t> number of columns: %v\r\n", len(table.Columns))
		}

		tableName, content, tblErr := CreateTableStructString(settings, db, table)
		if tblErr != nil {
			if !settings.Force {
				return tableContent, fmt.Errorf("could not create string for table %q: %w", table.Name, tblErr)
			}
			fmt.Printf("could not create string for table %q: %v\n", table.Name, tblErr)
			continue
		}
		tableContent[tableName] = content
	}
	return tableContent, nil
}
