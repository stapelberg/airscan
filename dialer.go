package airscan

import (
	"context"
	"log"
	"net"
	"sync"
)

type fallbackDialer struct {
	underlying *net.Dialer
	mu         sync.Mutex
	hostports  []string
}

func (d *fallbackDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	if debug {
		log.Printf("DialContext(%s, %s)", network, addr)
	}
	var lastErr error
	d.mu.Lock()
	defer d.mu.Unlock()
	for idx, hostport := range d.hostports {
		if debug {
			log.Printf("-> trying %s", hostport)
		}
		conn, err := d.underlying.DialContext(ctx, network, hostport)
		if err != nil {
			lastErr = err
			continue
		}
		d.hostports = append(append([]string{hostport}, d.hostports[:idx]...), d.hostports[idx+1:]...)
		return conn, nil
	}
	return nil, lastErr
}
