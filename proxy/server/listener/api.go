package listener

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/jwhittle933/gateway/logger"
)

type GracefulListener struct {
	ln          net.Listener
	maxWaitTime time.Duration
	done        chan struct{}
	connCount   uint64
	shutdown    uint64
	logger      logger.Logger
}

type gracefulConn struct {
	net.Conn
	ln *GracefulListener
}

func New(ln net.Listener, maxWait time.Duration, logger logger.Logger) *GracefulListener {
	return &GracefulListener{
		ln:          ln,
		maxWaitTime: maxWait,
		done:        make(chan struct{}),
		logger:      logger,
	}
}

func (gl *GracefulListener) Accept() (net.Conn, error) {
	c, err := gl.ln.Accept()
	if err != nil {
		return nil, err
	}

	atomic.AddUint64(&gl.connCount, 1)
	return &gracefulConn{
		Conn: c,
		ln:   gl,
	}, nil
}

func (gl *GracefulListener) Addr() net.Addr {
	return gl.ln.Addr()
}

func (gl *GracefulListener) Close() error {
	if err := gl.ln.Close(); err != nil {
		return err
	}

	return gl.waitForZeroConns()
}

func (gl *GracefulListener) waitForZeroConns() error {
	atomic.AddUint64(&gl.shutdown, 1)

	if atomic.LoadUint64(&gl.connCount) == 0 {
		close(gl.done)
		return nil
	}

	select {
	case <-gl.done:
		break
	case <-time.After(gl.maxWaitTime):
		return fmt.Errorf("cannot complete graceful listener in %s", gl.maxWaitTime)
	}

	return nil
}

func (gl *GracefulListener) closeConn() {
	connCount := atomic.AddUint64(&gl.connCount, ^uint64(0))

	if atomic.LoadUint64(&gl.shutdown) != 0 && connCount == 0 {
		close(gl.done)
	}
}

func (gl *GracefulListener) AwaitShutdown(server *fasthttp.Server) {
	listenErr := make(chan error, 1)
	go func() {
		listenErr <- server.Serve(gl)
	}()

	gl.logger.Server("Listening at %s", gl.Addr())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-listenErr:
			if err != nil {
				gl.logger.Kill("%s\n", err.Error())
			}
		case exit := <-c:
			gl.logger.Server("Received shutdown %+v\n", exit)
			server.DisableKeepalive = true

			if err := gl.Close(); err != nil {
				gl.logger.Kill("Error closing listener: %s\n", err.Error())
			}

			gl.logger.Server("Closed")
		}
	}
}

func (gc *gracefulConn) Close() error {
	err := gc.Conn.Close()

	if err != nil {
		return err
	}

	gc.ln.closeConn()
	return nil
}
