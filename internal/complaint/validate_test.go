package complaint

import (
	"errors"
	"testing"
)

func TestInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		in      Input
		wantErr bool
	}{
		{
			name: "valid other complaint",
			in:   Input{Category: CategoryOther, Message: "Преподаватель угрожает отчислением"},
		},
		{
			name: "valid organization complaint",
			in:   Input{Category: CategoryOrg, Organization: "Student Union", Message: "Нарушен устав организации"},
		},
		{
			name: "valid staff complaint with subject",
			in: Input{
				Category: CategoryStaff, School: "ШИТиИ", TeacherName: "Иванов И.И.",
				Subject: "Алгоритмы", Message: "Занижает оценки без объяснения",
			},
		},
		{
			name: "valid staff complaint without subject",
			in: Input{
				Category: CategoryStaff, School: "БШ", TeacherName: "Иванов И.И.",
				Message: "Занижает оценки без объяснения",
			},
		},
		{
			name: "valid cafeteria complaint",
			in:   Input{Category: CategoryCafeteria, Message: "Отказали в доступе без причины"},
		},
		{
			name: "message of exactly 10 Cyrillic runes (20 bytes)",
			in:   Input{Category: CategoryOther, Message: "Привет мир"},
		},
		{
			name:    "message of 9 Cyrillic runes is too short even though it has 18 bytes",
			in:      Input{Category: CategoryOther, Message: "Привет ми"},
			wantErr: true,
		},
		{
			name:    "empty message",
			in:      Input{Category: CategoryOther, Message: ""},
			wantErr: true,
		},
		{
			name:    "whitespace-only message",
			in:      Input{Category: CategoryOther, Message: "              "},
			wantErr: true,
		},
		{
			name:    "message is long only because of spaces around it",
			in:      Input{Category: CategoryOther, Message: "     abc     "},
			wantErr: true,
		},
		{
			name:    "unknown category",
			in:      Input{Category: Category("music"), Message: "Достаточно длинное сообщение"},
			wantErr: true,
		},
		{
			name:    "empty category",
			in:      Input{Message: "Достаточно длинное сообщение"},
			wantErr: true,
		},
		{
			name:    "organization without a name",
			in:      Input{Category: CategoryOrg, Message: "Достаточно длинное сообщение"},
			wantErr: true,
		},
		{
			name:    "organization name of one character",
			in:      Input{Category: CategoryOrg, Organization: "A", Message: "Достаточно длинное сообщение"},
			wantErr: true,
		},
		{
			name: "staff with unknown school",
			in: Input{
				Category: CategoryStaff, School: "Unknown", TeacherName: "Иванов И.И.",
				Message: "Достаточно длинное сообщение",
			},
			wantErr: true,
		},
		{
			name: "staff without school",
			in: Input{
				Category: CategoryStaff, TeacherName: "Иванов И.И.",
				Message: "Достаточно длинное сообщение",
			},
			wantErr: true,
		},
		{
			name: "staff with too short teacher name",
			in: Input{
				Category: CategoryStaff, School: "ШГ", TeacherName: "Ив",
				Message: "Достаточно длинное сообщение",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Validate()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidInput) {
					t.Fatalf("Validate() error = %v, want an error wrapping ErrInvalidInput", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Validate() unexpected error: %v", err)
			}
		})
	}
}

func TestInput_Validate_DoesNotModifyInput(t *testing.T) {
	in := Input{Category: CategoryOther, Message: "   Достаточно длинное сообщение   "}
	if err := in.Validate(); err != nil {
		t.Fatal(err)
	}
	if in.Message != "   Достаточно длинное сообщение   " {
		t.Errorf("Validate changed the caller's input: %q", in.Message)
	}
}
