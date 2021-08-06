package request

import "GF_PROJECT_NAME/model"

type {{.StructName}}Search struct{
    model.{{.StructName}}
    PageInfo
}