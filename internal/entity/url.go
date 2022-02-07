package entity

type URL struct {
	Short string
	Full  string
}

func NewURL(full string) *URL {
	short := createShort(full)
	return &URL{
		Short: short,
		Full:  full,
	}
}

func createShort(full string) string {
	return "short_" + full
}
