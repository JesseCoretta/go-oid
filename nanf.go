package oid

/*
nanf.go deals with NameAndNumberForm syntax and viability
*/

type NameAndNumberForm struct {
	identifier        string
	primaryIdentifier uint
}

func (nanf NameAndNumberForm) IsZero() bool {
	return len(nanf.identifier)&int(nanf.primaryIdentifier) == 0
}

func (nanf NameAndNumberForm) Identifier() string {
	return nanf.identifier
}

func (nanf NameAndNumberForm) Decimal() int {
	return int(nanf.primaryIdentifier)
}

func (nanf NameAndNumberForm) String() (val string) {
	n := itoa(int(nanf.primaryIdentifier))
	if len(nanf.identifier) == 0 {
		return n
	}
	return sprintf("%s(%s)", nanf.identifier, n)
}

func (nanf NameAndNumberForm) Equal(n NameAndNumberForm) bool {
	return eq(nanf.identifier, n.identifier) &&
		nanf.primaryIdentifier == n.primaryIdentifier
}

func parseNaNFstr(x string) (nanf *NameAndNumberForm, err error) {
	if len(x) == 0 {
		err = errorf("No content for parseNaNFstr to read")
		return
	} else if x[len(x)-1] != ')' {
		err = errorf("No closing parenthesis for parseNaNFstr to read")
		return
	}

	idx := indexRune(x, '(')
	if idx == -1 {
		err = errorf("No opening parenthesis for parseNaNFstr to read")
		return
	}
	nanf = new(NameAndNumberForm)

	n := x[idx+1 : len(x)-1]
	if !isDigit(n) {
		err = errorf("Bad primaryIdentifier '%s'", n)
		return
	}

	f, _ := atoi(n)
	nanf.primaryIdentifier = uint(f)

	for c := 0; c < len(x[:idx-1]); c++ {
		ch := rune(x[c])

		if c == 0 {
			if !('a' <= ch && ch <= 'z') {
				err = errorf("Bad identifier '%s' at char #%d [%c] [hint: must only start with lowercase alpha]", x[:idx-1], c, ch)
				return
			}
		}

		if ('a' <= ch && ch <= 'z') ||
			('A' <= ch && ch <= 'Z') ||
			('0' <= ch && ch <= '9') || ch == '-' {
			// cool
		} else {
			err = errorf("Bad identifier '%s' at char #%d [%c], unsupported character(s) [hint: must be A-Z, a-z, 0-9 or '-']", x[:idx-1], c, ch)
			return
		}

		if c == idx-1 {
			if ch == '-' {
				err = errorf("Bad identifier '%s' at char #%d [%c] [hint: final identifier character cannot be a hyphen]", x[:idx-1], c, ch)
			}
		}
	}

	// identifier seems safe to assign
	nanf.identifier = x[:idx]
	return
}

func NewNameAndNumberForm(x any) (nanf *NameAndNumberForm, err error) {

	switch tv := x.(type) {
	case string:
		if !isDigit(tv) {
			nanf, err = parseNaNFstr(tv)
		} else {
			z, _ := atoi(tv)
			nanf, err = NewNameAndNumberForm(uint(z))
		}
	case uint:
		nanf = new(NameAndNumberForm)
		nanf.primaryIdentifier = tv
	case int:
		if tv < 0 {
			err = errorf("primaryIdentifier cannot be negative")
		} else {
			nanf, err = NewNameAndNumberForm(uint(tv))
		}
	default:
		err = errorf("Unsupported NameAndNumberForm input type '%T'", tv)
	}

	return
}
