package collective

import (
	"fmt"
	"strings"
	"sync"
)

type Device struct {
	Name string
}

type DeviceType struct {
	TypeString string
}

type CollGroupRuntimeDetails struct {
	CommunicatorKey string
}

func (c CollGroupRuntimeDetails) String() string {
	return fmt.Sprintf("CollGroupRuntimeDetails {communicator_key=%s}", c.CommunicatorKey)
}

type CollGroupParams struct {
	GroupKey            int
	GroupSize           int
	DeviceType          DeviceType
	NumTasks            int
	RuntimeDetails      CollGroupRuntimeDetails
	Members             []Device
	NumDevicesPerTask   map[int]int
}

func (p CollGroupParams) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("CollGroupParams {group_key=%d group_size=%d device_type=%s num_tasks=%d runtime_details=%s devices {",
		p.GroupKey, p.GroupSize, p.DeviceType.TypeString, p.NumTasks, p.RuntimeDetails.String()))
	for _, m := range p.Members {
		b.WriteString(m.Name + ",")
	}
	b.WriteString("} num_devices_per_task={")
	for k, v := range p.NumDevicesPerTask {
		b.WriteString(fmt.Sprintf("%d: %d, ", k, v))
	}
	b.WriteString("}")
	return b.String()
}

type ImplDetails struct {
	CollectiveName       string
	SubdivOffsets        []int
	SubdivPermutations   [][]int
	SubdivSourceRank     []int
	Dependencies         []string
}

type CollInstanceParams struct {
	InstanceKey int
	Type        string
	DataType    string
	Shape       string
	ImplDetails ImplDetails
	Devices     []string
	Permutation []int
}

func (p CollInstanceParams) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("CollInstanceParams { instance_key=%d type=%s data_type=%s shape=%s devices {",
		p.InstanceKey, p.Type, p.DataType, p.Shape))
	b.WriteString("} collective_name=" + p.ImplDetails.CollectiveName + ", subdiv_offsets={")
	for _, d := range p.ImplDetails.SubdivOffsets {
		b.WriteString(fmt.Sprintf("%d,", d))
	}
	b.WriteString("}, subdiv_perms={")
	for _, perm := range p.ImplDetails.SubdivPermutations {
		b.WriteString("{")
		for _, i := range perm {
			b.WriteString(fmt.Sprintf("%d,", i))
		}
		b.WriteString("}")
	}
	if len(p.ImplDetails.SubdivSourceRank) > 0 {
		b.WriteString(" subdiv_source_rank={")
		for _, r := range p.ImplDetails.SubdivSourceRank {
			b.WriteString(fmt.Sprintf("%d,", r))
		}
		b.WriteString("}")
	}
	if p.Type == "PERMUTE_COLLECTIVE" {
		b.WriteString("}, permute_devices {")
		for _, d := range p.Devices {
			b.WriteString(d + ",")
		}
		b.WriteString("}, permute_permutation {")
		for _, p := range p.Permutation {
			b.WriteString(fmt.Sprintf("%d,", p))
		}
		b.WriteString("}")
	}
	b.WriteString("}")
	return b.String()
}

type CollectiveParams struct {
	Name         string
	Group        CollGroupParams
	Instance     CollInstanceParams
	DefaultRank  int
	IsSource     bool
	SourceRank   int
	SubdivRank   []int
}

func (p CollectiveParams) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("CollectiveParams %s {%s %s", p.Name, p.Group.String(), p.Instance.String()))
	b.WriteString(fmt.Sprintf(" default_rank=%d is_source=%t source_rank=%d subdiv_rank={", p.DefaultRank, p.IsSource, p.SourceRank))
	for _, r := range p.SubdivRank {
		b.WriteString(fmt.Sprintf("%d,", r))
	}
	b.WriteString("}}")
	return b.String()
}

// ---- Registry logic ----

type CollectiveImplementation interface{}

type Factory func() CollectiveImplementation

type RegistrationInfo struct {
	Name                   string
	Factory                Factory
	ParamResolverInstance  CollectiveImplementation
}

var (
	registryMu sync.Mutex
	registry   []RegistrationInfo
)

func Register(name string, factory Factory) error {
	registryMu.Lock()
	defer registryMu.Unlock()

	for _, reg := range registry {
		if reg.Name == name {
			return fmt.Errorf("Already registered collective %s", name)
		}
	}
	registry = append(registry, RegistrationInfo{
		Name: name,
		Factory: factory,
		ParamResolverInstance: factory(),
	})
	return nil
}

func Lookup(name string) (CollectiveImplementation, error) {
	registryMu.Lock()
	defer registryMu.Unlock()
	for _, reg := range registry {
		if reg.Name == name {
			return reg.Factory(), nil
		}
	}
	return nil, fmt.Errorf("No such collective implementation: %s", name)
}

func LookupParamResolver(name string) (CollectiveImplementation, error) {
	registryMu.Lock()
	defer registryMu.Unlock()
	for _, reg := range registry {
		if reg.Name == name {
			return reg.ParamResolverInstance, nil
		}
	}
	return nil, fmt.Errorf("No param resolver for: %s", name)
}