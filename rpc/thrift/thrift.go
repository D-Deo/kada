package thrift

import (
	"crypto/tls"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
)

func NewServer() *Server {
	return &Server{}
}

type Server struct {
}

func (o *Server) Start(port int, processor thrift.TProcessor, protocol string, buffered bool, framed bool, secure bool) error {
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary":
	default:
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	}

	var transportFactory thrift.TTransportFactory
	if buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	addr := fmt.Sprintf(":%d", port)

	var transport thrift.TServerTransport
	var err error
	if secure {
		cfg := new(tls.Config)
		cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
		if err != nil {
			return err
		}
		cfg.Certificates = append(cfg.Certificates, cert)
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}

	if err != nil {
		return err
	}

	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	return server.Serve()
}
