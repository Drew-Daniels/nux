package config

import "testing"

func TestNormalizeSessionName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"my.project", "my_project"},
		{"my:project", "my_project"},
		{"my project", "my_project"},
		{"my..project", "my_project"},
		{"my. .project", "my_project"},
		{"-leading", "leading"},
		{"--leading", "leading"},
		{"trailing_", "trailing"},
		{"trailing__", "trailing"},
		{"conform.nvim", "conform_nvim"},
		{".hidden", "_hidden"},
		{"a:b.c d", "a_b_c_d"},
		{"---foo.bar..baz", "foo_bar_baz"},
		{"", ""},
		{"already_valid", "already_valid"},
		{"a", "a"},
		{"-", ""},
		{"_", ""},
	}
	for _, tt := range tests {
		got := NormalizeSessionName(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeSessionName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
