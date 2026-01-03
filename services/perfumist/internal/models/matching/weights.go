package matching

type Weights struct {
	UpperNotesWeight float64 `json:"upper_notes_weight,omitzero"`
	CoreNotesWeight  float64 `json:"core_notes_weight,omitzero"`
	BaseNotesWeight  float64 `json:"base_notes_weight,omitzero"`

	FamilyWeight float64 `json:"family_weight,omitzero"`
	NotesWeight  float64 `json:"notes_weight,omitzero"`
	TypeWeight   float64 `json:"type_weight,omitzero"`

	CharacteristicsWeight float64 `json:"characteristics_weight,omitzero"`
	TagsWeight            float64 `json:"tags_weight,omitzero"`
	OverlayWeight         float64 `json:"overlay_weight,omitzero"`
}

func NewBaseWeights(
	upperNotesWeight float64,
	coreNotesWeight float64,
	baseNotesWeight float64,
) *Weights {
	return &Weights{
		UpperNotesWeight: upperNotesWeight,
		CoreNotesWeight:  coreNotesWeight,
		BaseNotesWeight:  baseNotesWeight,
	}
}

func NewOverlayWeights(
	upperNotesWeight float64,
	coreNotesWeight float64,
	baseNotesWeight float64,
	familyWeight float64,
	notesWeight float64,
	typeWeight float64,
) *Weights {
	w := NewBaseWeights(upperNotesWeight, coreNotesWeight, baseNotesWeight)
	w.FamilyWeight = familyWeight
	w.NotesWeight = notesWeight
	w.TypeWeight = typeWeight
	return w
}

func NewSmartWeights(
	upperNotesWeight float64,
	coreNotesWeight float64,
	baseNotesWeight float64,
	characteristicsWeight float64,
	tagsWeight float64,
) *Weights {
	w := NewBaseWeights(upperNotesWeight, coreNotesWeight, baseNotesWeight)
	w.CharacteristicsWeight = characteristicsWeight
	w.TagsWeight = tagsWeight
	return w
}

func NewSmartEnhancedWeights(
	upperNotesWeight float64,
	coreNotesWeight float64,
	baseNotesWeight float64,
	characteristicsWeight float64,
	tagsWeight float64,
	overlayWeight float64,
) *Weights {
	w := NewSmartWeights(upperNotesWeight, coreNotesWeight, baseNotesWeight, characteristicsWeight, tagsWeight)
	w.OverlayWeight = overlayWeight
	return w
}
