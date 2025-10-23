package parameters

type RequestPerfume struct {
	Brand string
	Name  string
	UseAI bool
}

func (p *RequestPerfume) WithBrand(brand string) *RequestPerfume {
	p.Brand = brand
	return p
}

func (p *RequestPerfume) WithName(name string) *RequestPerfume {
	p.Name = name
	return p
}

func (p *RequestPerfume) WithUseAI(useAI bool) *RequestPerfume {
	p.UseAI = useAI
	return p
}

func NewGet() *RequestPerfume {
	return &RequestPerfume{UseAI: false}
}
