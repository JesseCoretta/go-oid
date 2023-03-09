package oid

import (
	"encoding/asn1"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
ObjectIdentifier facilitates the storage, and varied representation of, an ASN.1 object identifier in
a manner that goes beyond mere dotNotation and may be more convenient than using the asn1.ObjectIdentifier instance.
*/
type ObjectIdentifier struct {
	nANF, aka []string
	ints      []int
}

/*
ASN1 returns a populated instance of asn1.ObjectIdentifier using the contents of the receiver.
*/
func (o ObjectIdentifier) ASN1() asn1.ObjectIdentifier {
	return asn1.ObjectIdentifier(o.ints)
}

/*
Equal returns a boolean indicative of whether the provided type instance effectively matches the receiver.

This method supports asn1.ObjectIdentifier, []int, string and []string type instances for comparison. In the case of string input, a dotNotation match is attempted first, followed by an ASN.1 NameAndNumberForm sequence match and lastly a case folded string match of any alternative names by which the OID may be known.
*/
func (o ObjectIdentifier) Equal(x any) bool {
	switch tv := x.(type) {
	case asn1.ObjectIdentifier:
		return intSliceEqual([]int(tv), o.ints)
	case []int:
		return intSliceEqual(tv, o.ints)
	case string:
		if o.ASN1().String() == tv {
			// dotNotation
			return true
		} else if o.String() == tv {
			// ASN.1 NameAndNumberForm sequence
			return true
		} else {
			// alt names
			for i := 0; i < len(o.aka); i++ {
				if strings.EqualFold(o.aka[i], tv) {
					return true
				}
			}
		}
	case []string:
		return strSliceEqual(o.nANF, tv)
	}

	return false
}

/*
String returns the ASN.1 NameAndNumberForm sequence stored within the receiver in full, e.g.:

	{ iso(1) identified-organization(3) dod(6) }
*/
func (o ObjectIdentifier) String() string {
	return `{ ` + strings.Join(o.nANF, ` `) + ` }`
}

/*
NameAndNumberForm returns the effective ASN.1 NameAndNumberForm arc value for the receiver as an instance of [2]string. The first slice contains the identifier name (and is optional), while the second slice contains the integer string representation and is absolutely required in all cases.

For example, for { iso(1) identified-organization(3) dod(6) }, this method would return [2]string{`dod`, `6`}.
*/
func (o ObjectIdentifier) NameAndNumberForm() (nanf [2]string) {
	if len(o.nANF) == 0 || len(o.ints) == 0 {
		return
	}

	number := strconv.Itoa(o.ints[len(o.ints)-1])
	nanf = [2]string{"", number}
	name := o.nANF[len(o.nANF)-1]
	idx := strings.IndexRune(name, '(')
	if idx == -1 {
		// number only is fine.
		return
	}

	// we found a name, so add it to the payload
	nanf[0] = name[:idx]

	return
}

/*
IsZero checks the receiver for nilness and returns a boolean indicative of the result.
*/
func (oid *ObjectIdentifier) IsZero() bool {
	return oid == nil
}

/*
Valid returns a boolean value indicative of whether the receiver's length is greater than or equal to one (1) slice member.
*/
func (o ObjectIdentifier) Valid() bool {
	if o.IsZero() {
		return false
	}

	if len(o.nANF) != len(o.ints) {
		return false
	}

	// no negative numbers!
	for _, arc := range o.ints {
		if arc < 0 {
			return false
		}
	}

	// If the first arc is 0, 1 or 2,
	// then we passed verification.
	return 0 <= o.ints[0] && o.ints[0] <= 2
}

/*
oidToIntSlices converts a variety of dotNotation-based types into an []int instance.
*/
func oidToIntSlices(oid any) []int {
	switch tv := oid.(type) {
	case *ObjectIdentifier:
		return oidToIntSlices(tv.ints)
	case ObjectIdentifier:
		return oidToIntSlices(tv.ints)
	case string:
		arcs := strings.Split(tv, `.`)
		O := make([]int, len(arcs), len(arcs))
		for idx, arc := range arcs {
			u, err := strconv.Atoi(arc)
			if err != nil {
				return []int{}
			}
			O[idx] = u
		}
		return O
	case asn1.ObjectIdentifier:
		return []int(tv)
	case []int:
		return tv
	}

	return []int{}
}

/*
SetAltNames assigns alternative names by which the receiver may be known in the wild in addition to its "principal" name. Duplicates are filtered out.

One example of an alternate name in the wild is the OID `id-kp-serverAuth(1)` (1.3.6.1.5.5.7.3.1), which is also known simply as 'serverAuth'.
*/
func (o *ObjectIdentifier) SetAltNames(name ...string) {
	for i := 0; i < len(name); i++ {
		for j := 0; i < len(o.aka); i++ {
			if strings.EqualFold(name[j], o.aka[i]) {
				continue
			}
		}
		o.aka = append(o.aka, name[i])
	}

	return
}

/*
AltNames returns slices of string values, each representing an alternate name by which the receiver OID may be known in the wild.
*/
func (o ObjectIdentifier) AltNames() []string { return o.aka }

/*
NewObjectIdentifier creates an instance of ObjectIdentifier and returns it alongside an error.

The correct raw input syntax is the ASN.1 NameAndNumberForm sequence syntax, i.e.:

	{ iso(1) identified-organization(3) dod(6) }

Not all NameAndNumberForm values (arcs) require actual names; they can be numbers alone or in the so-called nameAndNumber syntax (name(Number)). For example:

	{ iso(1) identified-organization(3) 6 }

... is perfectly valid but generally not recommended when clarity is desired.
*/
func NewObjectIdentifier(raw string) (o *ObjectIdentifier, err error) {
	if len(raw) < 6 {
		err = errorf("Provide the proper ASN.1 notation (e.g.: '{ iso(1) org(3) dod(6) }')")
		return
	}

	t := new(ObjectIdentifier)
	o = new(ObjectIdentifier)
	f := strings.Fields(strings.TrimRight(strings.TrimLeft(raw, `{ `), ` }`))

	for i := 0; i < len(f); i++ {
		if len(f[i]) == 0 {
			err = errorf("Bad ASN.1 notation field '%s'", f[i])
			return
		}

		switch isDigit(f[i]) {
		case true:
			// we know its a number, so no need to check error
			num, _ := strconv.Atoi(f[i])
			t.ints = append(t.ints, num)
		default:
			// last char should be closing paren [ ")" ]
			idxr := strings.IndexRune(f[i], ')')
			if idxr != len(f[i])-1 {
				err = errorf("Bad identifier '%s' (raw: %#v)", f[i], f)
				return
			}

			// check everything before opening paren [ "(" ]
			idxl := strings.IndexRune(f[i], '(')
			for c := 0; c < idxl; c++ {
				ch := rune(f[i][c])
				if ('a' <= ch && ch <= 'z') ||
					('A' <= ch && ch <= 'Z') ||
					('0' <= ch && ch <= '9') || ch == '-' {
					// cool
				} else {
					err = errorf("Bad identifier '%s' at char #%d [%c]", f[i], c, ch)
					return
				}
			}

			// check everything between parens to ensure digits only [ "(%d)" ]
			if !isDigit(f[i][idxl+1 : idxr-1]) {
				err = errorf("Bad identifier '%s' at chars #%d[%c] through #%d[%c]", f[i], idxl, f[i][idxl], idxr, f[i][idxr])
				return
			}

			// we know its a number, so no need to check error
			num, _ := strconv.Atoi(f[i][idxl+1 : idxr])
			t.ints = append(t.ints, num)
		}

		// Seems safe to append current slice
		t.nANF = append(t.nANF, f[i])
	}

	if !t.Valid() {
		err = errorf("%T instance did not pass validity checks: %#v", t, *t)
		return
	}

	*o = *t

	return
}

/*
quick check to see if a string is effectively an integer.
*/
func isDigit(val string) bool {
	for _, c := range val {
		if '0' <= c && c <= '9' {
			continue
		}
		return false
	}
	return true
}

func errorf(msg any, x ...any) error {
	switch tv := msg.(type) {
	case string:
		return errors.New(fmt.Sprintf(tv, x...))
	case error:
		return errors.New(fmt.Sprintf(tv.Error(), x...))
	}

	return nil
}

/*
compare slice members of two (2) []int instances.
*/
func intSliceEqual(s1, s2 []int) (equal bool) {
	if len(s1)|len(s2) == 0 || len(s1) != len(s2) {
		return
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return
		}
	}

	equal = true
	return
}

/*
compare slice members of two (2) []string instances.
*/
func strSliceEqual(s1, s2 []string) (equal bool) {
	if len(s1)|len(s2) == 0 || len(s1) != len(s2) {
		return
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return
		}
	}

	equal = true
	return
}

