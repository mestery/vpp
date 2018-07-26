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

//ProcessNodeData reads data sent to the cache channel.
//It decides how to process the data received based on the type of Data Transfer Object.
//Then it updates the node with the name from the DTO with the specific data from the DTO.
func (p *ContivTelemetryProcessor) ProcessNodeData(nodename string) {
	for _, data := range p.dtoMap[nodename] {
		switch data.(type) {
		case NodeLivenessDTO:
			nlDto := data.(NodeLivenessDTO)
			p.Cache.SetNodeLiveness(nlDto.NodeName, nlDto.NodeInfo)
		case NodeInterfacesDTO:
			niDto := data.(NodeInterfacesDTO)
			p.Cache.SetNodeInterfaces(niDto.NodeName, niDto.NodeInfo)
		case NodeBridgeDomainsDTO:
			nbdDto := data.(NodeBridgeDomainsDTO)
			p.Cache.SetNodeBridgeDomain(nbdDto.NodeName, nbdDto.NodeInfo)
		case NodeL2FibsDTO:
			nl2fDto := data.(NodeL2FibsDTO)
			p.Cache.SetNodeL2Fibs(nl2fDto.NodeName, nl2fDto.NodeInfo)
		case NodeTelemetryDTO:
			ntDto := data.(NodeTelemetryDTO)
			p.Cache.SetNodeTelemetry(ntDto.NodeName, ntDto.NodeInfo)
		case NodeIPArpDTO:
			nipaDto := data.(NodeIPArpDTO)
			p.Cache.SetNodeIPARPs(nipaDto.NodeName, nipaDto.NodeInfo)
		default:
			p.Log.Error("Unknown data type")
		}
	}
	node, err := p.Cache.GetNode(nodename)
	if err != nil {
		p.Log.Error(err)
	}
	p.Cache.PopulateNodeMaps(node)
}

//ProcessNodeResponses will read the nodeDTO map and make sure that each node has
//enough DTOS to fully process information. It then clears the node DTO map after it
//is finished with it.
func (p *ContivTelemetryProcessor) ProcessNodeResponses() {
	for nodename := range p.nodeResponseChannel {
		if len(p.dtoMap[nodename]) == numDTOs {
			p.ProcessNodeData(nodename)
		}
		haveAllNetworkData := true
		for nodename := range p.dtoMap {
			if len(p.dtoMap[nodename]) != numDTOs {
				haveAllNetworkData = false
				break
			}
		}
		if haveAllNetworkData {
			p.ValidateNodeInfo()
			for nodename := range p.dtoMap {
				p.dtoMap[nodename] = p.dtoMap[nodename][0:0]
			}
		}
	}

}
