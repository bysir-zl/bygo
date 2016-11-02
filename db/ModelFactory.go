package db

import (
	"strings"
	"fmt"
	"reflect"
	"errors"
	"log"
	"regexp"
	"math"
	"time"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/bean"
)

type orderItem struct {
	Field string
	Desc  string
}

type ModelFactory struct {
	dbConfigs     map[string]DbConfig

	out           interface{} // 要输出的Model,可以是slice也可以是单个model
	modelFieldMap util.FieldTagMapper

	table         string
	connect       DbConfig
	debug         bool

	fields        []string
	where         map[string]([]interface{})
	limit         [2]int
	order         []orderItem
}

func (p *ModelFactory) Field(fields ...string) *ModelFactory {
	p.fields = fields
	return p
}

func (p *ModelFactory) Limit(skip int, len int) *ModelFactory {
	p.limit = [2]int{skip, len}
	return p
}

func (p *ModelFactory) OrderBy(field string, desc bool) *ModelFactory {
	if p.order == nil {
		p.order = []orderItem{}
	}
	item := orderItem{}
	item.Field = field
	if desc {
		item.Desc = "DESC"
	} else {
		item.Desc = ""
	}

	p.order = append(p.order, item)
	return p
}
func (p *ModelFactory) Debug(debug bool) *ModelFactory {
	p.debug = debug
	return p
}

func (p *ModelFactory) Where(condition string, values ...interface{}) *ModelFactory {
	if p.where == nil {
		p.where = map[string]([]interface{}){}
	}
	p.where[condition] = values

	return p
}

//清除已经有的Where条件
func (p *ModelFactory) reSet() *ModelFactory {
	p.where = map[string]([]interface{}){}
	p.fields = []string{}
	p.limit = [2]int{0, 0}
	p.order = []orderItem{}
	return p
}

func (p *ModelFactory) WhereIn(field string, params ...interface{}) *ModelFactory {
	if p.where == nil {
		p.where = map[string]([]interface{}){}
	}
	dataHolder := strings.Repeat(",?", len(params))
	dataHolder = dataHolder[1:]

	p.where[field + " IN (" + dataHolder + ")"] = params
	return p
}


//指定连接哪个数据库
//connect 必须是已经是在config/db配置好了的
func (p *ModelFactory) Db(connect string) *ModelFactory {
	c := p.dbConfigs[connect]
	if c.Port == 0 {
		err := errors.New("the connect `" + connect + "` is undefined in dbConfig")
		panic(err)
	}
	p.connect = c
	return p
}

//指定连接哪个数据库
func (p *ModelFactory) DbConfig(connect string, config DbConfig) *ModelFactory {

	p.dbConfigs[connect] = config
	p.connect = config

	return p
}
//指定表
func (p *ModelFactory) Table(table string) *ModelFactory {
	p.table = table
	return p
}

//原始方法 查询sql返回map
func (p *ModelFactory) Query(sql string, args ...interface{}) (data []map[string]interface{}, err error) {
	data = nil
	dbDriver, err := Singleton(p.connect)

	if err != nil {
		return
	}
	out, err := dbDriver.Query(sql, args...)
	if err != nil {
		return
	}
	data = out

	return
}

//查询主方法,返回[]Map原数据给get和first使用
func (p *ModelFactory) QueryToMap() (data []map[string]interface{}, err error) {

	//字段->数据库字段映射
	fieldTagMap := p.modelFieldMap.GetFieldMapByTagName("name")
	modelConfig := p.modelFieldMap.GetFieldMapByTagName("db")

	//检查fields是否在Model里
	if p.fields != nil&&len(p.fields) != 0 {
		if isOk, msg := util.ArrayInMapKey(p.fields, fieldTagMap); !isOk {
			err = errors.New("filed`s `" + msg + "` is not in the model fields")
			return
		}
	}

	if p.table == "" {
		err = errors.New("the model has not `Table` field or Tag.name")
		return
	}

	//没有指定链接
	if p.connect.Port == 0 {
		dbConnect := modelConfig["Connect"];
		if dbConnect == "" {
			dbConnect = "default"
		}
		p.connect = p.dbConfigs[dbConnect]
		//没有找到connect配置
		if p.connect.Port == 0 {
			err = errors.New("the connect `" + dbConnect + "` is undefined in dbConfig")
			return
		}
	}

	sql, args, e := buildSelectSql(p.fields, p.table, p.where, p.order, p.limit, fieldTagMap);
	if e != nil {
		err = e
		return
	}

	if p.debug {
		log.Println("p query sql is : " + sql, " ", args)
	}

	dataMap, err2 := p.Query(sql, args...)
	if err2 != nil {
		err = err2
		return
	}
	if len(dataMap) == 0 {

		return
	}

	// {'report_id':'Id','report_content':'Content'}
	//fieldsMap2 := Util.ReverseMap(fieldTagMap)

	//dms := []map[string]interface{}{}
	//for _, md := range dataMap {
	//    // 根据fieldsMap,将db field => struct field
	//    dm := map[string]interface{}{}
	//    for key, value := range md {
	//        k2 := fieldsMap2[key]
	//        // 在struct字段中,没有 查询出来的dbMap 上的key
	//        if k2 == "" {
	//            continue
	//        }
	//        dm[k2] = value
	//    }
	//    dms = append(dms, dm)
	//}
	data = dataMap

	return
}

func (p *ModelFactory)GetAndLink(out interface{}, pk string, outPk string) (err error) {

	data, e := p.QueryToMap()
	if e != nil {
		err = e
		return
	}

	util.MapListToObjList(p.out, data, "name")
	LinkMapListToObjList(p.dbConfigs, data, out, pk, outPk)

	return
}

//查询返回一个数组
func (p *ModelFactory) Get() (err error) {
	//从数组interface中获取一个元素
	mo := reflect.ValueOf(p.out).Type().Elem()

	if mo.String()[0] != '[' {
		err = errors.New("Get function need one slice(model) param")
		return
	}
	data, err := p.QueryToMap()
	if err != nil {
		return
	}

	util.MapListToObjList(p.out, data, "name")
	return
}

func (p *ModelFactory) Page(page int, pageSize int) (pageData bean.Page, err error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 1
	} else if pageSize > 200 {
		pageSize = 40
	}

	p.limit = [2]int{(page - 1) * pageSize, pageSize}

	data, err := p.QueryToMap()
	if err != nil {
		return
	}

	util.MapListToObjList(p.out, data, "name")

	count, err := p.Count();
	if err != nil {
		return
	}
	pageTotal := int(math.Ceil(float64(count) / float64(pageSize)))
	pageData = bean.Page{Total:count, Page:page, PageSize:pageSize, PageTotal:pageTotal}

	return
}

func (p *ModelFactory)PageWithOutTotal(page int, pageSize int) (pageData bean.Page, err error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 1
	} else if pageSize > 200 {
		pageSize = 40
	}

	p.limit = [2]int{(page - 1) * pageSize, pageSize}

	data, err := p.QueryToMap()
	if err != nil {
		return
	}

	util.MapListToObjList(p.out, data, "name")

	pageData = bean.Page{Page:page, PageSize:pageSize}

	return
}

func (p *ModelFactory) First() (err error) {

	p.limit = [2]int{0, 1}
	datas, err := p.QueryToMap();

	if err != nil {
		return
	}
	if datas == nil || len(datas) == 0 {
		return
	}

	util.MapToObj(p.out, datas[0], "name")

	return
}

func (p *ModelFactory) Count() (count int64, err error) {

	//字段->数据库字段映射
	fieldTagMap := p.modelFieldMap.GetFieldMapByTagName("name")
	modelConfig := p.modelFieldMap.GetFieldMapByTagName("db")

	if p.table == "" {
		err = errors.New("the struct `" + reflect.TypeOf(p.out).String() + "` has not `Table` field or Tag.name")
		return
	}

	//没有指定链接
	if p.connect.Port == 0 {
		dbConnect := modelConfig["Connect"];
		if dbConnect == "" {
			dbConnect = "default"
		}

		p.connect = p.dbConfigs[dbConnect]

		//没有找到connect配置
		if p.connect.Port == 0 {
			err = errors.New("the `" + reflect.TypeOf(p.out).String() + "` connect `" + dbConnect + "` is undefined in dbConfig")
			return
		}
	}

	sql, args, e := buildCountSql(p.table, p.where, fieldTagMap)
	if e != nil {
		err = e
		return
	}

	if p.debug {
		log.Println("p query sql is : " + sql, " ", args)
	}

	data, e := p.Query(sql, args...)
	if e != nil {
		err = e
		return
	}
	if len(data) == 0 {
		err = errors.New("query sql success,but return len is 0")
	}
	count = data[0]["count"].(int64)

	return
}

func (p *ModelFactory) GetPk(types string) string {
	pk := ""
	pkMap := p.modelFieldMap.GetFieldMapByTagName("pk")
	if pkMap != nil&&len(pkMap) != 0 {
		for key, value := range pkMap {
			if value == types {
				pk = key
			}
		}
	}
	return pk
}

// 取得在method操作时需要自动填充的字段与值
func (p *ModelFactory) GetAutoSetField(method string) (needSet  map[string]interface{}, err error) {
	autoField := p.modelFieldMap.GetFieldMapByTagName("auto")
	if len(autoField) != 0 {
		needSet = map[string]interface{}{};
		for field, tagVal := range autoField {
			// time,insert
			tyAme := strings.Split(tagVal, ",");
			methods := tyAme[1];

			if util.ItemInArray(method, strings.Split(methods, "|")) {
				if tyAme[0] == "timestr" {
					needSet[field] = time.Now().Format("2006-01-02 15:04:05")
				} else if tyAme[0] == "timeint" {
					needSet[field] = time.Now().Unix()
				}
			}
		}
	}

	return
}

func (p *ModelFactory) Insert() (err error) {

	//字段->数据库字段映射
	fieldTagMap := p.modelFieldMap.GetFieldMapByTagName("name")
	modelConfig := p.modelFieldMap.GetFieldMapByTagName("db")

	//检查fields是否在Model里
	if p.fields != nil&&len(p.fields) != 0 {
		if isOk, msg := util.ArrayInMapKey(p.fields, fieldTagMap); !isOk {
			err = errors.New("filed`s `" + msg + "` is not in the struct `" + reflect.TypeOf(p.out).String() + "` fields")
			return
		}
	}

	if p.table == "" {
		err = errors.New("the struct `" + reflect.TypeOf(p.out).String() + "` has not `Table` field or Tag.name")
		return
	}

	//没有指定链接
	if p.connect.Port == 0 {
		dbConnect := modelConfig["Connect"];
		if dbConnect == "" {
			dbConnect = "default"
		}

		p.connect = p.dbConfigs[dbConnect]

		//没有找到connect配置
		if p.connect.Port == 0 {
			err = errors.New("the `" + reflect.TypeOf(p.out).String() + "` connect `" + dbConnect + "` is undefined in dbConfig")
			return
		}
	}

	//获取应该存数据库的键值对 数据库字段->值映射
	saveData := map[string]interface{}{}

	//字段->值映射
	mapper := util.ObjToMap(p.out, "")
	for key, value := range mapper {
		k := fieldTagMap[key]
		//没有 name 属性的字段,说明不是数据库字段 , 比如Table,Connect等配置字段
		if k == "" {
			continue
		}

		//指定了fields 就只更新指定字段
		if p.fields != nil {
			if !util.ItemInArray(key, p.fields) {
				continue
			}
		}
		//在插入的时候过滤空值
		if util.IsEmptyValue(value) {
			continue
		}
		saveData[k] = value
	}
	autoSet, e := p.GetAutoSetField("insert");
	if e != nil {
		err = e
		return
	}

	if autoSet != nil&&len(autoSet) != 0 {
		for k, v := range autoSet {
			k = fieldTagMap[k];
			saveData[k] = v
		}

		//将自动添加的字段附加到model里，方便返回
		util.MapToObj(p.out, autoSet, "")
	}

	sql, args, e := buildInsertSql(p.table, saveData, p.where);
	if e != nil {
		err = e
		return
	}

	if p.debug {
		log.Println("p query sql is : " + sql, " ", args)
	}

	_, _insertId, _err := p.Exec(sql, args...)
	if _err != nil {
		err = _err
		return
	}

	//找到主键，并且赋值为lastInsertId

	pk := p.GetPk("auto");
	if pk != "" {
		ma := map[string]interface{}{}
		ma[pk] = _insertId

		util.MapToObj(p.out, ma, "")
	}

	return
}

func (p *ModelFactory) Update() (count int64, err error) {
	if p.where == nil {
		err = errors.New("you need set condition in Where()")
		return
	}


	//字段->数据库字段映射
	fieldTagMap := p.modelFieldMap.GetFieldMapByTagName("name")
	modelConfig := p.modelFieldMap.GetFieldMapByTagName("db")

	//检查fields是否在Model里
	if p.fields != nil&&len(p.fields) != 0 {
		if isOk, msg := util.ArrayInMapKey(p.fields, fieldTagMap); !isOk {
			err = errors.New("filed`s `" + msg + "` is not in the struct `" + reflect.TypeOf(p.out).String() + "` fields")
			return
		}
	}

	if p.table == "" {
		err = errors.New("the struct `" + reflect.TypeOf(p.out).String() + "` has not `Table` field or Tag.name")
		return
	}

	//没有指定链接
	if p.connect.Port == 0 {
		dbConnect := modelConfig["Connect"];
		if dbConnect == "" {
			dbConnect = "default"
		}

		p.connect = p.dbConfigs[dbConnect]

		//没有找到connect配置
		if p.connect.Port == 0 {
			err = errors.New("the `" + reflect.TypeOf(p.out).String() + "` connect `" + dbConnect + "` is undefined in dbConfig")
			return
		}
	}



	//获取应该存数据库的键值对 数据库字段->值映射
	saveData := map[string]interface{}{}

	//获取主键
	pk := p.GetPk("auto")

	//字段->值映射
	mapper := util.ObjToMap(p.out, "")
	for key, value := range mapper {
		k := fieldTagMap[key]
		//没有 name 属性的字段,说明不是数据库字段 , 比如Table,Connect等配置字段
		//如果是主键,则跳过赋值
		if k == "" || key == pk {
			continue
		}
		//指定了fields 就只更新指定字段
		if p.fields != nil {
			if !util.ItemInArray(key, p.fields) {
				continue
			}
		}

		saveData[k] = value
	}

	autoSet, e := p.GetAutoSetField("update");
	if e != nil {
		err = e
		return
	}

	if autoSet != nil&&len(autoSet) != 0 {
		for k, v := range autoSet {
			k = fieldTagMap[k];
			saveData[k] = v
		}

		//将自动添加的字段附加到model里，方便返回
		util.MapToObj(p.out, autoSet, "")
	}

	sql, args, e := buildUpdateSql(p.table, saveData, p.where, fieldTagMap)
	if e != nil {
		err = e
		return
	}

	if p.debug {
		log.Println("p query sql is : " + sql, " ", args)
	}

	c, _, e := p.Exec(sql, args...)
	if e != nil {
		err = e
		return
	}
	count = c

	return
}
func (p *ModelFactory) Delete() (count int64, err error) {
	if p.where == nil {
		err = errors.New("you need set condition in Where()")
		return
	}


	//字段->数据库字段映射
	fieldTagMap := p.modelFieldMap.GetFieldMapByTagName("name")
	modelConfig := p.modelFieldMap.GetFieldMapByTagName("db")

	//检查fields是否在Model里
	if p.fields != nil&&len(p.fields) != 0 {
		if isOk, msg := util.ArrayInMapKey(p.fields, fieldTagMap); !isOk {
			err = errors.New("filed`s `" + msg + "` is not in the struct `" + reflect.TypeOf(p.out).String() + "` fields")
			return
		}
	}

	if p.table == "" {
		err = errors.New("the struct `" + reflect.TypeOf(p.out).String() + "` has not `Table` field or Tag.name")
		return
	}

	//没有指定链接
	if p.connect.Port == 0 {
		dbConnect := modelConfig["Connect"];
		if dbConnect == "" {
			dbConnect = "default"
		}

		p.connect = p.dbConfigs[dbConnect]

		//没有找到connect配置
		if p.connect.Port == 0 {
			err = errors.New("the `" + reflect.TypeOf(p.out).String() + "` connect `" + dbConnect + "` is undefined in dbConfig")
			return
		}
	}

	sql, args, e := buildDeleteSql(p.table, p.where, fieldTagMap)
	if e != nil {
		err = e
		return
	}

	if p.debug {
		log.Println("p query sql is : " + sql, " ", args)
	}

	c, _, e := p.Exec(sql, args...)
	if e != nil {
		err = e
		return
	}
	count = c

	return
}

func (p *ModelFactory)Exec(sql string, args ...interface{}) (affectCount int64, lastInsertId int64, err error) {

	dbDriver, err := Singleton(p.connect)
	if err != nil {
		return
	}
	att, insertId, err := dbDriver.Exec(sql, args...)
	if err != nil {
		return
	}

	lastInsertId = insertId
	affectCount = att

	return
}

func buildSelectSql(fields []string, tableName string, where map[string]([]interface{}), order []orderItem, limit [2]int, fieldMapper map[string]string) (sql string, args []interface{}, err error) {

	args = []interface{}{}
	sql = "SELECT "

	//field
	fieldString := ""
	if fields == nil || len(fields) == 0 {
		fieldString = "* "
	} else {
		//转换字段名
		for _, value := range fields {
			fieldString = fieldString + "," + fieldMapper[value]
		}
		fieldString = fieldString[1:];
	}

	sql = sql + fieldString + " "

	//table
	sql = sql + "FROM `" + tableName + "` "

	//where
	if where != nil {
		whereString, as, e := buildWhere(where, fieldMapper)
		if e != nil {
			err = e
			return
		}

		for _, a := range as {
			args = append(args, a)
		}

		sql = sql + "WHERE " + whereString + " "
	}

	//orderBy
	if order != nil {
		orderString := ""

		for _, value := range order {
			field := fieldMapper[value.Field]

			if field == "" {
				err = errors.New("the order field " + field + " is undefined in model")
				return
			}

			orderString = orderString + "," + field + " " + value.Desc
		}
		orderString = orderString[1:]

		sql = sql + "ORDER BY " + orderString + " "
	}

	//limit
	if limit[0] != 0 || limit[1] != 0 {
		sql = sql + "LIMIT " + fmt.Sprintf("%d,%d", limit[0], limit[1]) + " "
	}

	return
}

func buildInsertSql(tableName string, saveData map[string]interface{}, where map[string]([]interface{})) (sql string, args []interface{}, err error) {

	if len(saveData) == 0 {
		err = errors.New("no save data on INSERT")
		return
	}

	args = []interface{}{}
	sql = "INSERT INTO " + tableName + " ("

	fieldsStr := ""
	holderStr := ""

	for key, value := range saveData {
		fieldsStr = fieldsStr + ", " + key
		holderStr = holderStr + ", ?"

		args = append(args, value)
	}

	fieldsStr = fieldsStr[2:]
	holderStr = holderStr[2:]

	sql = sql + fieldsStr + " ) VALUES ( " + holderStr + " )"

	return
}

func buildUpdateSql(tableName string, saveData map[string]interface{}, where map[string]([]interface{}), fieldMapper map[string]string) (sql string, args []interface{}, err error) {

	if len(saveData) == 0 {
		err = errors.New("no save data on INSERT")
		return
	}

	args = []interface{}{}
	sql = "UPDATE " + tableName + " SET "

	//value
	fieldsStr := ""
	for key, value := range saveData {
		fieldsStr = fieldsStr + ", " + key + "= ?"

		args = append(args, value)
	}

	fieldsStr = fieldsStr[2:]
	sql = sql + fieldsStr + " "

	//where
	if where != nil {
		whereString, as, e := buildWhere(where, fieldMapper)
		if e != nil {
			err = e
			return
		}

		for _, a := range as {
			args = append(args, a)
		}
		sql = sql + "WHERE " + whereString + " "
	}

	return
}

func buildDeleteSql(tableName string, where map[string]([]interface{}), fieldMapper map[string]string) (sql string, args []interface{}, err error) {
	args = []interface{}{}
	sql = "DELETE FROM " + tableName + " "

	//where
	if where != nil {
		whereString, as, e := buildWhere(where, fieldMapper)
		if e != nil {
			err = e
			return
		}
		for _, a := range as {
			args = append(args, a)
		}
		sql = sql + "WHERE (" + whereString + ") "
	}

	return
}

func buildCountSql(tableName string, where map[string]([]interface{}), fieldMapper map[string]string) (sql string, args []interface{}, err error) {
	args = []interface{}{}
	sql = "SELECT COUNT(*) as count FROM " + tableName + " "

	//where
	if where != nil {
		whereString, as, e := buildWhere(where, fieldMapper)
		if e != nil {
			err = e
			return
		}
		for _, a := range as {
			args = append(args, a)
		}
		sql = sql + "WHERE (" + whereString + ") "
	}

	return
}

//生成where 条件
func buildWhere(where map[string]([]interface{}), fieldMapper map[string]string) (whereString string, args []interface{}, err error) {
	if where != nil {
		args = []interface{}{}
		whereString = " "

		for key, vaules := range where {
			whereString = whereString + " AND ( " + key + " )"
			for _, value := range vaules {
				args = append(args, value)
			}
		}

		whereString = whereString[5:]

		if fieldMapper != nil && len(fieldMapper) != 0 {

			//处理where的字段,但是必须在写条件的时候,必须在字段两边加上`符号 如:`Id` = ? AND `Name` = ? ,才能映射字段
			reg, _ := regexp.Compile("`(.+?)`")

			whereString = reg.ReplaceAllStringFunc(whereString, func(in string) string {
				k := string(in)[1:len(in) - 1] //去掉左右`号
				var ne string = fieldMapper[k]
				if ne == "" {
					err = errors.New("the where field(in '" + whereString + "') " + k + " is undefined in model")
				}
				return "`" + ne + "`"
			})
		}

	}

	return
}

func newModel(dbConfig map[string]DbConfig, m interface{}) *ModelFactory {
	modelFactory := &ModelFactory{
		dbConfigs:dbConfig,
	}

	if m != nil {
		modelFactory.out = m

		// 如果是slice
		// todo 这里可以用.Kind() 优化
		if (reflect.ValueOf(m).Type().Elem().String()[0] == '[') {
			mo := reflect.New(reflect.TypeOf(m).Elem().Elem()).Interface()
			modelFactory.modelFieldMap = util.GetTagMapperFromPool(mo)
		} else {
			modelFactory.modelFieldMap = util.GetTagMapperFromPool(m)
		}

		modelFactory.table = modelFactory.modelFieldMap.GetFieldMapByTagName("db")["Table"]
		modelFactory.connect = dbConfig[modelFactory.modelFieldMap.GetFieldMapByTagName("db")["Connect"]]
	} else {
		modelFactory.connect = dbConfig["default"]
	}

	return modelFactory;
}


