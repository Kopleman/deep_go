package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type isCallable = func() interface{}

// TODO В ДЗ не указанно должны ли быть два разных регистра зависимостей для синглтон/не синглтон зависимостей, сделал общий
type Container struct {
	deps map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		deps: make(map[string]interface{}),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.deps[name] = constructor
}

func (c *Container) RegisterSingletonType(name string, constructor interface{}) {
	callable, callableOk := constructor.(isCallable)
	if !callableOk {
		c.deps[name] = constructor
		return
	}
	c.deps[name] = callable()
}

func (c *Container) Resolve(name string) (interface{}, error) {
	constructor, ok := c.deps[name]
	if !ok {
		return nil, errors.New("no dependency for name: " + name)
	}
	callable, callableOk := constructor.(isCallable)
	if !callableOk {
		return constructor, nil
	}
	return callable(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}

func TestRegisterSingletonType(t *testing.T) {
	container := NewContainer()

	container.RegisterSingletonType("UserService", func() interface{} {
		return &UserService{NotEmptyStruct: true}
	})

	container.RegisterSingletonType("Config", "test-config")

	assert.NotNil(t, container.deps["UserService"])
	assert.Equal(t, "test-config", container.deps["Config"])
}

func TestResolveSingleton(t *testing.T) {
	container := NewContainer()

	container.RegisterSingletonType("UserService", func() interface{} {
		return &UserService{NotEmptyStruct: true}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	assert.NotNil(t, userService1)

	u1, ok := userService1.(*UserService)
	assert.True(t, ok)
	assert.True(t, u1.NotEmptyStruct)
}

func TestSingletonSameInstance(t *testing.T) {
	container := NewContainer()

	container.RegisterSingletonType("UserService", func() interface{} {
		return &UserService{NotEmptyStruct: true}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)

	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.True(t, u1 == u2, "Singleton should return the same instance")
}

func TestRegisterTypeWithNonCallable(t *testing.T) {
	container := NewContainer()

	container.RegisterType("Config", "test-config")

	config, err := container.Resolve("Config")
	assert.NoError(t, err)
	assert.Equal(t, "test-config", config)
}

func TestMixedTypes(t *testing.T) {
	container := NewContainer()

	container.RegisterType("UserService", func() interface{} {
		return &UserService{NotEmptyStruct: true}
	})

	container.RegisterSingletonType("MessageService", func() interface{} {
		return &MessageService{NotEmptyStruct: true}
	})

	user1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	user2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := user1.(*UserService)
	u2 := user2.(*UserService)
	assert.False(t, u1 == u2, "Regular types should return different instances")

	msg1, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	msg2, err := container.Resolve("MessageService")
	assert.NoError(t, err)

	m1 := msg1.(*MessageService)
	m2 := msg2.(*MessageService)
	assert.True(t, m1 == m2, "Singletons should return the same instance")
}

func TestRegisterSingletonTypeWithNonCallable(t *testing.T) {
	container := NewContainer()

	container.RegisterSingletonType("Config", "test-config")

	config, err := container.Resolve("Config")
	assert.NoError(t, err)
	assert.Equal(t, "test-config", config)

	config2, err := container.Resolve("Config")
	assert.NoError(t, err)
	assert.Equal(t, "test-config", config2)
	assert.True(t, config == config2)
}

func TestResolvePriority(t *testing.T) {
	container := NewContainer()

	container.RegisterType("Service", func() interface{} {
		return &UserService{NotEmptyStruct: false}
	})
	container.RegisterSingletonType("Service", func() interface{} {
		return &MessageService{NotEmptyStruct: true}
	})

	service, err := container.Resolve("Service")
	assert.NoError(t, err)

	messageService, ok := service.(*MessageService)
	assert.True(t, ok)
	assert.True(t, messageService.NotEmptyStruct)
}

func TestNilConstructor(t *testing.T) {
	container := NewContainer()

	container.RegisterType("NilService", nil)

	service, err := container.Resolve("NilService")
	assert.NoError(t, err)
	assert.Nil(t, service)
}

func TestEmptyStringName(t *testing.T) {
	container := NewContainer()

	container.RegisterType("", func() interface{} {
		return &UserService{}
	})

	service, err := container.Resolve("")
	assert.NoError(t, err)
	assert.NotNil(t, service)

	_, ok := service.(*UserService)
	assert.True(t, ok)
}

func TestMultipleRegistrations(t *testing.T) {
	container := NewContainer()

	container.RegisterType("Service", func() interface{} {
		return &UserService{NotEmptyStruct: false}
	})

	container.RegisterType("Service", func() interface{} {
		return &MessageService{NotEmptyStruct: true}
	})

	service, err := container.Resolve("Service")
	assert.NoError(t, err)

	messageService, ok := service.(*MessageService)
	assert.True(t, ok)
	assert.True(t, messageService.NotEmptyStruct)
}

func TestComplexObjectCreation(t *testing.T) {
	container := NewContainer()

	container.RegisterSingletonType("ComplexService", func() interface{} {
		return map[string]interface{}{
			"userService":    &UserService{NotEmptyStruct: true},
			"messageService": &MessageService{NotEmptyStruct: true},
			"config":         "test-config",
		}
	})

	complexService, err := container.Resolve("ComplexService")
	assert.NoError(t, err)
	assert.NotNil(t, complexService)

	serviceMap, ok := complexService.(map[string]interface{})
	assert.True(t, ok)

	userService, ok := serviceMap["userService"].(*UserService)
	assert.True(t, ok)
	assert.True(t, userService.NotEmptyStruct)

	messageService, ok := serviceMap["messageService"].(*MessageService)
	assert.True(t, ok)
	assert.True(t, messageService.NotEmptyStruct)

	assert.Equal(t, "test-config", serviceMap["config"])
}

func TestErrorHandling(t *testing.T) {
	container := NewContainer()

	_, err := container.Resolve("NonExistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no dependency for name: NonExistent")

	_, err = container.Resolve("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no dependency for name: ")
}

func TestMemoryEfficiency(t *testing.T) {
	container := NewContainer()

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Service%d", i)
		container.RegisterType(name, func() interface{} {
			return &UserService{NotEmptyStruct: true}
		})
	}

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Service%d", i)
		service, err := container.Resolve(name)
		assert.NoError(t, err)
		assert.NotNil(t, service)

		_, ok := service.(*UserService)
		assert.True(t, ok)
	}
}
