package api

type Repository interface {
    Lookup(recordType, name string) ([]string, error)
}
