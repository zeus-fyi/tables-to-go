package table_formatting

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/zeus-fyi/tables-to-go/pkg/settings"
)

type ColumnInfo struct {
	isNullable bool
	isTemporal bool
}

func (c ColumnInfo) isNullableOrTemporal() bool {
	return c.isNullable || c.isTemporal
}

// FormatColumnName checks for invalid characters and transforms a column name
// according to the provided settings.
func FormatColumnName(settings *settings.Settings, column, table string) (string, error) {

	// Replace any whitespace with underscores
	columnName := strings.Map(ReplaceSpace, column)
	columnName = Caser.String(columnName)

	if settings.IsOutputFormatCamelCase() {
		columnName = CamelCaseString(columnName)
	}
	if settings.ShouldInitialism() {
		columnName = ToInitialisms(columnName)
	}

	// Check that the column name doesn't contain any invalid characters for Go variables
	if !ValidVariableName(columnName) {
		return "", fmt.Errorf("column name %q in table %q contains invalid characters", column, table)
	}

	// First character of an identifier in Go must be letter or _
	// We want it to be an uppercase letter to be a public field
	if !unicode.IsLetter(rune(columnName[0])) {
		prefix := "X_"
		if settings.IsOutputFormatCamelCase() {
			prefix = "X"
		}
		if settings.ShouldInitialism() {
			// Note we use the original passed in name of the column here to
			// avoid the Title'izing of the first non-digit character as done
			// by cases.Caser. Eg: `1fish2fish` gets transformed to `X1Fish2fish`
			// but we want `X1fish2fish`.
			columnName = ToInitialisms(column)
		}
		if settings.Verbose {
			fmt.Printf("\t\t>column %q in table %q doesn't start with a letter; prepending with %q\n", column, table, prefix)
		}
		columnName = prefix + columnName
	}

	return columnName, nil
}
