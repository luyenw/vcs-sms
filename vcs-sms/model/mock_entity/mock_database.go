package mock_entity

import (
	"fmt"
	"vcs-sms/model/entity"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDatabase struct {
	mock.Mock
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{}
}

func (m *MockDatabase) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	value := args.Get(1)
	switch v := out.(type) {
	case *[]entity.Server:
		*v = *value.(*[]entity.Server)
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Create(value interface{}) *gorm.DB {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Save(value interface{}) *gorm.DB {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	if err := args.Error(1); err != nil {
		return &gorm.DB{Error: err}
	} else {
		if len(args) > 1 {
			value := args.Get(0)
			if value != nil {
				switch v := dest.(type) {
				case *entity.Server:
					*v = *value.(*entity.Server)
				}
			}
		}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Order(value interface{}) *gorm.DB {
	args := m.Called(value)
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Limit(limit int) *gorm.DB {
	args := m.Called(limit)
	fmt.Println(limit)
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Offset(offset int) *gorm.DB {
	args := m.Called(offset)
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Table(name string, o ...interface{}) *gorm.DB {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Update(column string, value interface{}) *gorm.DB {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Preload(column string, conditions ...interface{}) *gorm.DB {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Where(query interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(query, conds)
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Model(value interface{}) *gorm.DB {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}

func (m *MockDatabase) Association(column string) *gorm.Association {
	args := m.Called()
	if err := args.Error(0); err != nil {
		return &gorm.Association{Error: err}
	}
	return &gorm.Association{}
}
