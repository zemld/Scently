package models

type contextKey string

const UpdateParametersContextKey contextKey = "update_parameters"

type UpdateParameters struct {
	Perfumes []UpgradedPerfume `json:"perfumes"`
}

func NewUpdateParameters() *UpdateParameters {
	return &UpdateParameters{}
}

func (p *UpdateParameters) WithPerfumes(perfumes []UpgradedPerfume) *UpdateParameters {
	p.Perfumes = perfumes
	return p
}
