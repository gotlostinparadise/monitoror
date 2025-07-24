package repository

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/monitoror/monitoror/monitorables/port/api"
	"github.com/monitoror/monitoror/monitorables/port/config"
	pkgNet "github.com/monitoror/monitoror/pkg/net"
)

type (
	portRepository struct {
		config *config.Port
		dialer pkgNet.Dialer
	}
)

func NewPortRepository(conf *config.Port) api.Repository {
	timeout := time.Millisecond * time.Duration(conf.Timeout)
	return &portRepository{conf, &net.Dialer{Timeout: timeout}}
}

func (r *portRepository) OpenSocket(hostname string, port int, network string, payload []byte) (responding bool, duration time.Duration, err error) {
	start := time.Now()
	target := fmt.Sprintf("%s:%d", hostname, port)

	conn, err := r.dialer.Dial(network, target)
	if err != nil {
		return false, time.Since(start), err
	}
	if conn == nil {
		return false, time.Since(start), fmt.Errorf("no connection")
	}
	defer conn.Close()

	deadline := time.Now().Add(time.Millisecond * time.Duration(r.config.Timeout))
	_ = conn.SetDeadline(deadline)

	if len(payload) > 0 {
		if _, err = conn.Write(payload); err != nil {
			return false, time.Since(start), err
		}
	} else if network == "udp" {
		// send empty datagram to validate connection
		if _, err = conn.Write([]byte{}); err != nil {
			return false, time.Since(start), err
		}
	}

	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	duration = time.Since(start)
	if err == nil || err == io.EOF {
		responding = true
		err = nil
	} else if ne, ok := err.(net.Error); ok && ne.Timeout() {
		err = nil
	}

	return
}
