package api

type Methods struct {
	typo  Typo
	from  int `json:"from"`
	to    int `json:"to"`
	year  int `json:"year"`
	month int `json:"moth"`
	day   int `json:"day"`
}

func NewMethods(From, To, Year, Month, Day int, typo Typo) *Methods {
	return &Methods{
		typo:  typo,
		from:  From,
		to:    To,
		year:  Year,
		month: Month,
		day:   Day,
	}
}

func (m *Methods) From() int {
	return m.from
}

func (m *Methods) To() int {
	return m.to
}

func (m *Methods) Typo() Typo {
	return m.typo
}

func (m *Methods) Year() int {
	return m.year
}

func (m *Methods) Month() int {
	return m.month
}

func (m *Methods) Day() int {
	return m.day
}
