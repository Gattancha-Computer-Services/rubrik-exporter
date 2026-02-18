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
	"io/ioutil"
	"log"
	"net/url"
)

type VmStorageList struct {
	*ResultList
	Data []VmStorage `json:"data"`
}

type VmStorage struct {
	ID                     string
	Logicalbytes           float64 `json:"logicalBytes"`
	IngestedBytes          float64 `json:"ingestedBytes"`
	ExclusivePhysicalBytes float64 `json:"exclusivePhysicalBytes"`
	SharedPhysicalBytes    float64 `json:"sharedPhysicalBytes"`
	IndexStorageBytes      float64 `json:"indexStorageBytes"`
}

type SystemStorage struct {
	Total         int `json:"total"`
	Used          int `json:"used"`
	Available     int `json:"available"`
	Snapshot      int `json:"snapshot"`
	LiveMount     int `json:"liveMount"`
	Miscellaneous int `json:"miscellaneous"`
}

type DataLocationUsageList struct {
	*ResultList
	Data []DataLocationUsage `json:"data"`
}

type DataLocationUsage struct {
	LocationID                 string `json:"locationId"`
	DataDownloaded             int    `json:"dataDownloaded"`
	DataArchived               int    `json:"dataArchived"`
	NumVMsArchived             int    `json:"numVMsArchived"`
	NumFilesetsArchived        int    `json:"numFilesetsArchived"`
	NumLinuxFilesetsArchived   int    `json:"numLinuxFilesetsArchived"`
	NumWindowsFilesetsArchived int    `json:"numWindowsFilesetsArchived"`
	NumShareFilesetsArchived   int    `json:"numShareFilesetsArchived"`
	NumMssqlDbsArchived        int    `json:"numMssqlDbsArchived"`
	NumHypervVmsArchived       int    `json:"numHypervVmsArchived"`
	NumNutanixVmsArchived      int    `json:"numNutanixVmsArchived"`
	NumManagedVolumesArchived  int    `json:"numManagedVolumesArchived"`
}

// GetSystemStorage ...
func (r Rubrik) GetSystemStorage() SystemStorage {
	log.Printf("GetSystemStorage: Starting API call")
	// Try GraphQL first
	if r.graphqlClient != nil {
		log.Printf("GetSystemStorage: Trying GraphQL")
		var response SystemStorageResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), SystemStorageQuery, nil, &response)
		if err == nil {
			log.Printf("GetSystemStorage: GraphQL succeeded, total: %d, used: %d", response.System.Storage.Total, response.System.Storage.Used)
			return response.System.Storage
		}
		log.Printf("GraphQL GetSystemStorage failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	log.Printf("GetSystemStorage: Falling back to REST API")
	resp, err := r.makeRequest("GET", "/api/internal/stats/system_storage", RequestParams{})
	if err != nil || resp == nil {
		log.Printf("GetSystemStorage: REST API failed: %v", err)
		return SystemStorage{}
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	var d SystemStorage
	err = data.Decode(&d)
	if err != nil {
		log.Printf("GetSystemStorage: Failed to decode REST response: %v", err)
		return SystemStorage{}
	}
	log.Printf("GetSystemStorage: REST succeeded, total: %d, used: %d", d.Total, d.Used)
	return d
}

// GetPerVMStorage ...
func (r Rubrik) GetPerVMStorage() []VmStorage {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response PerVMStorageResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), PerVMStorageQuery, nil, &response)
		if err == nil {
			// Convert GraphQL response to VmStorage structs
			storages := make([]VmStorage, len(response.VmwareVms.Edges))
			for i, edge := range response.VmwareVms.Edges {
				storages[i] = VmStorage{
					ID:                     edge.Node.ID,
					Logicalbytes:           edge.Node.LogicalBytes,
					IngestedBytes:          edge.Node.IngestedBytes,
					ExclusivePhysicalBytes: edge.Node.ExclusivePhysicalBytes,
					SharedPhysicalBytes:    edge.Node.SharedPhysicalBytes,
					IndexStorageBytes:      edge.Node.IndexStorageBytes,
				}
			}
			return storages
		}
		log.Printf("GraphQL GetPerVMStorage failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/stats/per_vm_storage", RequestParams{})
	if err != nil || resp == nil {
		return []VmStorage{}
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	var d VmStorageList
	data.Decode(&d)

	return d.Data
}

// GetStreamCount ...
func (r Rubrik) GetStreamCount() int {
	log.Printf("GetStreamCount: Starting API call")
	// Try GraphQL first
	if r.graphqlClient != nil {
		log.Printf("GetStreamCount: Trying GraphQL")
		var response StreamsCountResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), StreamsCountQuery, nil, &response)
		if err == nil {
			log.Printf("GetStreamCount: GraphQL succeeded, count: %d", response.System.Streams.Count)
			return response.System.Streams.Count
		}
		log.Printf("GetStreamCount: GraphQL failed, falling back to REST: %v", err)
	} else {
		log.Printf("GetStreamCount: GraphQL client not available, using REST")
	}

	// Fallback to REST API
	log.Printf("GetStreamCount: Trying REST API")
	resp, err := r.makeRequest("GET", "/api/internal/stats/streams/count", RequestParams{})
	if err != nil || resp == nil {
		log.Printf("GetStreamCount: REST API failed: %v", err)
		return 0
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]int
	json.Unmarshal(body, &data)
	count := data["count"]
	log.Printf("GetStreamCount: REST succeeded, count: %d", count)
	return count
}

// GetDataLocationUsage ...
func (r Rubrik) GetDataLocationUsage() []DataLocationUsage {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response DataLocationUsageResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), DataLocationUsageQuery, nil, &response)
		if err == nil {
			// Convert GraphQL response to DataLocationUsage structs
			usages := make([]DataLocationUsage, len(response.ArchiveLocations.Edges))
			for i, edge := range response.ArchiveLocations.Edges {
				usages[i] = DataLocationUsage{
					LocationID:                 edge.Node.ID,
					DataDownloaded:             edge.Node.DataDownloaded,
					DataArchived:               edge.Node.DataArchived,
					NumVMsArchived:             edge.Node.NumVMsArchived,
					NumFilesetsArchived:        edge.Node.NumFilesetsArchived,
					NumLinuxFilesetsArchived:   edge.Node.NumLinuxFilesetsArchived,
					NumWindowsFilesetsArchived: edge.Node.NumWindowsFilesetsArchived,
					NumShareFilesetsArchived:   edge.Node.NumShareFilesetsArchived,
					NumMssqlDbsArchived:        edge.Node.NumMssqlDbsArchived,
					NumHypervVmsArchived:       edge.Node.NumHypervVmsArchived,
					NumNutanixVmsArchived:      edge.Node.NumNutanixVmsArchived,
					NumManagedVolumesArchived:  edge.Node.NumManagedVolumesArchived,
				}
			}
			return usages
		}
		log.Printf("GraphQL GetDataLocationUsage failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/stats/data_location/usage", RequestParams{})
	if err != nil || resp == nil {
		return []DataLocationUsage{}
	}
	defer resp.Body.Close()

	var data DataLocationUsageList
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data.Data
}

func (r Rubrik) GetPhysicalIngest() []TimeStat {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response PhysicalIngestTimeSeriesResponse
		variables := map[string]interface{}{
			"range": "-10min",
		}
		err := r.graphqlClient.ExecuteQuery(context.Background(), PhysicalIngestTimeSeriesQuery, variables, &response)
		if err == nil {
			// Convert GraphQL response to TimeStat structs
			stats := make([]TimeStat, len(response.System.PhysicalIngest.TimeSeries))
			for i, point := range response.System.PhysicalIngest.TimeSeries {
				stats[i] = TimeStat{
					Time: point.Date,
					Stat: int(point.Value),
				}
			}
			return stats
		}
		log.Printf("GraphQL GetPhysicalIngest failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/stats/physical_ingest/time_series", RequestParams{params: url.Values{"range": []string{"-10min"}}})
	if err != nil || resp == nil {
		return []TimeStat{}
	}
	defer resp.Body.Close()

	var data []TimeStat
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data
}

func (r Rubrik) GetArchivalBandwith(locationID string, timerange string) []TimeStat {
	if timerange == "" {
		timerange = "-1h"
	}

	// Try GraphQL first
	if r.graphqlClient != nil {
		var response ArchivalBandwidthTimeSeriesResponse
		variables := map[string]interface{}{
			"range": timerange,
		}
		err := r.graphqlClient.ExecuteQuery(context.Background(), ArchivalBandwidthTimeSeriesQuery, variables, &response)
		if err == nil {
			// Convert GraphQL response to TimeStat structs
			stats := make([]TimeStat, len(response.System.ArchivalBandwidth.TimeSeries))
			for i, point := range response.System.ArchivalBandwidth.TimeSeries {
				stats[i] = TimeStat{
					Time: point.Date,
					Stat: int(point.Value),
				}
			}
			return stats
		}
		log.Printf("GraphQL GetArchivalBandwith failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/stats/archival/bandwidth/time_series",
		RequestParams{params: url.Values{"data_location_id": []string{locationID}, "range": []string{timerange}}})
	if err != nil || resp == nil {
		return []TimeStat{}
	}
	defer resp.Body.Close()

	var data []TimeStat
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data
}

// GetRunawayRemaining - Get the number of days remaining before the system fills up.
func (r Rubrik) GetRunawayRemaining() int {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response RunwayRemainingResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), RunwayRemainingQuery, nil, &response)
		if err == nil {
			return response.System.RunwayRemaining
		}
		log.Printf("GraphQL GetRunawayRemaining failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/stats/runway_remaining", RequestParams{})
	if err != nil || resp == nil {
		return 0
	}
	defer resp.Body.Close()

	var data map[string]int
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data["days"]
}

// GetAverageStorageGrowthPerDay - Get average storage growth per day.
func (r Rubrik) GetAverageStorageGrowthPerDay() int {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response AverageStorageGrowthResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), AverageStorageGrowthQuery, nil, &response)
		if err == nil {
			return int(response.System.AverageStorageGrowthPerDay)
		}
		log.Printf("GraphQL GetAverageStorageGrowthPerDay failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/stats/average_storage_growth_per_day", RequestParams{})
	if err != nil || resp == nil {
		return 0
	}
	defer resp.Body.Close()

	var data map[string]int
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data["bytes"]
}
