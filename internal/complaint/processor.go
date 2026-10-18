package complaint

import "fmt"

type Processor interface {
	Process(in Input) (Result, error)
}

type KeywordProcessor struct {
	valid, invalid, urgent, medium keywordSet
}

func NewKeywordProcessor() *KeywordProcessor {
	return &KeywordProcessor{
		valid:   newKeywordSet(validKeywords),
		invalid: newKeywordSet(invalidKeywords),
		urgent:  newKeywordSet(urgentKeywords),
		medium:  newKeywordSet(mediumKeywords),
	}
}

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
