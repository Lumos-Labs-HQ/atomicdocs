package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/yourusername/atomicdocs/internal/types"
)

type HonoParser struct {
	appURL string
}

func NewHonoParser(appURL string) *HonoParser {
	return &HonoParser{appURL: appURL}
}

func (p *HonoParser) Parse() ([]types.RouteInfo, error) {
	resp, err := http.Get(p.appURL + "/__atomicdocs_routes")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routes: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var routes []types.RouteInfo
	if err := json.Unmarshal(body, &routes); err != nil {
		return nil, err
	}
	
	return routes, nil
}
