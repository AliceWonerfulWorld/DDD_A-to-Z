package brainfuck

import (
	"errors"
	"strings"
	"testing"
)

func TestExecuteBF(t *testing.T) {
	t.Run("input output", func(t *testing.T) {
		got, err := ExecuteBF(",.", "A")
		if err != nil {
			t.Fatalf("ExecuteBF() error = %v", err)
		}
		if got != "A" {
			t.Fatalf("ExecuteBF() = %q, want %q", got, "A")
		}
	})

	t.Run("loop and ignored characters", func(t *testing.T) {
		got, err := ExecuteBF("+++[>++<-]>.", "")
		if err != nil {
			t.Fatalf("ExecuteBF() error = %v", err)
		}
		if got != string([]byte{6}) {
			t.Fatalf("ExecuteBF() bytes = %v, want [6]", []byte(got))
		}
	})

	t.Run("nested loop", func(t *testing.T) {
		got, err := ExecuteBF("++[>++[>++<-]<-]>>.", "")
		if err != nil {
			t.Fatalf("ExecuteBF() error = %v", err)
		}
		if got != string([]byte{8}) {
			t.Fatalf("ExecuteBF() bytes = %v, want [8]", []byte(got))
		}
	})

	t.Run("comma without input writes zero", func(t *testing.T) {
		got, err := ExecuteBF(",.", "")
		if err != nil {
			t.Fatalf("ExecuteBF() error = %v", err)
		}
		if got != string([]byte{0}) {
			t.Fatalf("ExecuteBF() bytes = %v, want [0]", []byte(got))
		}
	})
}

func TestExecuteBFErrors(t *testing.T) {
	tests := []struct {
		name string
		code string
		err  error
	}{
		{name: "unmatched open bracket", code: "[", err: ErrBracketMismatch},
		{name: "unmatched close bracket", code: "]", err: ErrBracketMismatch},
		{name: "pointer underflow", code: "<", err: ErrPointerBounds},
		{name: "pointer overflow", code: strings.Repeat(">", tapeSize), err: ErrPointerBounds},
		{name: "step limit", code: "+[]", err: ErrStepLimit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ExecuteBF(tt.code, "")
			if !errors.Is(err, tt.err) {
				t.Fatalf("ExecuteBF() error = %v, want %v", err, tt.err)
			}
		})
	}
}
