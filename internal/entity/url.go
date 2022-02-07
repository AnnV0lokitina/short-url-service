package entity

type Url struct {
	Short string
	Full  string
}

func NewUrl(full string) *Url {
	short := createShort(full)
	return &Url{
		Short: short,
		Full:  full,
	}
}

func createShort(full string) string {
	return "short_" + full
}
