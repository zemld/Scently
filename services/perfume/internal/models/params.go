package models

import (
	"fmt"
	"strings"

	queries "github.com/zemld/PerfumeRecommendationSystem/perfume/internal/db/query"
)

type contextKey string

const UpdateParametersContextKey contextKey = "update_parameters"

type UpdateParameters struct {
	Perfumes []Perfume `json:"perfumes"`
}

func NewUpdateParameters() *UpdateParameters {
	return &UpdateParameters{}
}

func (p *UpdateParameters) WithPerfumes(perfumes []Perfume) *UpdateParameters {
	p.Perfumes = perfumes
	return p
}

const SelectParametersContextKey contextKey = "select_parameters"

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

func (p SelectParameters) GetQuery() string {
	choosingPerfumesQuery := p.GetChoosingPerfumesQuery()
	withClause := fmt.Sprintf(queries.WithSelect, choosingPerfumesQuery)
	return withClause + queries.EnrichSelectedPerfumes
}

func (p SelectParameters) GetChoosingPerfumesQuery() string {
	query := strings.TrimSpace(queries.SelectPerfumesBaseInfo)
	conditions := []string{}

	parametersCount := 1
	if p.Brand != "" {
		conditions = append(conditions, fmt.Sprintf("pb.brand = $%d", parametersCount))
		parametersCount++
	}
	if p.Name != "" {
		conditions = append(conditions, fmt.Sprintf("pb.name = $%d", parametersCount))
		parametersCount++
	}
	if p.Sex == "male" || p.Sex == "female" {
		conditions = append(conditions, fmt.Sprintf("(s.sex = $%d OR s.sex = 'unisex')", parametersCount))
		parametersCount++
	} else {
		conditions = append(conditions, "s.sex = 'unisex'")
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query = query + whereClause
	}

	return query
}

func (p SelectParameters) Unpack() []any {
	var args []any
	if p.Brand != "" {
		args = append(args, p.Brand)
	}
	if p.Name != "" {
		args = append(args, p.Name)
	}
	if p.Sex == "male" || p.Sex == "female" {
		args = append(args, p.Sex)
	}
	return args
}
