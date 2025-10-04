package util

type GetParameters struct {
	Brand string
	Name  string
}

func (p *GetParameters) WithBrand(brand string) *GetParameters {
	p.Brand = brand
	return p
}

func (p *GetParameters) WithName(name string) *GetParameters {
	p.Name = name
	return p
}

func NewGetParameters() *GetParameters {
	return &GetParameters{}
}
