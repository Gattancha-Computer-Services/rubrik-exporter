package rubrik

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
)

type ReportList struct {
	*ResultList
	Data []Report `json:"data"`
}

type Report struct {
	Name           string `json:"name"`
	ReportType     string `json:"reportType"`
	UpdateTime     string `json:"updateTime"`
	ID             string `json:"id"`
	ReportTemplate string `json:"reportTemplate"`
	UpdateStatus   string `json:"updateStatus"`
}

type ReportData struct {
	ID          string             `json:"id"`
	Attribute   string             `json:"attribute"`
	ChartType   string             `json:"chartType"`
	Name        string             `json:"name"`
	Measure     string             `json:"measure"`
	DataColumns []ReportDataColumn `json:"dataColumns"`
}

type ReportDataColumn struct {
	Label      string            `json:"label"`
	DataPoints []ReportDataPoint `json:"dataPoints"`
}

type ReportDataPoint struct {
	Measure string  `json:"measure"`
	Value   float64 `json:"value"`
}

func (r Rubrik) GetReports(params map[string]string) []Report {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response ReportsResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), ReportsQuery, nil, &response)
		if err == nil {
			// Convert GraphQL response to Report structs
			reports := make([]Report, len(response.Reports.Edges))
			for i, edge := range response.Reports.Edges {
				reports[i] = Report{
					ID:         edge.Node.ID,
					Name:       edge.Node.Name,
					ReportType: edge.Node.ReportType,
					UpdateStatus: edge.Node.Status,
				}
			}
			return reports
		}
		log.Printf("GraphQL GetReports failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	_params := &RequestParams{params: url.Values{}}
	for k, v := range params {
		_params.params[k] = []string{v}
	}

	resp, err := r.makeRequest("GET", "/api/internal/report", *_params)
	if err != nil || resp == nil {
		return []Report{}
	}
	defer resp.Body.Close()

	var l ReportList
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&l)

	return l.Data
}

// GetTaskDetails - Returned the reported TaskStatus in last 24h
// returns  map[succeeded:3 failed:1 canceled:2]
func (r Rubrik) GetTaskDetails() map[string]float64 {
	reports := r.GetReports(map[string]string{
		"type": "Canned", "report_template": "ProtectionTasksDetails",
	})
	
	result := make(map[string]float64)
	
	// Return empty map if no reports found
	if len(reports) == 0 {
		return result
	}
	
	report := reports[0]

	_params := &RequestParams{params: url.Values{"chart_id": []string{"chart0"}}}
	_url := fmt.Sprintf("/api/internal/report/%s/chart", report.ID)

	resp, err := r.makeRequest("GET", _url, *_params)
	if err != nil || resp == nil {
		return result
	}
	defer resp.Body.Close()

	var data []ReportData
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	// Return empty map if no data found
	if len(data) == 0 {
		return result
	}

	for _, c := range data[0].DataColumns {
		if len(c.DataPoints) > 0 {
			_key := strings.ToLower(c.Label)
			result[_key] = c.DataPoints[0].Value
		}
	}

	return result
}
