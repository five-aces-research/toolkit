package test_implementation

import (
	"errors"
	"fmt"
	"github.com/five-aces-research/toolkit/live/cle/cparser"
	"io"
)

type MockCommunicator struct {
	input        []byte
	output       []byte
	variables    map[string]cparser.Variable
	errorMessage error
}

func NewMockCommunicator() *MockCommunicator {
	return &MockCommunicator{
		variables: make(map[string]cparser.Variable),
	}
}

func (mc *MockCommunicator) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))
	return len(p), nil
}

func (mc *MockCommunicator) Read(p []byte) (n int, err error) {
	if len(mc.output) == 0 {
		return 0, io.EOF
	}
	n = copy(p, mc.output)
	mc.output = mc.output[n:]
	return n, nil
}

func (mc *MockCommunicator) AddVariable(name string, variable cparser.Variable) {
	if mc.variables == nil {
		mc.variables = make(map[string]cparser.Variable)
	}
	mc.variables[name] = variable
}

func (mc *MockCommunicator) GetVariable(name string) (*cparser.Variable, error) {
	if mc.variables == nil {
		return nil, errors.New("no variables set")
	}
	variable, ok := mc.variables[name]
	if !ok {
		return nil, fmt.Errorf("variable not found: %s", name)
	}
	return &variable, nil
}

func (mc *MockCommunicator) ErrorMessage(err error) {
	fmt.Println(err)
}
