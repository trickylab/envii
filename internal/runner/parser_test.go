package runner

import (
	"strings"
	"testing"

	"github.com/trickylab/envii/internal/model"
)

func TestParseDotenvCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []*model.Var
	}{
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "comments and blanks only",
			input: "# comment\n\n  \n\t# indented comment\n",
		},
		{
			name:  "basic key values",
			input: "APP_NAME=envii\nPORT=3000\n",
			want: []*model.Var{
				{Key: "APP_NAME", Value: "envii"},
				{Key: "PORT", Value: "3000"},
			},
		},
		{
			name:  "valid key characters",
			input: "_FOO=bar\nFOO2=baz\nCAFÉ=latte\n",
			want: []*model.Var{
				{Key: "_FOO", Value: "bar"},
				{Key: "FOO2", Value: "baz"},
				{Key: "CAFÉ", Value: "latte"},
			},
		},
		{
			name:  "whitespace around key and value",
			input: "  KEY  =  value  \n",
			want:  []*model.Var{{Key: "KEY", Value: "value", Secret: true}},
		},
		{
			name: "comments in unquoted values",
			input: "URL=https://example.com#fragment\n" +
				"SPACED=hello world # comment\n" +
				"HASH=# comment\n",
			want: []*model.Var{
				{Key: "URL", Value: "https://example.com#fragment"},
				{Key: "SPACED", Value: "hello world"},
				{Key: "HASH", Value: ""},
			},
		},
		{
			name:  "empty values",
			input: "EMPTY=\nEMPTY_QUOTED=\"\"\nEMPTY_SINGLE=''\n",
			want: []*model.Var{
				{Key: "EMPTY", Value: ""},
				{Key: "EMPTY_QUOTED", Value: ""},
				{Key: "EMPTY_SINGLE", Value: ""},
			},
		},
		{
			name:  "double quoted escapes",
			input: `DOUBLE="hello \"world\"\nnext\rline\ttab\\slash\x"`,
			want:  []*model.Var{{Key: "DOUBLE", Value: "hello \"world\"\nnext\rline\ttab\\slashx"}},
		},
		{
			name:  "quoted values with trailing comments",
			input: "DOUBLE=\"value\" # comment\nSINGLE='value # still value' # comment\n",
			want: []*model.Var{
				{Key: "DOUBLE", Value: "value"},
				{Key: "SINGLE", Value: "value # still value"},
			},
		},
		{
			name:  "export prefix",
			input: "export API_TOKEN=secret\nexport  SECRET_KEY=value\n",
			want: []*model.Var{
				{Key: "API_TOKEN", Value: "secret", Secret: true},
				{Key: "SECRET_KEY", Value: "value", Secret: true},
			},
		},
		{
			name:  "equals in value",
			input: "DSN=postgres://user:pass@host/db?ssl=true\nCHAIN=a=b=c\n",
			want: []*model.Var{
				{Key: "DSN", Value: "postgres://user:pass@host/db?ssl=true"},
				{Key: "CHAIN", Value: "a=b=c"},
			},
		},
		{
			name:  "preserves duplicate keys in file order",
			input: "A=first\nB=second\nA=third\n",
			want: []*model.Var{
				{Key: "A", Value: "first"},
				{Key: "B", Value: "second"},
				{Key: "A", Value: "third"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDotenv(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			assertVars(t, got, tt.want)
		})
	}
}

func TestParseDotenvErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "missing equals", input: "NO_VALUE"},
		{name: "empty key", input: "=value"},
		{name: "leading digit", input: "1BAD=value"},
		{name: "dash in key", input: "BAD-KEY=value"},
		{name: "space in key", input: "BAD KEY=value"},
		{name: "dot in key", input: "BAD.KEY=value"},
		{name: "unterminated double quote", input: "KEY=\"unterminated"},
		{name: "only opening double quote", input: "KEY=\""},
		{name: "unterminated single quote", input: "KEY='unterminated"},
		{name: "double quoted trailing junk", input: "KEY=\"value\" extra"},
		{name: "single quoted trailing junk", input: "KEY='value' extra"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ParseDotenv(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestParseDotenvScannerError(t *testing.T) {
	_, err := ParseDotenv(strings.NewReader("KEY=" + strings.Repeat("x", 70*1024)))
	if err == nil {
		t.Fatal("expected scanner error")
	}
}

func assertVars(t *testing.T, got, want []*model.Var) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d vars, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Key != want[i].Key || got[i].Value != want[i].Value || got[i].Secret != want[i].Secret {
			t.Fatalf("var %d mismatch:\n got: %+v\nwant: %+v", i, got[i], want[i])
		}
	}
}
