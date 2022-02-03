package database

type DummyImpl struct{}

var _ DatabaseProvider = (*DummyImpl)(nil)

func (DummyImpl) GetData() string {
	return "Data, yeah!"
}

func (DummyImpl) GetPrefix() string {
	return "di: "
}
