package alexa

import (
	"github.com/dasjott/alexa-sdk-go/dialog"
)

const (
	ConfirmationStatusNone      = "NONE"
	ConfirmationStatusDenied    = "DENIED"
	ConfirmationStatusConfirmed = "CONFIRMED"
)

// Slot is a simple representation of the slot object from the echo request
type Slot struct {
	// ID is the first ID of that slot value
	ID string
	// Value is the first value of slot, not the actual spoken value
	Value string
	// Value is the value of slot, not the actual spoken value
	Values []SlotValue
	// Spoken is the actual value, spoken by the user
	Spoken string
	// ConfirmationStatus is the status of the confirmation of this slot
	ConfirmationStatus string
	// Match is true, if the actual speech is a match to one of this slots values or its synonyms
	Match bool
}

type SlotValue struct {
	ID    string
	Value string
}

// Empty determines whether this slot is already filled
func (s *Slot) Empty() bool {
	return s.Spoken == ""
}

func slotFromEchoSlot(es *dialog.EchoSlot) *Slot {
	var resolution SlotValue
	var match bool
	var builtin = true
	var values []SlotValue

	if es.Resolutions != nil && len(es.Resolutions.ResolutionsPerAuthority) > 0 {
		match = es.Resolutions.ResolutionsPerAuthority[0].IsMatch()
		builtin = es.Resolutions.ResolutionsPerAuthority[0].IsBuiltIn()
		for _, val := range es.Resolutions.ResolutionsPerAuthority[0].Values {
			values = append(values, SlotValue{val.Value.ID, val.Value.Name})
		}
		if len(values) > 0 {
			resolution = values[0]
		}
	}

	if resolution.Value == "" && builtin {
		resolution.Value = es.Value
	}

	return &Slot{
		ID:                 resolution.ID,
		Value:              resolution.Value,
		Values:             values,
		Spoken:             es.Value,
		ConfirmationStatus: es.ConfirmationStatus,
		Match:              match,
	}
}
