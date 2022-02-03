package printer

type PrinterProvider interface {
	Print(v any)
}
