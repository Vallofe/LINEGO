package Auth

import (
	"../LINE/thrift"
	"golang.org/x/xerrors"
)

type ThriftClient struct {
	StandardClient *thrift.TStandardClient
	httpClient     *thrift.THttpClient
}

func NewThriftClient(url string) (*ThriftClient, error) {
	transport, err := thrift.NewTHttpClient(url)
	if err != nil {
		return nil, xerrors.Errorf("failed to generate HTTP client: %w", err)
	}

	protocol := thrift.NewTCompactProtocolConf(transport, &thrift.TConfiguration{})

	standardClient := thrift.NewTStandardClient(protocol, protocol)

	client := &ThriftClient{
		StandardClient: standardClient,
		httpClient:     transport.(*thrift.THttpClient),
	}

	return client, nil
}

func (tc *ThriftClient) SetHeader(key, value string) {
	tc.httpClient.SetHeader(key, value)
}

func (tc *ThriftClient) DelHeader(key string) {
	tc.httpClient.DelHeader(key)
}
