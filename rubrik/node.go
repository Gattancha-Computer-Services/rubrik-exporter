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
	"fmt"
	"log"
	"net/url"
)

type NodeList struct {
	*ResultList
	Data []Node `json:"data"`
}

// Node - Descripe a Rubrik node
type Node struct {
	ID              string `json:"id"`
	BrikID          string `json:"brikId"`
	Status          string `json:"status"`
	IPAddress       string `json:"ipAddress"`
	NeedsInspection bool   `json:"needsInspection"`
}

type NodeStat struct {
	ID              string       `json:"id"`
	BrikID          string       `json:"brikId"`
	Status          string       `json:"status"`
	IPAddress       string       `json:"ipAddress"`
	NeedsInspection bool         `json:"needsInspection"`
	NetworkStat     NetworkStat  `json:"networkStat"`
	Iops            Iops         `json:"iops"`
	IOThroughput    IoThroughput `json:"ioThroughput"`
	CPUStat         []TimeStat   `json:"cpuStat"`
}

// GetNodes - Returns the List of all Rubrik Nodes
func (r Rubrik) GetNodes() []Node {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response NodesResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), NodesQuery, nil, &response)
		if err == nil {
			// Convert GraphQL response to Node structs
			nodes := make([]Node, len(response.Nodes))
			for i, node := range response.Nodes {
				nodes[i] = Node{
					ID:              node.ID,
					BrikID:          node.Name, // GraphQL 'name' maps to 'brikId'
					Status:          node.Status,
					IPAddress:       node.IPAddress,
					NeedsInspection: node.NeedsInspection,
				}
			}
			return nodes
		}
		log.Printf("GraphQL GetNodes failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/node", RequestParams{})
	if err != nil || resp == nil {
		return []Node{}
	}
	defer resp.Body.Close()

	var l NodeList
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&l)

	return l.Data
}

// GetNodeStats ...
func (r Rubrik) GetNodeStats(id string) NodeStat {
	resp, err := r.makeRequest(
		"GET",
		fmt.Sprintf("/api/internal/node/%s/stats", id),
		RequestParams{params: url.Values{"range": []string{"-10min"}}})
	if err != nil || resp == nil {
		return NodeStat{}
	}
	defer resp.Body.Close()

	var result NodeStat
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&result)

	return result
}
