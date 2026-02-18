//
// rubrik-exporter
//
// Exports metrics from rubrik backup for prometheus
//
// License: Apache License Version 2.0,
// Organization: Claranet GmbH
// Author: Martin Weber <martin.weber@de.clara.net>
//

package rubrik

import (
	"context"
	"encoding/json"
	"log"
)

type LocationList struct {
	*ResultList
	Data []Location `json:"data"`
}

type Location struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	LocationType string `json:"locationType"`
	IsActive     bool   `json:"isActive"`
	IPAddress    string `json:"ipAddress"`
	Bucket       string `json:"bucket"`
}

// GetArchiveLocations ...
func (r Rubrik) GetArchiveLocations() []Location {
	log.Printf("GetArchiveLocations: Starting API call")
	// Try GraphQL first
	if r.graphqlClient != nil {
		log.Printf("GetArchiveLocations: Trying GraphQL")
		var response ArchiveLocationsResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), ArchiveLocationsQuery, nil, &response)
		if err == nil {
			// Convert GraphQL response to Location structs
			locations := make([]Location, len(response.ArchiveLocations.Edges))
			for i, edge := range response.ArchiveLocations.Edges {
				locations[i] = Location{
					ID:           edge.Node.ID,
					Name:         edge.Node.Name,
					LocationType: edge.Node.ArchivalLocationType,
					IsActive:     edge.Node.Status == "CONNECTED", // Map status to isActive
				}
			}
			log.Printf("GetArchiveLocations: GraphQL succeeded, found %d locations", len(locations))
			return locations
		}
		log.Printf("GetArchiveLocations: GraphQL failed, falling back to REST: %v", err)
	} else {
		log.Printf("GetArchiveLocations: GraphQL client not available, using REST")
	}

	// Fallback to REST API
	log.Printf("GetArchiveLocations: Trying REST API")
	resp, err := r.makeRequest("GET", "/api/internal/archive/location", RequestParams{})
	if err != nil || resp == nil {
		log.Printf("GetArchiveLocations: REST API failed: %v", err)
		return []Location{}
	}
	defer resp.Body.Close()
	var data LocationList
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Printf("GetArchiveLocations: Failed to decode REST response: %v", err)
		return []Location{}
	}
	log.Printf("GetArchiveLocations: REST succeeded, found %d locations", len(data.Data))
	return data.Data
}
