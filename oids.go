package oid

import "sync"

type ObjectIdentifierMap map[string]*ObjectIdentifier

func (o ObjectIdentifierMap) Exists(term any) (exists bool) {
	_, exists = o.Get(term)
	return
}

func (o ObjectIdentifierMap) Set(key string, x *ObjectIdentifier) {
	mut := &sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()

	o[key] = x
}

func (o ObjectIdentifierMap) New(key, nanf string) (err error) {
	// create preliminary instance
	var x *ObjectIdentifier
	if x, err = NewObjectIdentifier(nanf); err != nil {
		return
	}

	mut := &sync.Mutex{}
	mut.Lock()
	defer mut.Unlock()

	o[key] = x
	return
}

func (o ObjectIdentifierMap) Get(term any) (*ObjectIdentifier, bool) {
	for k, v := range o {
		// lookup various forms of oid and asn1
		if v.Equal(term) {
			return v, !v.IsZero()
		}

		// try to match the term with the current
		// key iteration (if term is a string).
		if assert, ok := term.(string); ok {
			if eq(k, assert) {
				return v, !v.IsZero()
			}
		}
	}

	return nil, false
}
