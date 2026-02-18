package rubrik

import (
	"context"
	"encoding/json"
	"log"
)

type ManagedVolumeList struct {
	*ResultList
	Data []ManagedVolume `json:"data"`
}

type ManagedVolume struct {
	ID                      string  `json:"id"`
	State                   string  `json:"state"`
	NumChannels             float64 `json:"numChannels"`
	ConfiguredSLADomainName string  `json:"configuredSlaDomainName"`
	EffectiveSLADomainID    string  `json:"effectiveSlaDomainId"`
	PrimaryClusterID        string  `json:"primaryClusterId"`
	UsedSize                float64 `json:"usedSize"`
	SLAAssignment           string  `json:"slaAssignment"`
	// MainExport              string  `json:"mainExport"`
	ConfiguredSLADomainID  string  `json:"configuredSlaDomainId"`
	IsWritable             string  `json:"isWritable"`
	VolumeSize             float64 `json:"volumeSize"`
	EffectiveSLADomainName string  `json:"effectiveSlaDomainName"`
	SnapshotCount          float64 `json:"snapshotCount"`
	PendingSnapshotCount   float64 `json:"pendingSnapshotCount"`
	IsRelic                string  `json:"isRelic"`
	Name                   string  `json:"name"`
	// HostPatterns            string  `json:"hostPatterns"`
	// Links string `json:"links"`

}

/* GetManagedVolumes
 *
 */
func (r Rubrik) GetManagedVolumes() []ManagedVolume {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response ManagedVolumesResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), ManagedVolumesQuery, nil, &response)
		if err == nil {
			// Convert GraphQL response to ManagedVolume structs
			volumes := make([]ManagedVolume, len(response.ManagedVolumes.Edges))
			for i, edge := range response.ManagedVolumes.Edges {
				effectiveSlaID := ""
				if edge.Node.EffectiveSlaDomain != nil {
					effectiveSlaID = edge.Node.EffectiveSlaDomain.ID
				}
				volumes[i] = ManagedVolume{
					ID:                      edge.Node.ID,
					Name:                    edge.Node.Name,
					State:                   edge.Node.State,
					NumChannels:             edge.Node.NumChannels,
					ConfiguredSLADomainName: edge.Node.ConfiguredSLADomainName,
					EffectiveSLADomainID:    effectiveSlaID,
					PrimaryClusterID:        edge.Node.PrimaryClusterID,
					UsedSize:                edge.Node.UsedSize,
					SLAAssignment:           edge.Node.SlaAssignment,
					ConfiguredSLADomainID:   edge.Node.ConfiguredSLADomainID,
					IsWritable:              edge.Node.IsWritable,
					VolumeSize:              edge.Node.VolumeSize,
					EffectiveSLADomainName:  edge.Node.EffectiveSLADomainName,
					SnapshotCount:           edge.Node.SnapshotCount,
					PendingSnapshotCount:    edge.Node.PendingSnapshotCount,
					IsRelic:                 edge.Node.IsRelic,
				}
			}
			return volumes
		}
		log.Printf("GraphQL GetManagedVolumes failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/managed_volume", RequestParams{})
	if err != nil || resp == nil {
		return []ManagedVolume{}
	}
	defer resp.Body.Close()

	var l ManagedVolumeList
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&l)

	return l.Data
}
