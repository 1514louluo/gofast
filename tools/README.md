##使用方法

###编译
    go build auto_code.go

###使用
    ./auto_code -sql: 指定sql文件路径
                -tpl_path: 指定模板文件路径, 默认是./tpl
                -search_types: 指定某个字段搜索是按=来查询，还是用LIKE模糊查询，格式 band:=,band_pinyin:LIKE,band_name:LIKE

###自动生成代码
    sh gen.sh

###文件移动
    1. 将当前目录生成的api、service、router、model、model/request对应的文件拷贝到tb_server对应目录下文件即可
    2. 注意覆盖的文件注意备份