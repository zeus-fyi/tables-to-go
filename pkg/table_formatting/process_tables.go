package table_formatting

import (
	"fmt"

	"github.com/zeus-fyi/tables-to-go/pkg/database"
	"github.com/zeus-fyi/tables-to-go/pkg/settings"
)

type TableContentMap struct {
	TableContent map[string]string
	TableMap     map[string]*database.Table
}

func NewTableContentMap() TableContentMap {
	return TableContentMap{
		TableContent: make(map[string]string),
		TableMap:     make(map[string]*database.Table),
	}
}

func (t *TableContentMap) ProcessTables(db database.Database, settings *settings.Settings, tables ...*database.Table) error {
	tableContent := NewTableContentMap()

	if settings.Verbose {
		fmt.Printf("> number of tables: %v\r\n", len(tables))
	}

	if err := db.PrepareGetColumnsOfTableStmt(); err != nil {
		return fmt.Errorf("could not prepare the get-column-statement: %w", err)
	}

	for _, table := range tables {
		tableContent.TableMap[table.Name] = table
		if settings.Verbose {
			fmt.Printf("> processing table %q\r\n", table.Name)
		}

		if err := db.GetColumnsOfTable(table); err != nil {
			if !settings.Force {
				return err
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
				return fmt.Errorf("could not create string for table %q: %w", table.Name, tblErr)
			}
			fmt.Printf("could not create string for table %q: %v\n", table.Name, tblErr)
			continue
		}
		tableContent.TableContent[tableName] = content
	}
	t.TableContent = tableContent.TableContent
	t.TableMap = tableContent.TableMap
	return nil
}
