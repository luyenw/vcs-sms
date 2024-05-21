package mock_entity

import (
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/stretchr/testify/mock"
)

type ESMock struct {
	mock.Mock
}

func (m *ESMock) Query(query string) (*esapi.Response, error) {
	args := m.Called(query)
	return args.Get(0).(*esapi.Response), args.Error(1)
}
