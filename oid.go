package oid

import (
	"encoding/asn1"
	"strings"
)

/*
ObjectIdentifier facilitates the storage, and varied representation of, an ASN.1 object identifier in
a manner that goes beyond mere dotNotation and may be more convenient than using the asn1.ObjectIdentifier instance.
*/
type ObjectIdentifier struct {
	nANF []NameAndNumberForm
	aka  []string
}

/*
ASN1 returns a populated instance of asn1.ObjectIdentifier using the contents of the receiver.
*/
func (o ObjectIdentifier) ASN1() (a asn1.ObjectIdentifier) {
	a = make(asn1.ObjectIdentifier, len(o.nANF), len(o.nANF))
	for i := 0; i < len(o.nANF); i++ {
		a[i] = int(o.nANF[i].primaryIdentifier)
	}
	return
}

/*
Equal returns a boolean indicative of whether the provided type instance effectively matches the receiver.

This method supports asn1.ObjectIdentifier, []int, string and []string type instances for comparison. In the case of string input, a dotNotation match is attempted first, followed by an ASN.1 NameAndNumberForm sequence match and lastly a case folded string match of any alternative names by which the OID may be known.
*/
func (o ObjectIdentifier) Equal(x any) bool {
	switch tv := x.(type) {
	case asn1.ObjectIdentifier:
		return intSliceEqual([]int(tv), []int(o.ASN1()))
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
		if len(o.nANF) != len(tv) {
			return false
		}

		for i := 0; i < len(o.nANF); i++ {
			if o.nANF[i].String() != tv[i] {
				return false
			}
		}
		return true
	}

	return false
}

/*
String returns the ASN.1 NameAndNumberForm sequence stored within the receiver in full, e.g.:

	{ iso(1) identified-organization(3) dod(6) }
*/
func (o ObjectIdentifier) String() (a string) {
	a = `{`
	for i := 0; i < len(o.nANF); i++ {
		a += sprintf(" %s", o.nANF[i])
	}
	a += ` }`

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

	// If the first arc is 0, 1 or 2,
	// then we passed verification.
	return 0 <= o.nANF[0].primaryIdentifier && o.nANF[0].primaryIdentifier <= 2
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

func (o ObjectIdentifier) len() int { return len(o.nANF) }

func (o ObjectIdentifier) NameAndNumberForm() (nanf NameAndNumberForm) {
	if o.len() == 0 {
		return
	}

	return o.nANF[len(o.nANF)-1]
}

/*
NewObjectIdentifier creates an instance of ObjectIdentifier and returns it alongside an error.

The correct raw input syntax is the ASN.1 NameAndNumberForm sequence syntax, i.e.:

	{ iso(1) identified-organization(3) dod(6) }

Not all NameAndNumberForm values (arcs) require actual names; they can be numbers alone or in the so-called nameAndNumber syntax (name(Number)). For example:

	{ iso(1) identified-organization(3) 6 }

... is perfectly valid but generally not recommended when clarity is desired.
*/
func NewObjectIdentifier(x any) (o *ObjectIdentifier, err error) {
	t := new(ObjectIdentifier)

	switch tv := x.(type) {
	case string:
		f := fields(trimR(trimL(tv, `{ `), ` }`))
		for i := 0; i < len(f); i++ {
			var nanf *NameAndNumberForm
			if nanf, err = NewNameAndNumberForm(f[i]); err != nil {
				return
			}
			t.nANF = append(t.nANF, *nanf)
		}
	case []string:
		for i := 0; i < len(tv); i++ {
			var nanf *NameAndNumberForm
			if nanf, err = NewNameAndNumberForm(tv[i]); err != nil {
				return
			}
			t.nANF = append(t.nANF, *nanf)
		}
	case []int:
		for i := 0; i < len(tv); i++ {
			var nanf *NameAndNumberForm
			if nanf, err = NewNameAndNumberForm(tv[i]); err != nil {
				return
			}
			t.nANF = append(t.nANF, *nanf)
		}
	default:
		err = errorf("Unsupported %T input type %T\n", *o, x)
		return
	}

	if !t.Valid() {
		err = errorf("%T instance did not pass validity checks: %#v", t, *t)
		return
	}

	o = new(ObjectIdentifier)
	*o = *t

	return
}
