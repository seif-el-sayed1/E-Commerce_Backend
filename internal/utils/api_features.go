package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginationResult struct {
	CurrentPage  int   `json:"currentPage"`
	Limit        int   `json:"limit"`
	TotalPages   int   `json:"totalPages"`
	TotalResults int64 `json:"totalResults"`
	HasNextPage  bool  `json:"hasNextPage"`
	HasPrevPage  bool  `json:"hasPrevPage"`
}

type ApiFeatures struct {
	DB         *gorm.DB
	Query      map[string]string
	ModelName  string
	Pagination *PaginationResult

	page  int
	limit int
}

var searchFields = map[string][]string{
	"Admin": {"first_name", "last_name", "email"},
	"User":  {"first_name", "last_name", "email", "phone"},
}

func New(c *gin.Context, baseDB *gorm.DB, modelName string) *ApiFeatures {
	query := make(map[string]string, len(c.Request.URL.Query()))
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			query[key] = values[0]
		}
	}

	return &ApiFeatures{
		DB:        baseDB,
		Query:     query,
		ModelName: modelName,
		page:      1,
		limit:     20,
	}
}

func (af *ApiFeatures) Search() *ApiFeatures {
	keyword := strings.TrimSpace(af.Query["search"])
	if keyword == "" {
		return af
	}
	like := "%" + keyword + "%"

	fields := searchFields[af.ModelName]
	if len(fields) == 0 {
		return af
	}

	var conds []string
	var args []interface{}
	for _, f := range fields {
		conds = append(conds, f+" ILIKE ?")
		args = append(args, like)
	}

	af.DB = af.DB.Where(strings.Join(conds, " OR "), args...)
	return af
}

var reservedKeys = map[string]bool{
	"search": true, "page": true, "limit": true, "sort": true, "select": true,
	"startDate": true, "endDate": true,
}

var operatorSuffixes = []string{"gte", "gt", "lte", "lt"} // e.g. created_at_gte=2026-01-01

func parseValue(v string) interface{} {
	switch v {
	case "true":
		return true
	case "false":
		return false
	}
	if n, err := strconv.ParseFloat(v, 64); err == nil {
		return n
	}
	return v
}

func (af *ApiFeatures) Filter() *ApiFeatures {
	if start := af.Query["startDate"]; start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			af.DB = af.DB.Where("created_at >= ?", t)
		}
	}
	if end := af.Query["endDate"]; end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			af.DB = af.DB.Where("created_at <= ?", t)
		}
	}

	sqlOp := map[string]string{"gt": ">", "gte": ">=", "lt": "<", "lte": "<="}

	for key, raw := range af.Query {
		if reservedKeys[key] {
			continue
		}
		matched := false
		for _, suf := range operatorSuffixes {
			suffix := "_" + suf
			if strings.HasSuffix(key, suffix) {
				field := strings.TrimSuffix(key, suffix)
				af.DB = af.DB.Where(fmt.Sprintf("%s %s ?", field, sqlOp[suf]), parseValue(raw))
				matched = true
				break
			}
		}
		if !matched {
			af.DB = af.DB.Where(fmt.Sprintf("%s = ?", key), parseValue(raw))
		}
	}

	return af
}

func (af *ApiFeatures) Sort() *ApiFeatures {
	if af.Query["sort"] == "oldest" {
		af.DB = af.DB.Order("created_at ASC")
	} else {
		af.DB = af.DB.Order("created_at DESC")
	}
	return af
}

func (af *ApiFeatures) Paginate() *ApiFeatures {
	page := 1
	if p, err := strconv.Atoi(af.Query["page"]); err == nil && p > 0 {
		page = p
	}
	limit := 20
	if l, err := strconv.Atoi(af.Query["limit"]); err == nil && l > 0 {
		limit = l
	}

	af.page = page
	af.limit = limit

	af.DB = af.DB.Offset((page - 1) * limit).Limit(limit)
	return af
}

func (af *ApiFeatures) CalculatePagination() (*ApiFeatures, error) {
	var total int64
	if err := af.DB.Session(&gorm.Session{}).Offset(-1).Limit(-1).Count(&total).Error; err != nil {
		return af, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(af.limit)))

	af.Pagination = &PaginationResult{
		CurrentPage:  af.page,
		Limit:        af.limit,
		TotalPages:   totalPages,
		TotalResults: total,
		HasNextPage:  af.page < totalPages,
		HasPrevPage:  af.page > 1,
	}

	return af, nil
}

func (af *ApiFeatures) Execute(dest interface{}) error {
	return af.DB.Find(dest).Error
}
