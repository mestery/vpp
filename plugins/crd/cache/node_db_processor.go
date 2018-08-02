// Copyright (c) 2018 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import "github.com/pkg/errors"

//ProcessNodeResponses will read the nodeDTO map and make sure that each node has
//enough DTOS to fully process information. It then clears the node DTO map after it
//is finished with it.
func (p *ContivTelemetryProcessor) ProcessNodeResponses() {
	for data := range p.nodeResponseChannel {
		nodelist := p.ContivTelemetryCache.Cache.GetAllNodes()
		p.dtoList = append(p.dtoList, data)
		if len(p.dtoList) == numDTOs*len(nodelist) {
			p.SetNodeData()
			p.ValidateNodeInfo()
			p.dtoList = p.dtoList[0:0]
			p.ContivTelemetryCache.ClearCache()
		}
	}
}

// SetNodeData will iterate through the dtoList, read the type of dto, and assign the dto info to the name
// associated with the DTO.
func (p *ContivTelemetryProcessor) SetNodeData() {
	for _, data := range p.dtoList {
		if data.err != nil {
			p.ContivTelemetryCache.Cache.report = append(p.ContivTelemetryCache.Cache.report, errors.Errorf(
				"Node %+v has nodeDTO %+v and http error %s", data.NodeName, data, data.err).Error())
			continue
		}
		switch data.NodeInfo.(type) {
		case *NodeLiveness:
			nl := data.NodeInfo.(*NodeLiveness)
			p.ContivTelemetryCache.Cache.SetNodeLiveness(data.NodeName, nl)
		case *nodeInterfacesMapType:
			niDto := data.NodeInfo.(*nodeInterfacesMapType)
			p.ContivTelemetryCache.Cache.SetNodeInterfaces(data.NodeName, *niDto)
		case *nodeBridgeDomainMapTypes:
			nbdDto := data.NodeInfo.(*nodeBridgeDomainMapTypes)
			p.ContivTelemetryCache.Cache.SetNodeBridgeDomain(data.NodeName, *nbdDto)
		case *nodeL2FibMapTypes:
			nl2fDto := data.NodeInfo.(*nodeL2FibMapTypes)
			p.ContivTelemetryCache.Cache.SetNodeL2Fibs(data.NodeName, *nl2fDto)
		case *nodeTelemetryMapTypes:
			ntDto := data.NodeInfo.(*nodeTelemetryMapTypes)
			p.ContivTelemetryCache.Cache.SetNodeTelemetry(data.NodeName, *ntDto)
		case *nodeIPARPMapTypes:
			nipaDto := data.NodeInfo.(*nodeIPARPMapTypes)
			p.ContivTelemetryCache.Cache.SetNodeIPARPs(data.NodeName, *nipaDto)
		default:
			p.Log.Errorf("Node %+v has unknown data type: %+v", data.NodeName, data.NodeInfo)
		}

	}

}
