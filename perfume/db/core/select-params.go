package core

type SelectParameters struct {
	Brand string
	Name  string
}

func NewSelectParameters() *SelectParameters {
	return &SelectParameters{}
}

func (p *SelectParameters) WithBrand(brand string) *SelectParameters {
	p.Brand = brand
	return p
}

func (p *SelectParameters) WithName(name string) *SelectParameters {
	p.Name = name
	return p
}
