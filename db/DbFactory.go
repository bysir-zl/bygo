package db

type DbFactory struct {
	dbConfigs map[string]DbConfig
}

var BFactory DbFactory

func (p *DbFactory) Model(m interface{}) (modelFactory *ModelFactory) {
	modelFactory = newModel(p.dbConfigs, m)
	return
}

func (p *DbFactory) SetConfig(c map[string]DbConfig) {
	p.dbConfigs = c
	return
}

// orm 入口, init后就可以全局使用BFactory
func Init(dbConfigs map[string]DbConfig) DbFactory {
	BFactory.dbConfigs = dbConfigs
	return BFactory
}

func NewDbFactory(dbConfigs map[string]DbConfig) DbFactory {
	return DbFactory{
		dbConfigs: dbConfigs,
	}
}
