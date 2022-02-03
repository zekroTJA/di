package main

import (
	"github.com/zekrotja/di"
	"github.com/zekrotja/di/example/basic/services/database"
	"github.com/zekrotja/di/example/basic/services/printer"
	"github.com/zekrotja/di/example/basic/services/program"
)

func main() {
	c := di.NewContainer()

	di.MustRegister[database.DatabaseProvider, database.DummyImpl](c)
	di.MustRegister[printer.PrinterProvider, printer.LogImpl](c)
	di.MustRegister[program.ProgramProvider, program.ProgramImpl](c)

	p := di.MustGet[program.ProgramProvider](c)

	p.Run()
}
