package validation

import "testing"

func TestAmount(t *testing.T) {
	cases := []struct {
		in      string
		want    float64
		wantErr bool
	}{
		{"1000", 1000, false},
		{"  2500.50 ", 2500.50, false},
		{"0", 0, true},
		{"-5", 0, true},
		{"", 0, true},
		{"abc", 0, true},
	}
	for _, c := range cases {
		got, err := Amount(c.in)
		if c.wantErr && err == nil {
			t.Errorf("Amount(%q): expected error, got nil", c.in)
			continue
		}
		if !c.wantErr && err != nil {
			t.Errorf("Amount(%q): unexpected error %v", c.in, err)
			continue
		}
		if !c.wantErr && got != c.want {
			t.Errorf("Amount(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestDate(t *testing.T) {
	if _, err := Date("2026-07-25"); err != nil {
		t.Errorf("valid date rejected: %v", err)
	}
	for _, bad := range []string{"", "2026-13-01", "25-07-2026", "not-a-date"} {
		if _, err := Date(bad); err == nil {
			t.Errorf("Date(%q): expected error, got nil", bad)
		}
	}
}

func TestRequired(t *testing.T) {
	if v, err := Required("  hello ", "name"); err != nil || v != "hello" {
		t.Errorf("Required trimmed value = %q err=%v", v, err)
	}
	if _, err := Required("   ", "name"); err == nil {
		t.Error("Required(blank): expected error")
	}
}

func TestUUIDField(t *testing.T) {
	if _, err := UUIDField("f47ac10b-58cc-4372-a567-0e02b2c3d479", "account"); err != nil {
		t.Errorf("valid uuid rejected: %v", err)
	}
	if _, err := UUIDField("nope", "account"); err == nil {
		t.Error("UUIDField(invalid): expected error")
	}
}

func TestIntField(t *testing.T) {
	if v, err := IntField(" 3 ", "type"); err != nil || v != 3 {
		t.Errorf("IntField = %d err=%v", v, err)
	}
	if _, err := IntField("x", "type"); err == nil {
		t.Error("IntField(invalid): expected error")
	}
}
