package main

import (
	"flag"
	"fmt"
	"GF_PROJECT_NAME/core"
	"GF_PROJECT_NAME/global"
	"GF_PROJECT_NAME/utils"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/xwb1989/sqlparser"
)

var (
	SqlGoTypeMap = map[string]string{
		"int":                "int",
		"integer":            "int",
		"tinyint":            "int8",
		"smallint":           "int16",
		"mediumint":          "int32",
		"bigint":             "int64",
		"int unsigned":       "uint",
		"integer unsigned":   "uint",
		"tinyint unsigned":   "uint8",
		"smallint unsigned":  "uint16",
		"mediumint unsigned": "uint32",
		"bigint unsigned":    "uint64",
		"bit":                "byte",
		"bool":               "bool",
		"enum":               "string",
		"set":                "string",
		"varchar":            "string",
		"char":               "string",
		"tinytext":           "string",
		"mediumtext":         "string",
		"text":               "string",
		"longtext":           "string",
		"blob":               "string",
		"tinyblob":           "string",
		"mediumblob":         "string",
		"longblob":           "string",
		"date":               "time.Time",
		"datetime":           "time.Time",
		"timestamp":          "time.Time",
		"time":               "time.Time",
		"float":              "float64",
		"double":             "float64",
		"decimal":            "float64",
		"binary":             "string",
		"varbinary":          "string",
	}
)


///--------------------------------------------------
// 初始版本自动化代码工具
type AutoCoder struct {
	TplPath        string   `json:"tplPath"` // 模板文件dir
	StructName     string   `json:"structName"`
	TableName      string   `json:"tableName"`
	Abbreviation   string   `json:"abbreviation"` //缩写
	ImportTime     bool     `json:"importTime"`
	ImportGorm     bool     `json:"ImportGorm"`
	GoStructString string   `json:"goStructString"`
	Fields         []*Field `json:"fields"`
}

type Field struct {
	FieldName       string `json:"fieldName"`
	FieldType       string `json:"fieldType"`
	ColumnName      string `json:"columnName"`
	FieldSearchType string `json:"fieldSearchType"`
}

type tplData struct {
	template         *template.Template
	locationPath     string
	autoCodePath     string
	autoMoveFilePath string
}

func NewAutoCoder(tplPath, sql string, colSearchTypeMap map[string]string)  (*AutoCoder, error) {
	// parse sql
	statement, err := sqlparser.ParseStrictDDL(sql)
	if err != nil {
		return nil, err
	}
	stmt, ok := statement.(*sqlparser.DDL)
	if !ok {
		return nil, fmt.Errorf("sql is not a create statment")
	}

	// init autoCoder
	autoCoder := &AutoCoder{
		TplPath: tplPath,
		ImportTime: false,
		ImportGorm : false,
		Fields: make([]*Field, 0),
	}

	// convert to Go struct
	tableName := stmt.NewName.Name.String()
	autoCoder.TableName = tableName

	primaryIdxMap, uniqueIdxMap, idxMap := buildIdxMaps(stmt.TableSpec.Indexes)

	builder := strings.Builder{}
	structName := snakeToCamel(tableName)
	autoCoder.StructName = structName
	autoCoder.Abbreviation = strings.ToLower(structName[0:1]) + structName[1:]

	structStart := fmt.Sprintf("type %s struct { \n", structName)
	builder.WriteString(structStart)

	builder.WriteString("\tglobal.MODEL\n")
	for _, col := range stmt.TableSpec.Columns {
		columnType := col.Type.Type
		if col.Type.Unsigned {
			columnType += " unsigned"
		}

		goType := SqlGoTypeMap[columnType]

		fieldName := snakeToCamel(col.Name.String())

		searchType, _ := colSearchTypeMap[col.Name.String()]
//		fmt.Println(searchType)

		oneField := &Field{
			FieldName: fieldName,
			FieldType: goType,
			ColumnName: col.Name.String(),
			FieldSearchType: searchType,
		}
		autoCoder.Fields = append(autoCoder.Fields, oneField)

		if fieldName == "CreatedAt" || fieldName == "UpdatedAt" || fieldName == "DeletedAt" || fieldName == "Id"{
			continue
		}
		gormStr := buildGormStr(col, primaryIdxMap, uniqueIdxMap, idxMap)
		
		// common fieldName
		builder.WriteString(fmt.Sprintf("\t%s\t%-25s `json:\"%s\" form:\"%s\" gorm:\"%s\"`\n", fieldName, goType, col.Name.String(), col.Name.String(), gormStr))

	}
	builder.WriteString("}\n")
	autoCoder.GoStructString = builder.String()

	return autoCoder, nil
}
// inner funcs
func buildIdxMaps(indexList []*sqlparser.IndexDefinition) (primaryIdxMap map[string]string, uniqueIdxMap map[string]string, idxMap map[string]string) {
	primaryIdxMap = make(map[string]string)
	uniqueIdxMap = make(map[string]string)
	idxMap = make(map[string]string)
	for idx, _ := range indexList {
		if indexList[idx].Info.Primary {
			for cIdx, _ := range indexList[idx].Columns {
				primaryIdxMap[indexList[idx].Columns[cIdx].Column.String()] = indexList[idx].Info.Name.String()
			}
		} else if indexList[idx].Info.Unique {
			for cIdx, _ := range indexList[idx].Columns {
				uniqueIdxMap[indexList[idx].Columns[cIdx].Column.String()] = indexList[idx].Info.Name.String()
			}
		} else {
			for cIdx, _ := range indexList[idx].Columns {
				idxMap[indexList[idx].Columns[cIdx].Column.String()] = indexList[idx].Info.Name.String()
			}
		}
	}
	return
}

func buildGormStr(col *sqlparser.ColumnDefinition, primaryIdxMap map[string]string, uniqueIdxMap map[string]string, idxMap map[string]string) string {
	builder := strings.Builder{}
	columnStr := fmt.Sprintf("column:%s", col.Name.String())
	builder.WriteString(columnStr)
	switch col.Type.Type {
	case "enum":
		enumBuilder := strings.Builder{}
		for idx, _ := range col.Type.EnumValues {
			if 0 == idx {
				enumBuilder.WriteString(col.Type.EnumValues[idx])
			} else {
				enumBuilder.WriteString("," + col.Type.EnumValues[idx])
			}
		}
		typeStr := fmt.Sprintf(";type:enum(%s)", enumBuilder.String())
		builder.WriteString(typeStr)
	default:
		if nil != col.Type.Length {
			switch int(col.Type.Length.Type) {
			case 1: // int
				typeStr := fmt.Sprintf(";type:%s(%s)", col.Type.Type, col.Type.Length.Val)
				builder.WriteString(typeStr)
			}
		} else {

			typeStr := fmt.Sprintf(";type:%s", col.Type.Type)
			builder.WriteString(typeStr)
		}
	}

	if col.Type.Unsigned {
		builder.WriteString(" unsigned")
	}

	if col.Type.Autoincrement {
		builder.WriteString(" auto_increment")
	}

	_, ok := primaryIdxMap[col.Name.String()]
	if ok {
		builder.WriteString(";primary_key")
	}
	_, ok = uniqueIdxMap[col.Name.String()]
	if ok {
		builder.WriteString(";unique")
	}


	if nil != col.Type.Default {
		defaultStr := ""

		if col.Type.Type == "string" {
			defaultStr = fmt.Sprintf(";default:'%s'", col.Type.Default.Val)
		}else{
			defaultStr = fmt.Sprintf(";default:%s", col.Type.Default.Val)
		}
		builder.WriteString(defaultStr)
	}

	if col.Type.NotNull {
		builder.WriteString(";not null")
	}

	idxName, ok := idxMap[col.Name.String()]
	if ok {
		indexStr := fmt.Sprintf(";index:%s", idxName)
		builder.WriteString(indexStr)
	}

	if nil != col.Type.Comment {
		commentStr := fmt.Sprintf(";comment:'%s'", col.Type.Comment.Val)
		builder.WriteString(commentStr)
	}

	return builder.String()
}

// In sql, table name often is snake_case
// In Go, struct name often is camel case
func snakeToCamel(str string) string {
	builder := strings.Builder{}
	index := 0
	if str[0] >= 'a' && str[0] <= 'z' {
		builder.WriteByte(str[0] - ('a' - 'A'))
		index = 1
	}
	for i := index; i < len(str); i++ {
		if str[i] == '_' && i+1 < len(str) {
			if str[i+1] >= 'a' && str[i+1] <= 'z' {
				builder.WriteByte(str[i+1] - ('a' - 'A'))
				i++
				continue
			}
		}
		builder.WriteByte(str[i])
	}
	return builder.String()
}

func getAllTplFile(pathName string, fileList []string) ([]string, error) {
	files, err := ioutil.ReadDir(pathName)
	for _, fi := range files {
		if fi.IsDir() {
			fileList, err = getAllTplFile(pathName+"/"+fi.Name(), fileList)
			if err != nil {
				return nil, err
			}
		} else {
			if strings.HasSuffix(fi.Name(), ".tpl") {
				fileList = append(fileList, pathName+"/"+fi.Name())
			}
		}
	}
	return fileList, err
}
// inner funcs

func (t *AutoCoder)CreateTemp() (err error) {
	dataList, needMkdir, err := t.getNeedList()
	if err != nil {
		return err
	}
	// 写入文件前，先创建文件夹
	if err = utils.CreateDir(needMkdir...); err != nil {
		return err
	}
	// 生成文件
	for _, value := range dataList {
		f, err := os.OpenFile(value.autoCodePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if err != nil {
			return err
		}
		if err = value.template.Execute(f, t); err != nil {
			return err
		}
		_ = f.Close()
	}
	return nil
}

func (t *AutoCoder)getNeedList() (dataList []tplData, needMkDirs []string, err error) {
	// 去除所有空格
	utils.TrimSpace(t)
	for _, field := range t.Fields {
		utils.TrimSpace(field)
	}
	// 获取 basePath 文件夹下所有tpl文件
	tplFileList, err := getAllTplFile(t.TplPath, nil)
	if err != nil {
		return nil, nil, err
	}
	dataList = make([]tplData, 0, len(tplFileList))
	needMkDirs = make([]string, len(tplFileList), len(tplFileList))
	// 根据文件路径生成 tplData 结构体，待填充数据
	for _, value := range tplFileList {
		dataList = append(dataList, tplData{locationPath: value})
	}
	// 生成 *Template, 填充 template 字段
	for index, value := range dataList {
		dataList[index].template, err = template.ParseFiles(value.locationPath)
		if err != nil {
			return nil, nil, err
		}
	}

	for index, value := range dataList {
		if strings.Contains(value.locationPath, "router") {
			dataList[index].autoCodePath = "router/" + t.TableName + ".go"
			needMkDirs[index] = "router"
		} else if strings.Contains(value.locationPath, "model") {
			dataList[index].autoCodePath = "model/" + t.TableName + ".go"
			needMkDirs[index] = "model"
		} else if strings.Contains(value.locationPath, "api") {
			dataList[index].autoCodePath = "api/" + t.TableName + ".go"
			needMkDirs[index] = "api"
		} else if strings.Contains(value.locationPath, "service") {
			dataList[index].autoCodePath = "service/" + t.TableName + ".go"
			needMkDirs[index] = "service"
		} else if strings.Contains(value.locationPath, "request") {
			dataList[index].autoCodePath = "model/request/" + t.TableName + ".go"
			needMkDirs[index] = "model/request"
		}
	}
	return dataList, needMkDirs, err
}

func main() {
	var inputFile string
	var tplPath string
	var searchTypes string

	flag.StringVar(&inputFile, "sql", "", "input sql file")
	flag.StringVar(&tplPath, "tpl_path", "", "tpl path")
	flag.StringVar(&searchTypes, "search_types", "", "search types")
	flag.Parse()

	if 0 == len(inputFile) || 0 == len(tplPath) {
		fmt.Println("some arg is empty!")
		return
	}
	fmt.Println("input sql: "+inputFile)
	fmt.Println("tpl path: "+tplPath)
	fmt.Println("search_types: "+searchTypes)

	global.LOG = core.Zap()

	colSearchTypeMap := make(map[string]string)
	searchList := strings.Split(searchTypes, ",")
	for _, elem := range searchList {
		elemList := strings.Split(elem, ":")
		if 2 == len(elemList) {
			colSearchTypeMap[elemList[0]] = elemList[1]
		}
	}

	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("read file err:%s\n", err.Error())
		return
	}
	sqlStmt := string(data)

	autoCoder, err := NewAutoCoder(tplPath, sqlStmt, colSearchTypeMap)
	if nil != err {
		panic(err)
	}
	err = autoCoder.CreateTemp()
	if nil != err {
		panic(err)
	}
}