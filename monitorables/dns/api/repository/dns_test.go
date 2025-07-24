package repository

import "testing"

func TestRepository_Lookup_Invalid(t *testing.T) {
    repo := NewDNSRepository()
    _, err := repo.Lookup("ZZ", "example.com")
    if err == nil {
        t.Fail()
    }
}
