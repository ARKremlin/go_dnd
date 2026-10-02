package domain

import "testing"

func TestNormalizeTelegramUsername(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "petr", want: "petr"},
		{name: "with @", in: "@petr", want: "petr"},
		{name: "mixed case", in: "@Petr_Ivanov", want: "petr_ivanov"},
		{name: "spaces", in: "   @PETR  ", want: "petr"},
		{name: "empty", in: "", want: ""},
		{name: "only @", in: "@", want: ""},
		{name: "@ doubled", in: "@@petr", want: "@petr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeTelegramUsername(tt.in)
			if got != tt.want {
				t.Fatalf("NormalizeTelegramUsername(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
