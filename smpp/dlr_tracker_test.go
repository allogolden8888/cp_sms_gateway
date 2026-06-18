package smpp

import (
	"testing"
)

func TestDLRTracker(t *testing.T) {
	dlr := DLR{
		MessageID: "abc",
	}
	id1 := "abc"
	id2 := "abc1"

	tests := []struct {
		name      string
		input     *DLR
		want      string
		id        string
		wantEmpty bool
	}{
		{
			name:  "expected id",
			input: &dlr,
			want:  "abc",
			id:    id1,
		},
		{
			name:      "unexpected id",
			input:     &dlr,
			want:      "",
			id:        id2,
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dlrTracker := NewDLRTracker()

			dlrTracker.Expect(tt.id)
			dlrTracker.Receive(tt.input)
			if tt.wantEmpty {
				select {
				case id := <-dlrTracker.Done():
					t.Errorf("not expected ID, got %v", id)
				default:
					// works correctly, test passed
				}
			} else {
				id := <-dlrTracker.Done()
				if id != tt.want {
					t.Errorf("Unexpected ID, expected %v, got %v", tt.want, id)
				}
			}

		})
	}
}
