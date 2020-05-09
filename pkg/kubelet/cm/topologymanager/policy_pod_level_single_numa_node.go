/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package topologymanager

import (
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager/bitmask"
)

type podLevelSingleNumaNodePolicy struct {
	//List of NUMA Nodes available on the underlying machine
	numaNodes []int
}

var _ Policy = &podLevelSingleNumaNodePolicy{}

// PolicyPodLevelSingleNumaNode policy name.
const PolicyPodLevelSingleNumaNode string = "pod-level-single-numa-node"

// NewPodLevelSingleNumaNodePolicy returns pod-level-single-numa-node policy.
func NewPodLevelSingleNumaNodePolicy(numaNodes []int) Policy {
	klog.Infof("pod level single numa policy is created")
	return &podLevelSingleNumaNodePolicy{numaNodes: numaNodes}
}

func (p *podLevelSingleNumaNodePolicy) Name() string {
	return PolicyPodLevelSingleNumaNode
}

func (p *podLevelSingleNumaNodePolicy) canAdmitPodResult(hint *TopologyHint) bool {
	return hint.Preferred
}

func (p *podLevelSingleNumaNodePolicy) Merge(providersHints []map[string][]TopologyHint) (TopologyHint, bool) {
	klog.Infof("[topologymanager] bgchun podLevelSingleNumaNodePolicy.Merge is called providersHints: %v", providersHints)
	filteredHints := filterProvidersHints(providersHints)

	//[[{01/T, 10/T, 11/F }], [nil/T], [nil/F]]

	// Filter to only include don't cares and hints with a single NUMA node.
	//singleNumaHints := filterSingleNumaHints(filteredHints)

	bestHint := mergeFilteredHints(p.numaNodes, filteredHints)
	klog.Infof("[topologymanager] bgchun bestHint is %v", bestHint)

	defaultAffinity, _ := bitmask.NewBitMask(p.numaNodes...)
	klog.Infof("[topologymanager] bgchun defaultAffinity: %v", defaultAffinity)
	if bestHint.NUMANodeAffinity.IsEqual(defaultAffinity) {
		bestHint = TopologyHint{nil, bestHint.Preferred}
		klog.Infof("[topologymanager] bgchun bestHint is replaced to %v", bestHint)
	}

	admit := p.canAdmitPodResult(&bestHint)
	return bestHint, admit
}
