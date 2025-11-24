package matching

import "log"

type AlgType string

const (
	OverlayAlg         AlgType = "overlay"
	TagsAlg            AlgType = "tags"
	CharacteristicsAlg AlgType = "characteristics"
	TagsOverlayAlg     AlgType = "tags_overlay"
	SmartAlg           AlgType = "smart"
)

func (a AlgType) String() string {
	return string(a)
}

func GetMatcherByAlg(alg AlgType, weights Weights) Matcher {
	switch alg {
	case OverlayAlg:
		return NewOverlay(weights)
	case TagsAlg:
		return NewTags(weights)
	case CharacteristicsAlg:
		return NewCharacteristicsMatcher(weights)
	case TagsOverlayAlg:
		return NewTagsMatcher(weights)
	case SmartAlg:
		return NewSmart(weights)
	default:
		log.Fatalf("unknown algorithm: %s", alg)
		return nil
	}
}
