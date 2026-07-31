package server

import (
	"godrop/internal/common"
	"math"
)

const (
	defaultPort      = 9999
	defaultTimes     = math.MaxUint
	defaultQrOptions = false
)

type ServerOptions struct {
	Password  *string
	QrEnabled bool
	NTimes    *uint
	Port      *uint16
}

type serverConfig struct {
	password  string
	qrEnabled bool
	nTimes    int
	port      int
}

func (opts *ServerOptions) validate() serverConfig {
	var (
		password  = common.DefaultPassword
		port      = defaultPort
		times     = defaultTimes
		qrEnabled = defaultQrOptions
	)
	if opts.Password != nil {
		password = *opts.Password
	}
	if opts.NTimes != nil {
		times = opts.
	}

}
