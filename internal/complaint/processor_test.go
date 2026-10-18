package complaint

import (
	"errors"
	"testing"
)

func TestKeywordProcessor_Process(t *testing.T) {
	p := NewKeywordProcessor()

	tests := []struct {
		name         string
		message      string
		wantValidity Validity
		wantPriority Priority
	}{
		{
			name:         "bribery by a teacher is valid and urgent",
			message:      "Преподаватель вымогает взятку за пересдачу экзамена",
			wantValidity: ValidityValid,
			wantPriority: PriorityHigh,
		},
		{
			name:         "blackmail is valid and urgent",
			message:      "Меня шантажируют угрозами",
			wantValidity: ValidityValid,
			wantPriority: PriorityHigh,
		},
		{
			name:         "unfair grading is valid with medium priority",
			message:      "Преподаватель занижает баллы без объяснения",
			wantValidity: ValidityValid,
			wantPriority: PriorityMedium,
		},
		{
			name:         "a question about elections is valid with low priority",
			message:      "Хотел бы узнать про регламент студенческих выборов",
			wantValidity: ValidityValid,
			wantPriority: PriorityLow,
		},
		{
			name:         "cafeteria complaint is invalid",
			message:      "Суп в столовой холодный, порции маленькие, меню однообразное",
			wantValidity: ValidityInvalid,
			wantPriority: PriorityLow,
		},
		{
			name:         "nothing conclusive goes to manual review",
			message:      "Хочу рассказать про ситуацию на кафедре, там всё очень странно",
			wantValidity: ValidityReview,
			wantPriority: PriorityLow,
		},
		{
			name:         "a keyword inside another word is not a match (кран in экран)",
			message:      "Сломался экран проектора в аудитории",
			wantValidity: ValidityReview,
			wantPriority: PriorityLow,
		},
		{
			name:         "case and punctuation do not matter",
			message:      "ВЗЯТКА!!! СРОЧНО",
			wantValidity: ValidityValid,
			wantPriority: PriorityHigh,
		},
		{
			name:         "urgent beats medium",
			message:      "Жалоба: угроза расправы",
			wantValidity: ValidityValid,
			wantPriority: PriorityHigh,
		},
		{
			name:         "Kazakh text is supported",
			message:      "Оқытушы мені қорқытты",
			wantValidity: ValidityValid,
			wantPriority: PriorityHigh,
		},
		{
			name:         "English text is supported",
			message:      "The teacher showed bias in grading",
			wantValidity: ValidityValid,
			wantPriority: PriorityMedium,
		},
		{
			// 1 invalid keyword is not enough to reject.
			name:         "one invalid keyword and no valid ones goes to review",
			message:      "soup",
			wantValidity: ValidityReview,
			wantPriority: PriorityLow,
		},
		{
			// 2 invalid > 1 valid, and at least 2 invalid.
			name:         "more invalid than valid keywords is invalid",
			message:      "student soup menu",
			wantValidity: ValidityInvalid,
			wantPriority: PriorityLow,
		},
		{
			// 2 invalid == 2 valid: a tie is not enough to reject.
			name:         "a tie between valid and invalid keywords is valid",
			message:      "student rights soup menu",
			wantValidity: ValidityValid,
			wantPriority: PriorityLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Process(Input{Category: CategoryOther, Message: tt.message})
			if err != nil {
				t.Fatalf("Process() unexpected error: %v", err)
			}
			if got.Validity != tt.wantValidity {
				t.Errorf("Validity = %s, want %s (valid hits %d, invalid hits %d)",
					got.Validity, tt.wantValidity, got.ValidHits, got.InvalidHits)
			}
			if got.Priority != tt.wantPriority {
				t.Errorf("Priority = %s, want %s", got.Priority, tt.wantPriority)
			}
		})
	}
}

func TestKeywordProcessor_HitCounts(t *testing.T) {
	p := NewKeywordProcessor()

	got, err := p.Process(Input{Message: "Преподаватель вымогает взятку за пересдачу экзамена"})
	if err != nil {
		t.Fatal(err)
	}
	// "преподаватель", "взятк", "пересдач"
	if got.ValidHits != 3 || got.InvalidHits != 0 {
		t.Errorf("hits = valid %d, invalid %d; want valid 3, invalid 0", got.ValidHits, got.InvalidHits)
	}
}

func TestKeywordProcessor_IsDeterministic(t *testing.T) {
	p := NewKeywordProcessor()
	in := Input{Message: "Преподаватель занижает баллы, это несправедливо"}

	first, err := p.Process(in)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		got, err := p.Process(in)
		if err != nil {
			t.Fatal(err)
		}
		if got != first {
			t.Fatalf("run %d returned %+v, first run returned %+v", i, got, first)
		}
	}
}

func TestKeywordProcessor_EmptyMessage(t *testing.T) {
	p := NewKeywordProcessor()

	for _, msg := range []string{"", "   ", "!!! ... ???"} {
		_, err := p.Process(Input{Message: msg})
		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("Process(%q) error = %v, want ErrInvalidInput", msg, err)
		}
	}
}

func TestDecideValidity(t *testing.T) {
	tests := []struct {
		valid, invalid int
		want           Validity
	}{
		{0, 0, ValidityReview},
		{0, 1, ValidityReview},
		{1, 0, ValidityValid},
		{1, 1, ValidityValid},
		{0, 2, ValidityInvalid},
		{1, 2, ValidityInvalid},
		{2, 2, ValidityValid},
		{3, 2, ValidityValid},
		{2, 5, ValidityInvalid},
	}

	for _, tt := range tests {
		if got := decideValidity(tt.valid, tt.invalid); got != tt.want {
			t.Errorf("decideValidity(valid=%d, invalid=%d) = %s, want %s", tt.valid, tt.invalid, got, tt.want)
		}
	}
}

// KeywordProcessor must satisfy the Processor interface.
var _ Processor = (*KeywordProcessor)(nil)
