package core

import (
	"fmt"
	"reflect"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/constants"
)

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

func (p *SelectParameters) getQuery() string {
	query := constants.Select
	if p.Brand != "" && p.Name != "" {
		return fmt.Sprintf("%s WHERE perfumes.brand = $1 AND perfumes.name = $2", query)
	}
	if p.Brand != "" {
		return fmt.Sprintf("%s WHERE perfumes.brand = $1", query)
	}
	if p.Name != "" {
		return fmt.Sprintf("%s WHERE perfumes.name = $1", query)
	}
	return query
}

func (p *SelectParameters) unpack() []any {
	var args []any
	v := reflect.ValueOf(*p)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() != "" {
			args = append(args, v.Field(i).Interface())
		}
	}
	return args
}
