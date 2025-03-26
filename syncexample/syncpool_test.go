package syncexample

import "testing"

func Test_syncPool(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "testSyncPoll"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syncPool()
		})
	}
}
