package table_formatting

import (
	"strings"

	"github.com/zeus-fyi/tables-to-go/pkg/settings"
)

func GenerateImports(content *strings.Builder, settings *settings.Settings, columnInfo ColumnInfo) {

	if !columnInfo.isNullableOrTemporal() && !settings.IsMastermindStructableRecorder {
		return
	}

	content.WriteString("import (\n")

	if columnInfo.isNullable && settings.IsNullTypeSQL() {
		content.WriteString("\t\"database/sql\"\n")
	}

	if columnInfo.isTemporal {
		content.WriteString("\t\"time\"\n")
	}

	if settings.IsMastermindStructableRecorder {
		content.WriteString("\t\n\"github.com/Masterminds/structable\"\n")
	}

	content.WriteString(")\n\n")
}
