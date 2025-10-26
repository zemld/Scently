package parameters

type RequestPerfume struct {
	Brand string
	Name  string
	Sex   string
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

func (p *RequestPerfume) WithSex(sex string) *RequestPerfume {
	if sex != "male" && sex != "female" {
		return p
	}
	p.Sex = sex
	return p
}

func NewGet() *RequestPerfume {
	return &RequestPerfume{UseAI: false}
}
