package domain

import (
	"testing"
)

func TestNewDNSName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "mydb", false},
		{"valid with hyphens", "my-database", false},
		{"valid with numbers", "db123", false},
		{"valid single char", "a", false},
		{"valid max length", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false},
		{"empty", "", true},
		{"too long", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", true},
		{"starts with hyphen", "-mydb", true},
		{"ends with hyphen", "mydb-", true},
		{"uppercase", "MyDB", true},
		{"underscore", "my_db", true},
		{"dot", "my.db", true},
		{"space", "my db", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewDNSName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDNSName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && v.String() != tt.input {
				t.Errorf("NewDNSName(%q).String() = %q, want %q", tt.input, v.String(), tt.input)
			}
		})
	}
}

func TestNewInteger(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid zero", "0", false},
		{"valid positive", "3", false},
		{"valid large", "100", false},
		{"empty", "", true},
		{"negative", "-1", true},
		{"float", "1.5", true},
		{"text", "abc", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewInteger(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInteger(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && v.String() != tt.input {
				t.Errorf("NewInteger(%q).String() = %q, want %q", tt.input, v.String(), tt.input)
			}
		})
	}
}

func TestNewNonEmptyString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "hello", false},
		{"valid with spaces", "hello world", false},
		{"empty", "", true},
		{"whitespace only", "   ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewNonEmptyString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNonEmptyString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && v.String() != tt.input {
				t.Errorf("NewNonEmptyString(%q).String() = %q, want %q", tt.input, v.String(), tt.input)
			}
		})
	}
}

func TestNewPort(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid min", "1", false},
		{"valid common", "8080", false},
		{"valid max", "65535", false},
		{"empty", "", true},
		{"zero", "0", true},
		{"too high", "65536", true},
		{"negative", "-1", true},
		{"text", "http", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewPort(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPort(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && v.String() != tt.input {
				t.Errorf("NewPort(%q).String() = %q, want %q", tt.input, v.String(), tt.input)
			}
		})
	}
}

func TestNewMemory(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid Mi", "256Mi", false},
		{"valid Gi", "4Gi", false},
		{"valid plain", "1024", false},
		{"valid decimal", "1.5Gi", false},
		{"empty", "", true},
		{"invalid suffix", "4GB", true},
		{"text", "lots", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewMemory(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMemory(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && v.String() != tt.input {
				t.Errorf("NewMemory(%q).String() = %q, want %q", tt.input, v.String(), tt.input)
			}
		})
	}
}

func TestNewCPU(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid millicore", "100m", false},
		{"valid whole", "2", false},
		{"valid decimal", "0.5", false},
		{"empty", "", true},
		{"text", "fast", true},
		{"invalid suffix", "2cores", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewCPU(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCPU(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && v.String() != tt.input {
				t.Errorf("NewCPU(%q).String() = %q, want %q", tt.input, v.String(), tt.input)
			}
		})
	}
}

func TestNewCronSchedule(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"every minute", "* * * * *", false},
		{"midnight daily", "0 0 * * *", false},
		{"weekdays at 5am", "0 5 * * 1-5", false},
		{"every 5 minutes", "*/5 * * * *", false},
		{"complex", "0,30 9-17 * * 1-5", false},
		{"specific day and time", "15 14 1 * *", false},
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"too few fields", "* * *", true},
		{"too many fields", "* * * * * *", true},
		{"invalid characters", "0 0 * * abc", true},
		{"not a cron", "not a cron", true},
		{"single field", "0", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewCronSchedule(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCronSchedule(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && v.String() != tt.input {
				t.Errorf("NewCronSchedule(%q).String() = %q, want %q", tt.input, v.String(), tt.input)
			}
		})
	}
}

func TestNewFilename(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid yaml", "manifest.yaml", false},
		{"valid with hyphens", "my-file.yml", false},
		{"valid no extension", "backup", false},
		{"valid with spaces", "file with spaces.yaml", false},
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"path traversal", "../etc/passwd", true},
		{"forward slash", "sub/file.yaml", true},
		{"backslash", `\file.yaml`, true},
		{"dotdot only", "..", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewFilename(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilename(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && v.String() != tt.input {
				t.Errorf("NewFilename(%q).String() = %q, want %q", tt.input, v.String(), tt.input)
			}
		})
	}
}

func TestNewPath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid relative", "templates", false},
		{"valid absolute", "/absolute/path", false},
		{"valid dot relative", "./relative", false},
		{"empty", "", true},
		{"whitespace only", "   ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPath(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && v.String() != tt.input {
				t.Errorf("NewPath(%q).String() = %q, want %q", tt.input, v.String(), tt.input)
			}
		})
	}
}

func TestParseValue(t *testing.T) {
	knownTypes := []struct {
		typeName   string
		validValue string
	}{
		{"dns-name", "mydb"},
		{"integer", "42"},
		{"string", "hello"},
		{"port", "8080"},
		{"memory", "256Mi"},
		{"cpu", "100m"},
		{"cron-schedule", "0 0 * * *"},
		{"filename", "manifest.yaml"},
		{"path", "templates"},
	}
	for _, tt := range knownTypes {
		t.Run(tt.typeName, func(t *testing.T) {
			v, err := ParseValue(tt.typeName, tt.validValue)
			if err != nil {
				t.Errorf("ParseValue(%q, %q) returned error: %v", tt.typeName, tt.validValue, err)
			}
			if v == nil {
				t.Errorf("ParseValue(%q, %q) returned nil", tt.typeName, tt.validValue)
			}
			if v != nil && v.String() != tt.validValue {
				t.Errorf("ParseValue(%q, %q).String() = %q, want %q", tt.typeName, tt.validValue, v.String(), tt.validValue)
			}
		})
	}

	t.Run("unknown type", func(t *testing.T) {
		_, err := ParseValue("unknown", "value")
		if err == nil {
			t.Error("ParseValue(\"unknown\", \"value\") expected error, got nil")
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		_, err := ParseValue("dns-name", "INVALID")
		if err == nil {
			t.Error("ParseValue(\"dns-name\", \"INVALID\") expected error, got nil")
		}
	})
}
