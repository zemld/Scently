package parameters

import (
	"net/http"

	"github.com/zemld/Scently/models"
	"github.com/zemld/Scently/perfumist/internal/errors"
)

type contextKey string

const ParamsKey contextKey = "params"

const (
	BrandParamKey = "brand"
	NameParamKey  = "name"
	SexParamKey   = "sex"
)

type RequestPerfume struct {
	Brand string
	Name  string
	Sex   models.Sex
}

func (p *RequestPerfume) WithBrand(brand string) *RequestPerfume {
	p.Brand = brand
	return p
}

func (p *RequestPerfume) WithName(name string) *RequestPerfume {
	p.Name = name
	return p
}

func (p *RequestPerfume) WithSex(sex models.Sex) *RequestPerfume {
	if sex != models.Male && sex != models.Female {
		return p
	}
	p.Sex = sex
	return p
}

func NewGet() *RequestPerfume {
	return &RequestPerfume{Sex: models.Unisex}
}

func (p RequestPerfume) AddToQuery(r *http.Request) {
	addQueryParameter(r, "brand", p.Brand)
	addQueryParameter(r, "name", p.Name)
	if p.Sex == models.Male || p.Sex == models.Female {
		addQueryParameter(r, "sex", string(p.Sex))
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
