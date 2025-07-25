package repository

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/monitoror/monitoror/monitorables/rtsp/api"
	"github.com/monitoror/monitoror/monitorables/rtsp/config"
	pkgNet "github.com/monitoror/monitoror/pkg/net"
)

type rtspRepository struct {
	config *config.RTSP
	dialer pkgNet.Dialer
}

func NewRTSPRepository(conf *config.RTSP) api.Repository {
	timeout := time.Millisecond * time.Duration(conf.Timeout)
	return &rtspRepository{conf, &net.Dialer{Timeout: timeout}}
}

func parseDigest(header string) map[string]string {
	res := map[string]string{}
	header = strings.TrimSpace(header)
	if strings.HasPrefix(strings.ToLower(header), "digest") {
		header = strings.TrimSpace(header[len("digest"):])
	}
	parts := strings.Split(header, ",")
	for _, p := range parts {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 {
			res[strings.ToLower(strings.TrimSpace(kv[0]))] = strings.Trim(kv[1], "\"")
		}
	}
	return res
}

func computeDigest(username, realm, password, method, uri, nonce string) string {
	ha1 := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", username, realm, password)))
	ha2 := md5.Sum([]byte(fmt.Sprintf("%s:%s", method, uri)))
	resp := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", hex.EncodeToString(ha1[:]), nonce, hex.EncodeToString(ha2[:]))))
	return hex.EncodeToString(resp[:])
}

func (r *rtspRepository) Authenticate(hostname string, port int, path, method, username, password string) (bool, time.Duration, error) {
	start := time.Now()
	if path == "" {
		path = "/"
	}
	if method == "" {
		method = "OPTIONS"
	}
	target := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := r.dialer.Dial("tcp", target)
	if err != nil {
		return false, time.Since(start), err
	}
	if conn == nil {
		return false, time.Since(start), fmt.Errorf("no connection")
	}
	defer conn.Close()

	deadline := time.Now().Add(time.Millisecond * time.Duration(r.config.Timeout))
	_ = conn.SetDeadline(deadline)

	uri := fmt.Sprintf("rtsp://%s:%d%s", hostname, port, path)
	req1 := fmt.Sprintf("%s %s RTSP/1.0\r\nCSeq: 1\r\n\r\n", method, uri)
	if _, err = conn.Write([]byte(req1)); err != nil {
		return false, time.Since(start), err
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		return false, time.Since(start), err
	}
	resp := string(buf[:n])
	var authHeader string
	for _, line := range strings.Split(resp, "\r\n") {
		if strings.HasPrefix(strings.ToLower(line), "www-authenticate:") {
			authHeader = strings.TrimSpace(line[len("WWW-Authenticate:"):])
			break
		}
	}
	if authHeader == "" {
		return false, time.Since(start), fmt.Errorf("no auth header")
	}
	params := parseDigest(authHeader)
	realm := params["realm"]
	nonce := params["nonce"]
	response := computeDigest(username, realm, password, method, uri, nonce)

	req2 := fmt.Sprintf("%s %s RTSP/1.0\r\nCSeq: 2\r\nAuthorization: Digest username=\"%s\", realm=\"%s\", nonce=\"%s\", uri=\"%s\", response=\"%s\"\r\n\r\n",
		method, uri, username, realm, nonce, uri, response)
	if _, err = conn.Write([]byte(req2)); err != nil {
		return false, time.Since(start), err
	}

	n, err = conn.Read(buf)
	duration := time.Since(start)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		return false, duration, err
	}
	resp2 := string(buf[:n])
	if strings.HasPrefix(resp2, "RTSP/1.0 200") {
		return true, duration, nil
	}
	return false, duration, nil
}
