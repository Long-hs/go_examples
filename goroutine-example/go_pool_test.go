package goroutine_example

import "testing"

func Test_goroutineExample(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "testGoroutineExample"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goroutineExample()
		})
	}
}
