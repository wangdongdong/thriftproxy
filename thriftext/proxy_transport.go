package thriftext

import (
	"encoding/binary"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// ProxyTransport append a Uint32 hash key to thrift.TFramedTransport to
// provide a general way to route thrift request
// format: ProxyTransportFrame ::= FrameTransportFrame HashKey
//         HashKey ::= Uint32
type TProxyTransport struct {
	*thrift.TFramedTransport
	hashKey uint32
	buf     []byte
}

func NewTProxyTransport(transport thrift.TTransport) *TProxyTransport {
	framedTransport, ok := transport.(*thrift.TFramedTransport)
	if !ok {
		panic("FramedTransport is required")
	}
	return &TProxyTransport{
		TFramedTransport: framedTransport,
		buf:              make([]byte, 4),
	}
}

func (p *TProxyTransport) SetHashKey(hashKey uint32) {
	p.hashKey = hashKey
}

func (p *TProxyTransport) HashKey() uint32 {
	return p.hashKey
}

func (p *TProxyTransport) Flush() error {
	binary.BigEndian.PutUint32(p.buf, p.hashKey)
	_, err := p.TFramedTransport.Write(p.buf)
	if err != nil {
		return thrift.NewTTransportExceptionFromError(err)
	}
	return p.TFramedTransport.Flush()
}

type tProxyTransportFactory struct {
	factory thrift.TTransportFactory
}

func NewTProxyTransportFactory(factory thrift.TTransportFactory) thrift.TTransportFactory {
	return &tProxyTransportFactory{
		factory: factory,
	}
}

func (p *tProxyTransportFactory) GetTransport(base thrift.TTransport) thrift.TTransport {
	transport := p.factory.GetTransport(base)
	if framedTransport, ok := transport.(*thrift.TFramedTransport); !ok {
		return nil
	} else {
		return NewTProxyTransport(framedTransport)
	}
}
