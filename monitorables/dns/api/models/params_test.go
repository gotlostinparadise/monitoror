package models

import "testing"

func TestDNSParams_Validate(t *testing.T) {
    p := &DNSParams{RecordType: "A", Name: "example.com"}
    if len(p.Validate()) == 0 {
        t.Fatal("expected error when expected value and pattern are empty")
    }
}
