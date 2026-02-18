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

type VirtualMachine struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	EffectiveSLADomainID string `json:"effectiveSlaDomainId"`
}

type VirtualMachineList struct {
	*ResultList
	Data []VirtualMachine `json:"data"`
}

// ListAllVM retrieves a list of all Virtual Machine ID and Name
// for All kinds of hypervisors (vmware, nutanix, hyperv)
func (r Rubrik) ListAllVM() []VirtualMachine {
	var list []VirtualMachine
	list = append(list, r.ListVmwareVM()...)
	list = append(list, r.ListNutanixVM()...)
	list = append(list, r.ListHypervVM()...)

	return list
}

// ListVmwareVM retrieve a List of all known VMware VM's
func (r Rubrik) ListVmwareVM() []VirtualMachine {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response VMwareVMsResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), VMwareVMsQuery, nil, &response)
		if err == nil {
			// Convert GraphQL response to VirtualMachine structs
			vms := make([]VirtualMachine, len(response.VmwareVms.Edges))
			for i, edge := range response.VmwareVms.Edges {
				slaID := ""
				if edge.Node.EffectiveSlaDomain != nil {
					slaID = edge.Node.EffectiveSlaDomain.ID
				}
				vms[i] = VirtualMachine{
					ID:                   edge.Node.ID,
					Name:                 edge.Node.Name,
					EffectiveSLADomainID: slaID,
				}
			}
			return vms
		}
		log.Printf("GraphQL ListVmwareVM failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/v1/vmware/vm", RequestParams{})
	if err != nil || resp == nil {
		return []VirtualMachine{}
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	var s VirtualMachineList
	data.Decode(&s)
	return s.Data
}

// ListNutanixVM retrieve a List of all known Nutanix VM's
func (r Rubrik) ListNutanixVM() []VirtualMachine {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response NutanixVMsResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), NutanixVMsQuery, nil, &response)
		if err == nil {
			// Convert GraphQL response to VirtualMachine structs
			vms := make([]VirtualMachine, len(response.NutanixVms.Edges))
			for i, edge := range response.NutanixVms.Edges {
				slaID := ""
				if edge.Node.EffectiveSlaDomain != nil {
					slaID = edge.Node.EffectiveSlaDomain.ID
				}
				vms[i] = VirtualMachine{
					ID:                   edge.Node.ID,
					Name:                 edge.Node.Name,
					EffectiveSLADomainID: slaID,
				}
			}
			return vms
		}
		log.Printf("GraphQL ListNutanixVM failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/nutanix/vm", RequestParams{})
	if err != nil || resp == nil {
		return []VirtualMachine{}
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	var s VirtualMachineList
	data.Decode(&s)
	return s.Data
}

// ListHypervVM retrieve a List of all known Hyper-V VM's
func (r Rubrik) ListHypervVM() []VirtualMachine {
	// Try GraphQL first
	if r.graphqlClient != nil {
		var response HypervVMsResponse
		err := r.graphqlClient.ExecuteQuery(context.Background(), HypervVMsQuery, nil, &response)
		if err == nil {
			// Convert GraphQL response to VirtualMachine structs
			vms := make([]VirtualMachine, len(response.HypervVms.Edges))
			for i, edge := range response.HypervVms.Edges {
				slaID := ""
				if edge.Node.EffectiveSlaDomain != nil {
					slaID = edge.Node.EffectiveSlaDomain.ID
				}
				vms[i] = VirtualMachine{
					ID:                   edge.Node.ID,
					Name:                 edge.Node.Name,
					EffectiveSLADomainID: slaID,
				}
			}
			return vms
		}
		log.Printf("GraphQL ListHypervVM failed, falling back to REST: %v", err)
	}

	// Fallback to REST API
	resp, err := r.makeRequest("GET", "/api/internal/hyperv/vm", RequestParams{})
	if err != nil || resp == nil {
		return []VirtualMachine{}
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	var s VirtualMachineList
	data.Decode(&s)
	return s.Data
}
