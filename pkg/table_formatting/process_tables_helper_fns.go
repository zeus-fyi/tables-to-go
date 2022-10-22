package table_formatting

import (
	"fmt"
	"strings"

	"github.com/zeus-fyi/tables-to-go/pkg/database"
	"github.com/zeus-fyi/tables-to-go/pkg/settings"
	"github.com/zeus-fyi/tables-to-go/pkg/tagger"
)

func CreateTableStructString(settings *settings.Settings, db database.Database, table *database.Table) (string, string, error) {
	var structFields strings.Builder
	tableName := Caser.String(settings.Prefix + table.Name + settings.Suffix)
	// Replace any whitespace with underscores
	tableName = strings.Map(ReplaceSpace, tableName)
	if settings.IsOutputFormatCamelCase() {
		tableName = CamelCaseString(tableName)
	}

	tagger := tagger.NewTaggers(settings)

	// Check that the table name doesn't contain any invalid characters for Go variables
	if !ValidVariableName(tableName) {
		return "", "", fmt.Errorf("table name %q contains invalid characters", table.Name)
	}

	columnInfo := ColumnInfo{}
	columns := map[string]struct{}{}

	for _, column := range table.Columns {
		columnName, err := FormatColumnName(settings, column.Name, table.Name)
		if err != nil {
			return "", "", err
		}

		// ISSUE-4: if columns are part of multiple constraints
		// then the sql returns multiple rows per column name.
		// Therefore, we check if we already added a column with
		// that name to the struct, if so, skip.
		if _, ok := columns[columnName]; ok {
			continue
		}
		columns[columnName] = struct{}{}

		if settings.VVerbose {
			fmt.Printf("\t\t> %v\r\n", column.Name)
		}

		columnType, col := MapDbColumnTypeToGoType(settings, db, column)

		// save that we saw types of columns at least once
		if !columnInfo.isTemporal {
			columnInfo.isTemporal = col.isTemporal
		}
		if !columnInfo.isNullable {
			columnInfo.isNullable = col.isNullable
		}

		structFields.WriteString(columnName)
		structFields.WriteString(" ")
		structFields.WriteString(columnType)
		structFields.WriteString(" ")
		structFields.WriteString(tagger.GenerateTag(db, column))
		structFields.WriteString("\n")
	}

	if settings.IsMastermindStructableRecorder {
		structFields.WriteString("\t\nstructable.Recorder\n")
	}

	var fileContent strings.Builder

	// write header infos
	fileContent.WriteString("package ")
	fileContent.WriteString(settings.PackageName)
	fileContent.WriteString("\n\n")

	// write imports
	GenerateImports(&fileContent, settings, columnInfo)

	// write struct with fields
	fileContent.WriteString("type ")
	fileContent.WriteString(tableName)
	fileContent.WriteString(" struct {\n")
	fileContent.WriteString(structFields.String())
	fileContent.WriteString("}")

	return tableName, fileContent.String(), nil
}
