package core

import (
	"fmt"
	"strings"

	"github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/constants"
)

type SelectParameters struct {
	Brand string
	Name  string
	Sex   string
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

func (p *SelectParameters) WithSex(sex string) *SelectParameters {
	p.Sex = sex
	return p
}

func (p *SelectParameters) getQuery() string {
	query := constants.Select
	conditions := []string{}

	parametersCount := 1
	if p.Brand != "" {
		conditions = append(conditions, fmt.Sprintf("perfumes.brand = $%d", parametersCount))
		parametersCount++
	}
	if p.Name != "" {
		conditions = append(conditions, fmt.Sprintf("perfumes.name = $%d", parametersCount))
		parametersCount++
	}
	if p.Sex != "" && p.Sex != "unisex" {
		conditions = append(conditions, fmt.Sprintf("(sexes.sex = $%d OR sexes.sex = 'unisex')", parametersCount))
		parametersCount++
	}

	if len(conditions) > 0 {
		return fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	return query
}

func (p *SelectParameters) unpack() []any {
	var args []any
	if p.Brand != "" {
		args = append(args, p.Brand)
	}
	if p.Name != "" {
		args = append(args, p.Name)
	}
	if p.Sex != "" && p.Sex != "unisex" {
		args = append(args, p.Sex)
	}
	return args
}
