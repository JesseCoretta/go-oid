# go-oid

[![Go Report Card](https://goreportcard.com/badge/JesseCoretta/go-oid)](https://goreportcard.com/report/github.com/JesseCoretta/go-oid) [![GoDoc](https://godoc.org/github.com/JesseCoretta/go-oid?status.svg)](https://godoc.org/github.com/JesseCoretta/go-oid)

Package oid offers a convenient ASN.1 Object Identifier type and associated methods.

ASN.1 Object Identifiers encompass information that goes beyond their dotted representation. This tiny package merely facilitates the handling of ASN.1 NameAndNumberForm values and alternate names that may be associated with a given OID in the wild.

## Example

```
package main
import (
        "log"

        "github.com/JesseCoretta/go-oid"
)

func main() {
        chkerr := func(err error) {
                if err != nil {
                        log.Fatal(err)
                }
        }

        // Test 1
        internet, err := oid.NewObjectIdentifier(`{ iso(1) identified-organization(3) dod(6) internet(1) }`)
        chkerr(err)

        log.SetFlags(0)
        log.Printf("dotNotation: %s\n", internet.ASN1().String())
        log.Printf("ASN.1 Notation: %s\n", internet.String())
        log.Printf("NameAndNumberForm: %s\n", internet.NameAndNumberForm())

        // Test 2
        var sa *oid.ObjectIdentifier
        sa, err = oid.NewObjectIdentifier(`1 3 6 1 5 5 7 3 1`) // legal but less than helpful
        chkerr(err)
        sa.SetAltNames(`serverAuth`)

        log.Printf("dotNotation: %s\n", sa.ASN1().String())
        log.Printf("Alt. Names: %v\n", sa.AltNames())

}
```

## Result
```
dotNotation: 1.3.6.1
ASN.1 Notation: { iso(1) identified-organization(3) dod(6) internet(1) }
NameAndNumberForm: [internet 1]
dotNotation: 1.3.6.1.5.5.7.3.1
Alt. Names: [serverAuth]
```
