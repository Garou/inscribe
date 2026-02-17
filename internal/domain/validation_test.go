package domain

import (
	"testing"
)

func TestValidateDNSName(t *testing.T) {
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
			err := ValidateDNSName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDNSName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateInteger(t *testing.T) {
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
			err := ValidateInteger(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInteger(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateString(t *testing.T) {
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
			err := ValidateString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
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
			err := ValidatePort(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateMemory(t *testing.T) {
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
			err := ValidateMemory(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMemory(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateCPU(t *testing.T) {
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
			err := ValidateCPU(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCPU(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestGetValidator(t *testing.T) {
	knownTypes := []string{"dns-name", "integer", "string", "port", "memory", "cpu"}
	for _, vt := range knownTypes {
		t.Run(vt, func(t *testing.T) {
			v, err := GetValidator(vt)
			if err != nil {
				t.Errorf("GetValidator(%q) returned error: %v", vt, err)
			}
			if v == nil {
				t.Errorf("GetValidator(%q) returned nil validator", vt)
			}
		})
	}

	t.Run("unknown type", func(t *testing.T) {
		_, err := GetValidator("unknown")
		if err == nil {
			t.Error("GetValidator(\"unknown\") expected error, got nil")
		}
	})
}
