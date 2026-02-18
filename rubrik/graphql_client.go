package rubrik

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/machinebox/graphql"
)

// GraphQLClient wraps the GraphQL client with authentication
type GraphQLClient struct {
	client   *graphql.Client
	endpoint string
	token    string
}

// NewGraphQLClient creates a new GraphQL client
func NewGraphQLClient(endpoint, token string) *GraphQLClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	client := graphql.NewClient(endpoint, graphql.WithHTTPClient(httpClient))

	return &GraphQLClient{
		client:   client,
		endpoint: endpoint,
		token:    token,
	}
}

// ExecuteQuery executes a GraphQL query with authentication
func (g *GraphQLClient) ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	req := graphql.NewRequest(query)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.token))

	// Add variables if provided
	for key, value := range variables {
		req.Var(key, value)
	}

	log.Printf("Executing GraphQL query (first 100 chars): %.100s", query)

	// Execute query
	if err := g.client.Run(ctx, req, result); err != nil {
		log.Printf("GraphQL query failed: %v", err)
		log.Printf("GraphQL endpoint: %s", g.endpoint)
		log.Printf("Auth token present: %t", g.token != "")
		return err
	}

	// Check if result is empty
	resultJSON, _ := json.Marshal(result)
	if len(resultJSON) <= 2 { // {} or null
		log.Printf("GraphQL query succeeded but result is empty: %s", string(resultJSON))
	} else {
		log.Printf("GraphQL query succeeded with data (size: %d bytes)", len(resultJSON))
	}

	return nil
}

// Example GraphQL queries for Rubrik CDM
const (
	// Get cluster information
	ClusterInfoQuery = `
	query ClusterInfo {
		cluster {
			id
			name
			version
			status
		}
	}`

	// Get system storage stats
	SystemStorageQuery = `
	query SystemStorage {
		system {
			storage {
				total
				used
				available
				snapshot
				liveMount
				miscellaneous
			}
		}
	}`

	// Get nodes
	NodesQuery = `
	query Nodes {
		nodes {
			id
			name
			status
			ipAddress
			cluster {
				id
				name
			}
		}
	}`

	// Get VMs
	VMwareVMsQuery = `
	query VMwareVMs {
		vmwareVms {
			edges {
				node {
					id
					name
					effectiveSlaDomain {
						id
						name
					}
				}
			}
		}
	}`

	// Get Nutanix VMs
	NutanixVMsQuery = `
	query NutanixVMs {
		nutanixVms {
			edges {
				node {
					id
					name
					effectiveSlaDomain {
						id
						name
					}
				}
			}
		}
	}`

	// Get Hyper-V VMs
	HypervVMsQuery = `
	query HypervVMs {
		hypervVms {
			edges {
				node {
					id
					name
					effectiveSlaDomain {
						id
						name
					}
				}
			}
		}
	}`

	// Get archive locations
	ArchiveLocationsQuery = `
	query ArchiveLocations {
		archiveLocations {
			edges {
				node {
					id
					name
					archivalLocationType
					status
				}
			}
		}
	}`

	// Get managed volumes
	ManagedVolumesQuery = `
	query ManagedVolumes {
		managedVolumes {
			edges {
				node {
					id
					name
					state
					numChannels
					configuredSlaDomainName
					effectiveSlaDomain {
						id
						name
					}
					primaryClusterId
					usedSize
					slaAssignment
					configuredSlaDomainId
					isWritable
					volumeSize
					effectiveSlaDomainName
					snapshotCount
					pendingSnapshotCount
					isRelic
				}
			}
		}
	}`

	// Get per VM storage stats
	PerVMStorageQuery = `
	query PerVMStorage {
		vmwareVms {
			edges {
				node {
					id
					name
					logicalBytes
					ingestedBytes
					exclusivePhysicalBytes
					sharedPhysicalBytes
					indexStorageBytes
				}
			}
		}
	}`

	// Get streams count
	StreamsCountQuery = `
	query StreamsCount {
		system {
			streams {
				count
			}
		}
	}`

	// Get data location usage
	DataLocationUsageQuery = `
	query DataLocationUsage {
		archiveLocations {
			edges {
				node {
					id
					name
					dataDownloaded
					dataArchived
					numVMsArchived
					numFilesetsArchived
					numLinuxFilesetsArchived
					numWindowsFilesetsArchived
					numShareFilesetsArchived
					numMssqlDbsArchived
					numHypervVmsArchived
					numNutanixVmsArchived
					numManagedVolumesArchived
				}
			}
		}
	}`

	// Get physical ingest time series
	PhysicalIngestTimeSeriesQuery = `
	query PhysicalIngestTimeSeries($range: String!) {
		system {
			physicalIngest {
				timeSeries(range: $range) {
					date
					value
				}
			}
		}
	}`

	// Get archival bandwidth time series
	ArchivalBandwidthTimeSeriesQuery = `
	query ArchivalBandwidthTimeSeries($range: String!) {
		system {
			archivalBandwidth {
				timeSeries(range: $range) {
					date
					value
				}
			}
		}
	}`

	// Get runway remaining
	RunwayRemainingQuery = `
	query RunwayRemaining {
		system {
			runwayRemaining
		}
	}`

	// Get average storage growth
	AverageStorageGrowthQuery = `
	query AverageStorageGrowth {
		system {
			averageStorageGrowthPerDay
		}
	}`

	// Get reports
	ReportsQuery = `
	query Reports {
		reports {
			edges {
				node {
					id
					name
					reportType
					status
				}
			}
		}
	}`
)

// Example response structures
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GraphQLError  `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message string `json:"message"`
	Path    []string `json:"path,omitempty"`
}

// Cluster response
type ClusterResponse struct {
	Cluster struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Version string `json:"version"`
		Status  string `json:"status"`
	} `json:"cluster"`
}

// System storage response
type SystemStorageResponse struct {
	System struct {
		Storage struct {
			Total         int `json:"total"`
			Used          int `json:"used"`
			Available     int `json:"available"`
			Snapshot      int `json:"snapshot"`
			LiveMount     int `json:"liveMount"`
			Miscellaneous int `json:"miscellaneous"`
		} `json:"storage"`
	} `json:"system"`
}

// Nodes response
type NodesResponse struct {
	Nodes []struct {
		ID              string `json:"id"`
		Name            string `json:"name"` // GraphQL uses 'name', REST uses 'brikId'
		Status          string `json:"status"`
		IPAddress       string `json:"ipAddress"`
		NeedsInspection bool   `json:"needsInspection"`
		Cluster         struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"cluster"`
	} `json:"nodes"`
}

// VMware VMs response
type VMwareVMsResponse struct {
	VmwareVms struct {
		Edges []struct {
			Node struct {
				ID                  string `json:"id"`
				Name                string `json:"name"`
				EffectiveSlaDomain  *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"effectiveSlaDomain"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"vmwareVms"`
}

// Nutanix VMs response
type NutanixVMsResponse struct {
	NutanixVms struct {
		Edges []struct {
			Node struct {
				ID                  string `json:"id"`
				Name                string `json:"name"`
				EffectiveSlaDomain  *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"effectiveSlaDomain"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"nutanixVms"`
}

// Hyper-V VMs response
type HypervVMsResponse struct {
	HypervVms struct {
		Edges []struct {
			Node struct {
				ID                  string `json:"id"`
				Name                string `json:"name"`
				EffectiveSlaDomain  *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"effectiveSlaDomain"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"hypervVms"`
}

// Archive locations response
type ArchiveLocationsResponse struct {
	ArchiveLocations struct {
		Edges []struct {
			Node struct {
				ID                    string `json:"id"`
				Name                  string `json:"name"`
				ArchivalLocationType  string `json:"archivalLocationType"`
				Status                string `json:"status"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"archiveLocations"`
}

// Managed volumes response
type ManagedVolumesResponse struct {
	ManagedVolumes struct {
		Edges []struct {
			Node struct {
				ID                      string  `json:"id"`
				Name                    string  `json:"name"`
				State                   string  `json:"state"`
				NumChannels             float64 `json:"numChannels"`
				ConfiguredSLADomainName string  `json:"configuredSlaDomainName"`
				EffectiveSlaDomain      *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"effectiveSlaDomain"`
				PrimaryClusterID        string  `json:"primaryClusterId"`
				UsedSize                float64 `json:"usedSize"`
				SlaAssignment           string  `json:"slaAssignment"`
				ConfiguredSLADomainID   string  `json:"configuredSlaDomainId"`
				IsWritable              string  `json:"isWritable"`
				VolumeSize              float64 `json:"volumeSize"`
				EffectiveSLADomainName  string  `json:"effectiveSlaDomainName"`
				SnapshotCount           float64 `json:"snapshotCount"`
				PendingSnapshotCount    float64 `json:"pendingSnapshotCount"`
				IsRelic                 string  `json:"isRelic"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"managedVolumes"`
}

// Per VM storage response
type PerVMStorageResponse struct {
	VmwareVms struct {
		Edges []struct {
			Node struct {
				ID                     string  `json:"id"`
				Name                   string  `json:"name"`
				LogicalBytes           float64 `json:"logicalBytes"`
				IngestedBytes          float64 `json:"ingestedBytes"`
				ExclusivePhysicalBytes float64 `json:"exclusivePhysicalBytes"`
				SharedPhysicalBytes    float64 `json:"sharedPhysicalBytes"`
				IndexStorageBytes      float64 `json:"indexStorageBytes"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"vmwareVms"`
}

// Streams count response
type StreamsCountResponse struct {
	System struct {
		Streams struct {
			Count int `json:"count"`
		} `json:"streams"`
	} `json:"system"`
}

// Data location usage response
type DataLocationUsageResponse struct {
	ArchiveLocations struct {
		Edges []struct {
			Node struct {
				ID                        string `json:"id"`
				Name                      string `json:"name"`
				DataDownloaded            int    `json:"dataDownloaded"`
				DataArchived              int    `json:"dataArchived"`
				NumVMsArchived            int    `json:"numVMsArchived"`
				NumFilesetsArchived       int    `json:"numFilesetsArchived"`
				NumLinuxFilesetsArchived  int    `json:"numLinuxFilesetsArchived"`
				NumWindowsFilesetsArchived int   `json:"numWindowsFilesetsArchived"`
				NumShareFilesetsArchived  int    `json:"numShareFilesetsArchived"`
				NumMssqlDbsArchived       int    `json:"numMssqlDbsArchived"`
				NumHypervVmsArchived      int    `json:"numHypervVmsArchived"`
				NumNutanixVmsArchived     int    `json:"numNutanixVmsArchived"`
				NumManagedVolumesArchived int    `json:"numManagedVolumesArchived"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"archiveLocations"`
}

// Time series data point
type TimeSeriesPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// Physical ingest time series response
type PhysicalIngestTimeSeriesResponse struct {
	System struct {
		PhysicalIngest struct {
			TimeSeries []TimeSeriesPoint `json:"timeSeries"`
		} `json:"physicalIngest"`
	} `json:"system"`
}

// Archival bandwidth time series response
type ArchivalBandwidthTimeSeriesResponse struct {
	System struct {
		ArchivalBandwidth struct {
			TimeSeries []TimeSeriesPoint `json:"timeSeries"`
		} `json:"archivalBandwidth"`
	} `json:"system"`
}

// Runway remaining response
type RunwayRemainingResponse struct {
	System struct {
		RunwayRemaining int `json:"runwayRemaining"`
	} `json:"system"`
}

// Average storage growth response
type AverageStorageGrowthResponse struct {
	System struct {
		AverageStorageGrowthPerDay float64 `json:"averageStorageGrowthPerDay"`
	} `json:"system"`
}

// Reports response
type ReportsResponse struct {
	Reports struct {
		Edges []struct {
			Node struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				ReportType string `json:"reportType"`
				Status     string `json:"status"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"reports"`
}