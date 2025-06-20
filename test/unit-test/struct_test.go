package unittest

import (
	"reflect"
	"testing"
)

func TestRefMode(t *testing.T) {
	value := NewStructRef(false)
	typeName := reflect.TypeOf(value).String()
	if typeName != "*unittest.StructRef" {
		t.Errorf("expected *unittest.StructRef, but got %v", typeName)
	}
}

func TestValueMode(t *testing.T) {
	value := NewStructValue(false)
	typeName := reflect.TypeOf(value).String()
	if typeName != "unittest.StructValue" {
		t.Errorf("expected unittest.StructValue, but got %v", typeName)
	}
}

func TestInitMode(t *testing.T) {
	value := NewStructWithInit(false)
	if value.Debug != true {
		t.Errorf("NewStructWithInit should calling init method")
	}
}

func TestValueInitMode(t *testing.T) {
	value := NewStructValueWithInit(false)
	if value.Debug != true {
		t.Errorf("NewStructValueWithInit should calling init method")
	}
	typeName := reflect.TypeOf(value).String()
	if typeName != "unittest.StructValueWithInit" {
		t.Errorf("expected unittest.StructValueWithInit, but got %v", typeName)
	}
}

func TestSkipped(t *testing.T) {
	value := NewSkipeed("msg")
	if value.Status != 403 {
		t.Errorf("NewSkipeed should calling init method")
	}
}

func TestStructWithInitError(t *testing.T) {
	// 测试正常情况
	s, err := NewStructWithInitError(true, "test message")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if s == nil {
		t.Error("Expected non-nil struct")
	}
	if !s.Debug {
		t.Error("Expected Debug to be true")
	}
	if s.Msg != "test message" {
		t.Errorf("Expected Msg to be 'test message', got %s", s.Msg)
	}

	// 测试错误情况
	s2, err := NewStructWithInitError(true, "")
	if err == nil {
		t.Error("Expected error for empty message")
	}
	if s2 != nil {
		t.Error("Expected nil struct when error occurs")
	}
}

func TestStructValueWithInitError(t *testing.T) {
	// 测试正常情况
	s, err := NewStructValueWithInitError(true, "test message")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !s.Debug {
		t.Error("Expected Debug to be true")
	}
	if s.Msg != "test message" {
		t.Errorf("Expected Msg to be 'test message', got %s", s.Msg)
	}

	// 测试错误情况
	s2, err := NewStructValueWithInitError(true, "")
	if err == nil {
		t.Error("Expected error for empty message")
	}
	// 对于值类型，错误时应该返回零值
	if s2.Debug || s2.Msg != "" {
		t.Error("Expected zero value struct when error occurs")
	}
}
