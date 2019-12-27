package redis_adapter

type Adapter struct {
}

func (a *Adapter) LoadPolicy(model interface{}) error {
	panic("implement me")
}

func (a *Adapter) SavePolicy(model interface{}) error {
	panic("implement me")
}

func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	panic("implement me")
}

func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	panic("implement me")
}

func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	panic("implement me")
}
