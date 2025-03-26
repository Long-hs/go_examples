package goroutine_example

import (
	"testing"
)

func TestRunExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "testPoolPackaged"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunExample()
		})
	}
}
