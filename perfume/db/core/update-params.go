package core

type UpdateParameters struct {
	IsTruncate bool
}

func NewUpdateParameters() *UpdateParameters {
	return &UpdateParameters{
		IsTruncate: false,
	}
}

func (p *UpdateParameters) WithTruncate() *UpdateParameters {
	p.IsTruncate = true
	return p
}
