package di

type Strategy string

const (
	Singleton   = Strategy("singleton")
	Transistent = Strategy("transistent")
)

func getStrategy(v string) (Strategy, bool) {
	switch v {
	case "singleton", "single", "s":
		return Singleton, true
	case "transistent", "trans", "t":
		return Transistent, true
	default:
		return "", false
	}
}
