package models

type GluedPerfume struct {
	Brand      string            `json:"brand"`
	Name       string            `json:"name"`
	Properties PerfumeProperties `json:"properties"`
	Links      map[int]string    `json:"links"`
}

func NewGluedPerfume(p Perfume) GluedPerfume {
	return GluedPerfume{
		Brand:      p.Brand,
		Name:       p.Name,
		Properties: p.getProperties(),
		Links:      map[int]string{p.Volume: p.Link},
	}
}
