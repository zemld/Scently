package parameters

import (
	"net/http"

	"github.com/zemld/Scently/perfumist/internal/errors"
)

type contextKey string

const ParamsKey contextKey = "params"

const (
	BrandParamKey = "brand"
	NameParamKey  = "name"
	SexParamKey   = "sex"
	UseAIParamKey = "use_ai"
)

const (
	SexMale   = "male"
	SexFemale = "female"
	SexUnisex = "unisex"
)

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

func (p RequestPerfume) AddToQuery(r *http.Request) {
	addQueryParameter(r, "brand", p.Brand)
	addQueryParameter(r, "name", p.Name)
	if p.Sex == "male" || p.Sex == "female" {
		addQueryParameter(r, "sex", p.Sex)
	}
}

func addQueryParameter(r *http.Request, key string, value string) {
	if value == "" {
		return
	}
	updatedQuery := r.URL.Query()
	updatedQuery.Set(key, value)
	r.URL.RawQuery = updatedQuery.Encode()
}

func (p RequestPerfume) Validate() error {
	if p.Brand == "" {
		return errors.NewValidationError("brand", "is required")
	}
	if p.Name == "" {
		return errors.NewValidationError("name", "is required")
	}
	return nil
}
