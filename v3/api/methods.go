package api

type MethodsOpts func(m *Methods) error

type Methods struct {
	Type  Type `json:"typo"`
	From  int  `json:"from"`
	To    int  `json:"to"`
	Year  int  `json:"year"`
	Month int  `json:"moth"`
	Day   int  `json:"day"`
}

func NewMethods(opts ...MethodsOpts) (*Methods, error) {
	mts := &Methods{}
	for _, op := range opts {
		err := op(mts)
		if err != nil {
			return nil, err
		}
	}
	return mts, nil
}

func (m *Methods) GetFrom() int {
	return m.From
}

func (m *Methods) GetTo() int {
	return m.To
}

func (m *Methods) GetType() Type {
	return m.Type
}

func (m *Methods) GetYear() int {
	return m.Year
}

func (m *Methods) GetMonth() int {
	return m.Month
}

func (m *Methods) GetDay() int {
	return m.Day
}

func OptsDay(value int) MethodsOpts {
	return func(m *Methods) error {
		m.Day = value
		return nil
	}
}

func OptsMonth(value int) MethodsOpts {
	return func(m *Methods) error {
		m.Month = value
		return nil
	}
}

func OptsYear(value int) MethodsOpts {
	return func(m *Methods) error {
		m.Year = value
		return nil
	}
}

func OptsType(value Type) MethodsOpts {
	return func(m *Methods) error {
		m.Type = value
		return nil
	}
}

func OptsTo(value int) MethodsOpts {
	return func(m *Methods) error {
		m.To = value
		return nil
	}
}

func OptsFrom(value int) MethodsOpts {
	return func(m *Methods) error {
		m.From = value
		return nil
	}
}
