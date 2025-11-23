package matching

type Weights struct {
	upperNotesWeight float64
	coreNotesWeight  float64
	baseNotesWeight  float64

	familyWeight float64
	notesWeight  float64
	typeWeight   float64

	characteristicsWeight float64
	tagsWeight            float64
}

func NewBaseWeights(
	upperNotesWeight float64,
	coreNotesWeight float64,
	baseNotesWeight float64,
) *Weights {
	return &Weights{
		upperNotesWeight: upperNotesWeight,
		coreNotesWeight:  coreNotesWeight,
		baseNotesWeight:  baseNotesWeight,
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
	w.familyWeight = familyWeight
	w.notesWeight = notesWeight
	w.typeWeight = typeWeight
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
	w.characteristicsWeight = characteristicsWeight
	w.tagsWeight = tagsWeight
	return w
}
