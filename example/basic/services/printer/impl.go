package printer

import (
	"log"

	"github.com/zekrotja/di/example/basic/services/database"
)

type LogImpl struct {
	Db database.DatabaseProvider
}

var _ PrinterProvider = (*LogImpl)(nil)

func (p *LogImpl) Print(v any) {
	log.Printf("%s%s", p.Db.GetPrefix(), v)
}
