package db

type DbFactory struct {
    dbConfigs map[string]DbConfig
}

var Factory = DbFactory{}

func (p *DbFactory)Model(m interface{}) (modelFactory *ModelFactory) {
    modelFactory = NewModel(p.dbConfigs, m)
    return
}

func (p *DbFactory)SetConfig(c map[string]DbConfig) {
    p.dbConfigs = c
    return
}

func NewDbFactory(dbConfigs map[string]DbConfig) (DbFactory) {
    return DbFactory{
        dbConfigs:dbConfigs,
    }
}

