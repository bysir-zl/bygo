package db

type DbFactory struct {
    dbConfigs map[string]DbConfig
}

func (p *DbFactory)Model(m interface{}) (modelFactory *ModelFactory) {
    modelFactory = NewModel(p.dbConfigs, m);
    return
}

func NewDbFactory(dbConfigs map[string]DbConfig) (DbFactory) {
    return DbFactory{
        dbConfigs:dbConfigs,
    }
}
