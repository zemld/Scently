package models

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/zemld/Scently/models"
	queries "github.com/zemld/Scently/perfume-hub/internal/db/query"
)

const DefaultItemsPerPage = 500

type contextKey string

const UpdateParametersContextKey contextKey = "update_parameters"

type UpdateParameters struct {
	Perfumes []models.Perfume `json:"perfumes"`
}

func NewUpdateParameters() *UpdateParameters {
	return &UpdateParameters{}
}

func (p *UpdateParameters) WithPerfumes(perfumes []models.Perfume) *UpdateParameters {
	p.Perfumes = perfumes
	return p
}

const SelectParametersContextKey contextKey = "select_parameters"

type SelectParameters struct {
	Brand           string
	Name            string
	Sex             string
	Page            int
	parametersCount int
}

func NewSelectParameters() *SelectParameters {
	return &SelectParameters{parametersCount: 1}
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

func (p *SelectParameters) WithPage(page int) *SelectParameters {
	if page <= 0 {
		page = 1
	}
	p.Page = page
	return p
}

func (p *SelectParameters) GetQuery() string {
	choosingPerfumesQuery := p.GetChoosingPerfumesQuery()
	withClause := fmt.Sprintf(queries.WithSelect, choosingPerfumesQuery)
	log.Printf("withClause: %s", withClause)
	return withClause + queries.EnrichSelectedPerfumes
}

func (p *SelectParameters) GetChoosingPerfumesQuery() string {
	if p.Brand == "" && p.Name == "" {
		return p.GetAllPerfumesQueryWithPage()
	}
	return p.GetConcretePerfumeQuery()
}

func (p *SelectParameters) GetAllPerfumesQueryWithPage() string {
	query := p.updateQueryWithSexFilter(strings.TrimSpace(queries.SelectPerfumesBaseInfo))
	log.Printf("query: %s", query)
	return query + fmt.Sprintf(" AND page_number = $%d", p.parametersCount)
}

func (p *SelectParameters) GetConcretePerfumeQuery() string {
	query := p.updateQueryWithSexFilter(strings.TrimSpace(queries.SelectPerfumesBaseInfo))
	if p.Brand != "" {
		query += fmt.Sprintf(" AND canonized_brand = $%d", p.parametersCount)
		p.parametersCount++
	}
	if p.Name != "" {
		query += fmt.Sprintf(" AND canonized_name = $%d", p.parametersCount)
		p.parametersCount++
	}
	return query
}

func (p *SelectParameters) updateQueryWithSexFilter(query string) string {
	query += " WHERE"
	if p.Sex == "male" || p.Sex == "female" {
		query += fmt.Sprintf(" (sex = 'unisex' OR sex = $%d)", p.parametersCount)
		p.parametersCount++
	} else {
		query += " sex = 'unisex'"
	}
	return query
}

func (p SelectParameters) Unpack() []any {
	var args []any
	if p.Sex == "male" || p.Sex == "female" {
		args = append(args, p.Sex)
	}
	if p.Brand != "" {
		args = append(args, canonize(p.Brand))
	}
	if p.Name != "" {
		args = append(args, canonize(p.Name))
	}
	if len(args) <= 1 {
		args = append(args, p.Page)
	}
	return args
}

func canonize(s string) string {
	canonized := strings.Builder{}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			continue
		}
		canonized.WriteRune(unicode.ToLower(r))
	}
	return canonized.String()
}
