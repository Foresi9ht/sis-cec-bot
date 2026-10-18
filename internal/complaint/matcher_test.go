package complaint

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Wi-Fi не работает!", "wi fi не работает"},
		{"  ВЗЯТКА!!!   СРОЧНО  ", "взятка срочно"},
		{"Қазақша сөйлем, Әділетсіз", "қазақша сөйлем әділетсіз"},
		{"123 abc", "123 abc"},
		{"   ...   ", ""},
		{"", ""},
	}

	for _, tt := range tests {
		if got := normalize(tt.in); got != tt.want {
			t.Errorf("normalize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestKeywordSet_Count(t *testing.T) {
	tests := []struct {
		name     string
		keywords []string
		text     string
		want     int
	}{
		{"stem matches the start of a word", []string{"нарушен"}, "нарушения прав", 1},
		{"stem does not match inside a word", []string{"кран"}, "сломался экран", 0},
		{"stem matches a longer word", []string{"кран"}, "кранов нет", 1},
		{"whole-word keyword matches the word", []string{"суп$"}, "суп холодный", 1},
		{"whole-word keyword at the end of text", []string{"суп$"}, "холодный суп", 1},
		{"whole-word keyword ignores punctuation", []string{"суп$"}, "суп, холодный", 1},
		{"whole-word keyword does not match a longer word", []string{"суп$"}, "супер идея", 0},
		{"multi-word keyword", []string{"качество еды"}, "качество еды ужасное", 1},
		{"multi-word keyword needs both words", []string{"качество еды"}, "качество обучения", 0},
		{"keyword with a hyphen matches the same text", []string{"wi-fi"}, "Wi-Fi не работает", 1},
		{"duplicates are counted once", []string{"нарушен", "нарушен"}, "нарушения", 1},
		{"distinct keywords are counted separately", []string{"нарушен", "угроз"}, "угрозы и нарушения", 2},
		{"empty keywords are ignored", []string{"", "$", "-"}, "любой текст", 0},
		{"no keywords", nil, "любой текст", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := newKeywordSet(tt.keywords)
			if got := set.count(pad(normalize(tt.text))); got != tt.want {
				t.Errorf("count = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestNewKeywordSet_RemovesDuplicates(t *testing.T) {
	set := newKeywordSet([]string{"нарушен", "Нарушен", "нарушен", "угроз"})
	if len(set) != 2 {
		t.Errorf("len = %d, want 2", len(set))
	}
}

func TestKeywordSet_Any(t *testing.T) {
	set := newKeywordSet([]string{"взятк", "срочно"})
	if !set.any(pad(normalize("Это СРОЧНО"))) {
		t.Error("any = false, want true")
	}
	if set.any(pad(normalize("Ничего особенного"))) {
		t.Error("any = true, want false")
	}
}
