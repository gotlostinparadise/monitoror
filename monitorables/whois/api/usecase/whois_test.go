package usecase

import (
	"testing"
	"time"

	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/monitorables/whois/api/models"
	"github.com/stretchr/testify/assert"
)

type fakeRepo struct {
	expiry time.Time
	raw    string
	err    error
}

func (f *fakeRepo) DomainExpiration(domain string) (time.Time, string, error) {
	return f.expiry, f.raw, f.err
}

func TestWHOIS_DisplayDomain(t *testing.T) {
	expiry := time.Now().Add(24 * time.Hour)
	repo := &fakeRepo{expiry: expiry, raw: "Domain Name: example.com"}
	usecase := NewWHOISUsecase(repo)

	params := &models.WHOISParams{Domain: "example.com", Display: "domain"}

	tile, err := usecase.WHOIS(params)
	if assert.NoError(t, err) {
		assert.Equal(t, "example.com", tile.Message)
		assert.Equal(t, coreModels.SuccessStatus, tile.Status)
	}
}

func TestWHOIS_DisplayRegexNotFound(t *testing.T) {
	expiry := time.Now().Add(24 * time.Hour)
	repo := &fakeRepo{expiry: expiry, raw: "Domain Name: example.com"}
	usecase := NewWHOISUsecase(repo)

	params := &models.WHOISParams{Domain: "example.com", Display: "Registrar:(.*)"}

	tile, err := usecase.WHOIS(params)
	if assert.NoError(t, err) {
		assert.Equal(t, "поле не найдено", tile.Message)
	}
}
