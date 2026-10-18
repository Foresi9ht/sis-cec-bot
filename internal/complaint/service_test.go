package complaint

import (
	"errors"
	"testing"
)

type stubProcessor struct {
	result Result
}

func (s stubProcessor) Process(Input) (Result, error) { return s.result, nil }

type failingProcessor struct {
	err error
}

func (f failingProcessor) Process(Input) (Result, error) { return Result{}, f.err }

func newTestService() *Service {
	return NewService(NewKeywordProcessor())
}

func otherInput(message string) Input {
	return Input{Category: CategoryOther, Message: message}
}

func TestService_Create_Success(t *testing.T) {
	svc := newTestService()

	in := Input{
		Category:     CategoryOrg,
		Organization: "  Student Union  ",
		Message:      "   Руководство нарушает устав организации   ",
	}
	c, err := svc.Create(in)
	if err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if c.ID != 1 {
		t.Errorf("ID = %d, want 1", c.ID)
	}
	if c.Status != StatusPending {
		t.Errorf("Status = %s, want %s", c.Status, StatusPending)
	}
	if c.Result != nil {
		t.Errorf("Result = %+v, want nil before processing", c.Result)
	}
	if c.Input.Category != CategoryOrg {
		t.Errorf("Category = %s, want %s", c.Input.Category, CategoryOrg)
	}
	if c.Input.Organization != "Student Union" {
		t.Errorf("Organization = %q, want it trimmed", c.Input.Organization)
	}
	if c.Input.Message != "Руководство нарушает устав организации" {
		t.Errorf("Message = %q, want it trimmed", c.Input.Message)
	}
}

func TestService_Create_InvalidInput(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name string
		in   Input
	}{
		{"empty message", otherInput("")},
		{"too short message", otherInput("коротко")},
		{"unknown category", Input{Category: "music", Message: "Достаточно длинное сообщение"}},
		{"organization without a name", Input{Category: CategoryOrg, Message: "Достаточно длинное сообщение"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := svc.Create(tt.in)
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("Create() error = %v, want ErrInvalidInput", err)
			}
			if c.ID != 0 {
				t.Errorf("returned complaint has ID %d, want a zero value", c.ID)
			}
		})
	}

	if got := len(svc.List()); got != 0 {
		t.Errorf("%d complaints stored after only invalid input, want 0", got)
	}

	c, err := svc.Create(otherInput("Достаточно длинное сообщение"))
	if err != nil {
		t.Fatal(err)
	}
	if c.ID != 1 {
		t.Errorf("first stored complaint has ID %d, want 1", c.ID)
	}
}

func TestService_Get_Success(t *testing.T) {
	svc := newTestService()
	created, err := svc.Create(otherInput("Преподаватель угрожает отчислением"))
	if err != nil {
		t.Fatal(err)
	}

	got, err := svc.Get(created.ID)
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
	if got.ID != created.ID || got.Input != created.Input || got.Status != StatusPending {
		t.Errorf("Get() = %+v, want the created complaint %+v", got, created)
	}
}

func TestService_Get_NotFound(t *testing.T) {
	svc := newTestService()
	if _, err := svc.Create(otherInput("Достаточно длинное сообщение")); err != nil {
		t.Fatal(err)
	}

	for _, id := range []ID{0, 2, 42, -1} {
		_, err := svc.Get(id)
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("Get(%d) error = %v, want ErrNotFound", id, err)
		}
	}
}

func TestService_Process_Success(t *testing.T) {
	svc := newTestService()
	created, err := svc.Create(otherInput("Преподаватель вымогает взятку за пересдачу экзамена"))
	if err != nil {
		t.Fatal(err)
	}

	got, err := svc.Process(created.ID)
	if err != nil {
		t.Fatalf("Process() unexpected error: %v", err)
	}

	if got.Result == nil {
		t.Fatal("Result is nil after processing")
	}
	want := Result{Validity: ValidityValid, Priority: PriorityHigh, ValidHits: 3, InvalidHits: 0}
	if *got.Result != want {
		t.Errorf("Result = %+v, want %+v", *got.Result, want)
	}

	if got.Status != StatusCompleted {
		t.Errorf("Status = %s, want %s", got.Status, StatusCompleted)
	}

	stored, err := svc.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Status != StatusCompleted || stored.Result == nil || *stored.Result != want {
		t.Errorf("stored complaint = %+v, want status completed with result %+v", stored, want)
	}
}

func TestService_Process_NotFound(t *testing.T) {
	svc := newTestService()
	_, err := svc.Process(5)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Process() error = %v, want ErrNotFound", err)
	}
}

func TestService_Process_TwiceIsRejected(t *testing.T) {
	svc := newTestService()
	created, err := svc.Create(otherInput("Меня шантажируют угрозами"))
	if err != nil {
		t.Fatal(err)
	}
	first, err := svc.Process(created.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Process(created.ID)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("second Process() error = %v, want ErrInvalidTransition", err)
	}

	after, err := svc.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if after.Status != StatusCompleted {
		t.Errorf("Status = %s after a rejected operation, want %s", after.Status, StatusCompleted)
	}
	if after.Result == nil || *after.Result != *first.Result {
		t.Errorf("Result changed after a rejected operation: %+v, want %+v", after.Result, first.Result)
	}
}

func TestService_Process_ProcessorFailure(t *testing.T) {
	boom := errors.New("boom")
	svc := NewService(failingProcessor{err: boom})
	created, err := svc.Create(otherInput("Достаточно длинное сообщение"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.Process(created.ID)
	if !errors.Is(err, boom) {
		t.Fatalf("Process() error = %v, want it to wrap the processor error", err)
	}

	got, err := svc.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusFailed {
		t.Errorf("Status = %s, want %s", got.Status, StatusFailed)
	}
	if got.Result != nil {
		t.Errorf("Result = %+v, want nil for a failed complaint", got.Result)
	}

	if _, err := svc.Process(created.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("Process() of a failed complaint error = %v, want ErrInvalidTransition", err)
	}
}

func TestService_Process_UsesTheInjectedProcessor(t *testing.T) {
	want := Result{Validity: ValidityReview, Priority: PriorityMedium, ValidHits: 7, InvalidHits: 9}
	svc := NewService(stubProcessor{result: want})
	created, err := svc.Create(otherInput("Достаточно длинное сообщение"))
	if err != nil {
		t.Fatal(err)
	}

	got, err := svc.Process(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Result == nil || *got.Result != want {
		t.Errorf("Result = %+v, want the stub result %+v", got.Result, want)
	}
}

func TestService_MultipleIndependentComplaints(t *testing.T) {
	svc := newTestService()

	inputs := []Input{
		{Category: CategoryOrg, Organization: "Student Union", Message: "Руководство организации игнорирует студентов"},
		{Category: CategoryStaff, School: "ШИТиИ", TeacherName: "Иванов Иван", Message: "Преподаватель вымогает взятку"},
		{Category: CategoryCafeteria, Message: "Суп холодный, порции маленькие, меню однообразное"},
	}
	wantValidity := []Validity{ValidityValid, ValidityValid, ValidityInvalid}
	wantPriority := []Priority{PriorityLow, PriorityHigh, PriorityLow}

	ids := make([]ID, len(inputs))
	seen := make(map[ID]bool)
	for i, in := range inputs {
		c, err := svc.Create(in)
		if err != nil {
			t.Fatalf("Create(%d): %v", i, err)
		}
		if seen[c.ID] {
			t.Fatalf("duplicate ID %d", c.ID)
		}
		seen[c.ID] = true
		ids[i] = c.ID
	}

	if _, err := svc.Process(ids[1]); err != nil {
		t.Fatal(err)
	}
	for _, i := range []int{0, 2} {
		c, err := svc.Get(ids[i])
		if err != nil {
			t.Fatal(err)
		}
		if c.Status != StatusPending || c.Result != nil {
			t.Errorf("complaint %d changed while another was processed: %+v", ids[i], c)
		}
	}

	for _, i := range []int{0, 2} {
		if _, err := svc.Process(ids[i]); err != nil {
			t.Fatal(err)
		}
	}
	for i, id := range ids {
		c, err := svc.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		if c.Input != inputs[i] {
			t.Errorf("complaint %d input = %+v, want %+v", id, c.Input, inputs[i])
		}
		if c.Status != StatusCompleted || c.Result == nil {
			t.Fatalf("complaint %d = %+v, want completed with a result", id, c)
		}
		if c.Result.Validity != wantValidity[i] || c.Result.Priority != wantPriority[i] {
			t.Errorf("complaint %d result = %s/%s, want %s/%s",
				id, c.Result.Validity, c.Result.Priority, wantValidity[i], wantPriority[i])
		}
	}
}

func TestService_Get_ReturnsACopy(t *testing.T) {
	svc := newTestService()
	created, err := svc.Create(otherInput("Меня шантажируют угрозами"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Process(created.ID); err != nil {
		t.Fatal(err)
	}

	got, err := svc.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	got.Status = StatusPending
	got.Input.Message = "changed"
	got.Result.Priority = PriorityLow

	again, err := svc.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if again.Status != StatusCompleted {
		t.Errorf("stored Status = %s, the caller must not be able to change it", again.Status)
	}
	if again.Input.Message != "Меня шантажируют угрозами" {
		t.Errorf("stored Message = %q, the caller must not be able to change it", again.Input.Message)
	}
	if again.Result.Priority != PriorityHigh {
		t.Errorf("stored Priority = %s, the caller must not be able to change it", again.Result.Priority)
	}
}

func TestService_List_IsOrderedByID(t *testing.T) {
	svc := newTestService()
	const n = 30
	for i := 0; i < n; i++ {
		if _, err := svc.Create(otherInput("Достаточно длинное сообщение")); err != nil {
			t.Fatal(err)
		}
	}

	list := svc.List()
	if len(list) != n {
		t.Fatalf("len(List()) = %d, want %d", len(list), n)
	}
	for i, c := range list {
		if c.ID != ID(i+1) {
			t.Fatalf("List()[%d].ID = %d, want %d", i, c.ID, i+1)
		}
	}
}
