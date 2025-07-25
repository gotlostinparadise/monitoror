package repository

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	rtspConfig "github.com/monitoror/monitoror/monitorables/rtsp/config"
	pkgNet "github.com/monitoror/monitoror/pkg/net"
	"github.com/monitoror/monitoror/pkg/net/mocks"

	"github.com/stretchr/testify/assert"
	. "github.com/stretchr/testify/mock"
)

func initRepository(t *testing.T, dialer pkgNet.Dialer) *rtspRepository {
	conf := &rtspConfig.RTSP{Timeout: 1000}
	repo := NewRTSPRepository(conf)
	r, ok := repo.(*rtspRepository)
	if assert.True(t, ok) {
		r.dialer = dialer
		return r
	}
	return nil
}

func md5Hex(b []byte) string { return hex.EncodeToString(md5.Sum(b)[:]) }

func expectedDigest(user, realm, pass, method, uri, nonce string) string {
	ha1 := md5Hex([]byte(fmt.Sprintf("%s:%s:%s", user, realm, pass)))
	ha2 := md5Hex([]byte(fmt.Sprintf("%s:%s", method, uri)))
	final := md5Hex([]byte(fmt.Sprintf("%s:%s:%s", ha1, nonce, ha2)))
	return final
}

func TestRepository_Authenticate_Digest(t *testing.T) {
	challenge := "RTSP/1.0 401 Unauthorized\r\nCSeq: 1\r\nWWW-Authenticate: Digest realm=\"realm\", nonce=\"nonce\"\r\n\r\n"
	success := "RTSP/1.0 200 OK\r\nCSeq: 2\r\n\r\n"

	mockConn := new(mocks.Conn)
	mockConn.On("Close").Return(nil)
	mockConn.On("SetDeadline", AnythingOfType("time.Time")).Return(nil)
	mockConn.On("Write", Anything).Return(func(b []byte) (int, error) { return len(b), nil }).Once()
	mockConn.On("Read", Anything).Run(func(args Arguments) { copy(args.Get(0).([]byte), []byte(challenge)) }).Return(len(challenge), nil).Once()

	var second []byte
	mockConn.On("Write", Anything).Return(func(b []byte) (int, error) { second = append([]byte(nil), b...); return len(b), nil }).Once()
	mockConn.On("Read", Anything).Run(func(args Arguments) { copy(args.Get(0).([]byte), []byte(success)) }).Return(len(success), nil).Once()

	mockDialer := new(mocks.Dialer)
	mockDialer.On("Dial", AnythingOfType("string"), AnythingOfType("string")).Return(mockConn, nil)

	repo := initRepository(t, mockDialer)
	if repo != nil {
		ok, dur, err := repo.Authenticate("test", 554, "/", "OPTIONS", "user", "pass")
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NotZero(t, dur)

		digest := expectedDigest("user", "realm", "pass", "OPTIONS", "rtsp://test:554/", "nonce")
		assert.Contains(t, string(second), fmt.Sprintf("response=\"%s\"", digest))
	}
}
