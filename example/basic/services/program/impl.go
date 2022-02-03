package program

import (
	"github.com/zekrotja/di/example/basic/services/database"
	"github.com/zekrotja/di/example/basic/services/printer"
)

type ProgramImpl struct {
	Db database.DatabaseProvider
	Pr printer.PrinterProvider
}

var _ ProgramProvider = (*ProgramImpl)(nil)

func (p ProgramImpl) Run() {
	data := p.Db.GetData()
	p.Pr.Print(data)
}
