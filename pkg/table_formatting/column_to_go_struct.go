package table_formatting

import (
	"github.com/fraenky8/tables-to-go/pkg/database"
	"github.com/fraenky8/tables-to-go/pkg/settings"
)

func MapDbColumnTypeToGoType(s *settings.Settings, db database.Database, column database.Column) (goType string, columnInfo ColumnInfo) {
	if db.IsInteger(column) {
		goType = "int"
		if db.IsNullable(column) {
			goType = GetNullType(s, "*int", "sql.NullInt64")
			columnInfo.isNullable = true
		}
	} else if db.IsFloat(column) {
		goType = "float64"
		if db.IsNullable(column) {
			goType = GetNullType(s, "*float64", "sql.NullFloat64")
			columnInfo.isNullable = true
		}
	} else if db.IsTemporal(column) {
		if !db.IsNullable(column) {
			goType = "time.Time"
			columnInfo.isTemporal = true
		} else {
			goType = GetNullType(s, "*time.Time", "sql.NullTime")
			columnInfo.isTemporal = s.Null == settings.NullTypeNative
			columnInfo.isNullable = true
		}
	} else {
		// TODO handle special data types
		switch column.DataType {
		case "boolean":
			goType = "bool"
			if db.IsNullable(column) {
				goType = GetNullType(s, "*bool", "sql.NullBool")
				columnInfo.isNullable = true
			}
		default:
			// Everything else we cannot detect defaults to (nullable) string.
			goType = "string"
			if db.IsNullable(column) {
				goType = GetNullType(s, "*string", "sql.NullString")
				columnInfo.isNullable = true
			}
		}
	}

	return goType, columnInfo
}
