// 自动生成模板{{.StructName}}
package model

import (
	"GF_PROJECT_NAME/global"
)

{{.GoStructString}}
{{ if .TableName }}
func ({{.StructName}}) TableName() string {
  return "{{.TableName}}"
}
{{ end }}
