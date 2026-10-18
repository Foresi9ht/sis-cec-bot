package complaint

import "fmt"

// Processor turns the input of a complaint into a result.
//
// It is the interface of the project: Service depends on it, not on a concrete
// classifier. That keeps Service independent of the way complaints are
// analysed and lets tests plug in a fake processor (for example one that
// always fails). A later checkpoint can run many processors in workers
// without changing Service's contract.
type Processor interface {
	Process(in Input) (Result, error)
}

// KeywordProcessor is the Processor used by the application. It classifies
// the message of a complaint with keyword lists.
//
// It has a pointer receiver because it owns precompiled keyword sets that
// should not be copied on every call. It is read-only after creation.
type KeywordProcessor struct {
	valid, invalid, urgent, medium keywordSet
}

// NewKeywordProcessor compiles the keyword lists and returns a ready processor.
func NewKeywordProcessor() *KeywordProcessor {
	return &KeywordProcessor{
		valid:   newKeywordSet(validKeywords),
		invalid: newKeywordSet(invalidKeywords),
		urgent:  newKeywordSet(urgentKeywords),
		medium:  newKeywordSet(mediumKeywords),
	}
}

// Process analyses in.Message and returns validity and priority.
//
// Validity rules:
//   - INVALID when there are more "invalid" than "valid" keywords and at least
//     two "invalid" ones (an operational matter, not a CEC case);
//   - otherwise VALID when at least one "valid" keyword is present;
//   - otherwise REVIEW (nothing conclusive, a human should look).
//
// Priority rules: HIGH when any urgent keyword is present, otherwise MEDIUM
// when any medium keyword is present, otherwise LOW.
func (p *KeywordProcessor) Process(in Input) (Result, error) {
	text := normalize(in.Message)
	if text == "" {
		return Result{}, fmt.Errorf("%w: message has no words to analyse", ErrInvalidInput)
	}
	padded := pad(text)

	validHits := p.valid.count(padded)
	invalidHits := p.invalid.count(padded)

	return Result{
		Validity:    decideValidity(validHits, invalidHits),
		Priority:    p.priority(padded),
		ValidHits:   validHits,
		InvalidHits: invalidHits,
	}, nil
}

func decideValidity(validHits, invalidHits int) Validity {
	switch {
	case invalidHits > validHits && invalidHits >= 2:
		return ValidityInvalid
	case validHits >= 1:
		return ValidityValid
	default:
		return ValidityReview
	}
}

func (p *KeywordProcessor) priority(padded string) Priority {
	switch {
	case p.urgent.any(padded):
		return PriorityHigh
	case p.medium.any(padded):
		return PriorityMedium
	default:
		return PriorityLow
	}
}
