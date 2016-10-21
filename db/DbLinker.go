package db

//连接(join)
import (
    "bygo/util"
)

func LinkMapListToObjList(dbConfigs map[string]DbConfig,data []map[string]interface{}, out interface{}, pk string, outPk string) (err error) {
    ids := []interface{}{}
    for _, item := range data {
        ids = append(ids, item[pk])
    }

    err = NewModel(dbConfigs,out).
        WhereIn(outPk, ids...).
        Get()

    return
}

func LinkObjListToObjList(dbConfigs map[string]DbConfig,data interface{}, out interface{}, pk string, outPk string,useTag string) (err error) {
    mapers := util.ObjListToMapList(data,useTag);
    err = LinkMapListToObjList(dbConfigs,mapers, out, pk, outPk);
    return
}