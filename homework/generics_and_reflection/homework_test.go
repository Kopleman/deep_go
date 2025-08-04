package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

// Дополнительные структуры для тестирования обобщенной функции
type Config struct {
	Host     string `properties:"host"`
	Port     int    `properties:"port"`
	Debug    bool   `properties:"debug"`
	LogLevel string `properties:"log_level,omitempty"`
}

type Product struct {
	ID          int     `properties:"id"`
	Name        string  `properties:"name"`
	Price       float64 `properties:"price"`
	InStock     bool    `properties:"in_stock"`
	Description string  `properties:"description,omitempty"`
}

type User struct {
	Username string `properties:"username"`
	Email    string `properties:"email,omitempty"`
	Active   bool   `properties:"active"`
	Role     string `properties:"role"`
}

// Serialize - обобщенная функция для сериализации любой структуры в .properties формат
func Serialize[T any](obj T) string {
	var result []string
	v := reflect.ValueOf(obj)
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return ""
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		tag := fieldType.Tag.Get("properties")
		if tag == "" || tag == "-" {
			continue
		}

		parts := strings.Split(tag, ",")
		key := parts[0]

		hasOmitEmpty := false
		for _, part := range parts[1:] {
			if strings.TrimSpace(part) == "omitempty" {
				hasOmitEmpty = true
				break
			}
		}

		if hasOmitEmpty && isEmptyValue(field) {
			continue
		}

		value := serializeValue(field)
		result = append(result, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(result, "\n")
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	default:
		return v.IsZero()
	}
}

func serializeValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
		"test case with all fields filled": {
			person: Person{
				Name:    "Alice Smith",
				Address: "New York",
				Age:     25,
				Married: false,
			},
			result: "name=Alice Smith\naddress=New York\nage=25\nmarried=false",
		},
		"test case with zero values and omitempty": {
			person: Person{
				Name:    "Bob Johnson",
				Address: "", // пустая строка с omitempty
				Age:     0,  // нулевое значение
				Married: false,
			},
			result: "name=Bob Johnson\nage=0\nmarried=false",
		},
		"test case with special characters in name": {
			person: Person{
				Name:    "John & Jane",
				Age:     35,
				Married: true,
			},
			result: "name=John & Jane\nage=35\nmarried=true",
		},
		"test case with empty name": {
			person: Person{
				Name:    "",
				Age:     40,
				Married: true,
			},
			result: "name=\nage=40\nmarried=true",
		},
		"test case with negative age": {
			person: Person{
				Name:    "Young Person",
				Age:     -5,
				Married: false,
			},
			result: "name=Young Person\nage=-5\nmarried=false",
		},
		"test case with large age": {
			person: Person{
				Name:    "Old Person",
				Age:     150,
				Married: true,
			},
			result: "name=Old Person\nage=150\nmarried=true",
		},
		"test case with only name field": {
			person: Person{
				Name: "Only Name",
			},
			result: "name=Only Name\nage=0\nmarried=false",
		},
		"test case with only address field (should be omitted)": {
			person: Person{
				Address: "Some Address",
			},
			result: "name=\naddress=Some Address\nage=0\nmarried=false",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestOmitEmptyBehavior(t *testing.T) {
	tests := []struct {
		name     string
		person   Person
		expected string
	}{
		{
			name: "empty string with omitempty",
			person: Person{
				Name:    "Test",
				Address: "",
				Age:     25,
				Married: false,
			},
			expected: "name=Test\nage=25\nmarried=false",
		},
		{
			name: "non-empty string with omitempty",
			person: Person{
				Name:    "Test",
				Address: "Valid Address",
				Age:     25,
				Married: false,
			},
			expected: "name=Test\naddress=Valid Address\nage=25\nmarried=false",
		},
		{
			name: "whitespace string with omitempty",
			person: Person{
				Name:    "Test",
				Address: "   ",
				Age:     25,
				Married: false,
			},
			expected: "name=Test\naddress=   \nage=25\nmarried=false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Serialize(tt.person)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataTypes(t *testing.T) {
	person := Person{
		Name:    "Type Test",
		Address: "Address Test",
		Age:     42,
		Married: true,
	}

	result := Serialize(person)

	assert.Contains(t, result, "name=Type Test")
	assert.Contains(t, result, "address=Address Test")
	assert.Contains(t, result, "age=42")
	assert.Contains(t, result, "married=true")
}

func TestFieldOrder(t *testing.T) {
	person := Person{
		Name:    "Order Test",
		Address: "Order Address",
		Age:     30,
		Married: true,
	}

	result := Serialize(person)
	lines := strings.Split(result, "\n")

	assert.Equal(t, "name=Order Test", lines[0])
	assert.Equal(t, "address=Order Address", lines[1])
	assert.Equal(t, "age=30", lines[2])
	assert.Equal(t, "married=true", lines[3])
}

func TestGenericSerializeWithConfig(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     8080,
		Debug:    true,
		LogLevel: "info",
	}

	result := Serialize(config)
	expected := "host=localhost\nport=8080\ndebug=true\nlog_level=info"
	assert.Equal(t, expected, result)
}

func TestGenericSerializeWithConfigOmitEmpty(t *testing.T) {
	config := Config{
		Host:  "localhost",
		Port:  8080,
		Debug: false,
	}

	result := Serialize(config)
	expected := "host=localhost\nport=8080\ndebug=false"
	assert.Equal(t, expected, result)
}

func TestGenericSerializeWithProduct(t *testing.T) {
	product := Product{
		ID:          1,
		Name:        "Laptop",
		Price:       999.99,
		InStock:     true,
		Description: "High-performance laptop",
	}

	result := Serialize(product)
	expected := "id=1\nname=Laptop\nprice=999.99\nin_stock=true\ndescription=High-performance laptop"
	assert.Equal(t, expected, result)
}

func TestGenericSerializeWithProductOmitEmpty(t *testing.T) {
	product := Product{
		ID:      2,
		Name:    "Mouse",
		Price:   29.99,
		InStock: false,
	}

	result := Serialize(product)
	expected := "id=2\nname=Mouse\nprice=29.99\nin_stock=false"
	assert.Equal(t, expected, result)
}

func TestGenericSerializeWithUser(t *testing.T) {
	user := User{
		Username: "john_doe",
		Email:    "john@example.com",
		Active:   true,
		Role:     "admin",
	}

	result := Serialize(user)
	expected := "username=john_doe\nemail=john@example.com\nactive=true\nrole=admin"
	assert.Equal(t, expected, result)
}

func TestGenericSerializeWithUserOmitEmpty(t *testing.T) {
	user := User{
		Username: "guest",
		Active:   false,
		Role:     "user",
	}

	result := Serialize(user)
	expected := "username=guest\nactive=false\nrole=user"
	assert.Equal(t, expected, result)
}

func TestGenericSerializeWithNonStruct(t *testing.T) {
	assert.Equal(t, "", Serialize("string"))
	assert.Equal(t, "", Serialize(42))
	assert.Equal(t, "", Serialize(true))

	person := &Person{Name: "Test"}
	assert.Equal(t, "", Serialize(person))
}
