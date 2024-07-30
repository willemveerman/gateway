// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package types

import (
	"fmt"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	resourcev3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"google.golang.org/protobuf/proto"

	"github.com/envoyproxy/gateway/internal/ir"
)

// XdsResources represents all the xds resources
type XdsResources = map[resourcev3.Type][]types.Resource

type EnvoyPatchPolicyStatuses []*ir.EnvoyPatchPolicyStatus

// ResourceVersionTable holds all the translated xds resources
type ResourceVersionTable struct {
	XdsResources
	EnvoyPatchPolicyStatuses
}

// DeepCopyInto copies the contents into the output object
// This was generated by controller-gen, moved from
// zz_generated.deepcopy.go and updated to use proto.Clone
// to deep copy the proto.Message
func (t *ResourceVersionTable) DeepCopyInto(out *ResourceVersionTable) {
	*out = *t
	if t.XdsResources != nil {
		in, out := &t.XdsResources, &out.XdsResources
		*out = make(map[string][]types.Resource, len(*in))
		for key, val := range *in {
			var outVal []types.Resource
			if val == nil {
				(*out)[key] = nil
			} else {
				// Snippet was generated by controller-gen
				// G601: Implicit memory aliasing in for loop.
				in, out := &val, &outVal //nolint:gosec,scopelint
				*out = make([]types.Resource, len(*in))
				for i := range *in {
					(*out)[i] = proto.Clone((*in)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy generates a deep copy of the ResourceVersionTable object.
// This was generated by controller-gen and moved over from
// zz_generated.deepcopy.go to this file.
func (t *ResourceVersionTable) DeepCopy() *ResourceVersionTable {
	if t == nil {
		return nil
	}
	out := new(ResourceVersionTable)
	t.DeepCopyInto(out)
	return out
}

// GetXdsResources retrieves the translated xds resources saved in the translator context.
func (t *ResourceVersionTable) GetXdsResources() XdsResources {
	return t.XdsResources
}

func (t *ResourceVersionTable) AddXdsResource(rType resourcev3.Type, xdsResource types.Resource) error {
	// It's a sanity check to make sure the xdsResource is not nil
	if xdsResource == nil {
		return fmt.Errorf("xds resource is nil")
	}

	// Perform type switch to handle different types of xdsResource
	switch rType {
	case resourcev3.ListenerType:
		// Handle Type specific operations
		if resourceOfType, ok := xdsResource.(*listenerv3.Listener); ok {
			if err := resourceOfType.ValidateAll(); err != nil {
				return fmt.Errorf("validation failed for xds resource %+v, err: %w", xdsResource, err)
			}
		} else {
			return fmt.Errorf("failed to cast xds resource %+v to Listener type", xdsResource)
		}
	case resourcev3.RouteType:
		// Handle Type specific operations
		if resourceOfType, ok := xdsResource.(*routev3.RouteConfiguration); ok {
			if err := resourceOfType.ValidateAll(); err != nil {
				return fmt.Errorf("validation failed for xds resource %+v, err: %w", xdsResource, err)
			}
		} else {
			return fmt.Errorf("failed to cast xds resource %+v to RouteConfiguration type", xdsResource)
		}

	case resourcev3.SecretType:
		// Handle specific operations
		if resourceOfType, ok := xdsResource.(*tlsv3.Secret); ok {
			if err := resourceOfType.ValidateAll(); err != nil {
				return fmt.Errorf("validation failed for xds resource %+v, err: %w", xdsResource, err)
			}
		} else {
			return fmt.Errorf("failed to cast xds resource %+v to Secret type", xdsResource)
		}

	case resourcev3.EndpointType:
		if resourceOfType, ok := xdsResource.(*endpointv3.ClusterLoadAssignment); ok {
			if err := resourceOfType.ValidateAll(); err != nil {
				return fmt.Errorf("validation failed for xds resource %+v, err: %w", xdsResource, err)
			}
		} else {
			return fmt.Errorf("failed to cast xds resource %+v to ClusterLoadAssignment type", xdsResource)
		}

	case resourcev3.ClusterType:
		// Handle specific operations
		if resourceOfType, ok := xdsResource.(*clusterv3.Cluster); ok {
			if err := resourceOfType.ValidateAll(); err != nil {
				return fmt.Errorf("validation failed for xds resource %+v, err: %w", xdsResource, err)
			}
		} else {
			return fmt.Errorf("failed to cast xds resource %+v to Cluster type", xdsResource)
		}
	case resourcev3.RateLimitConfigType:
		// Handle specific operations
		// cfg resource from runner.go is the RateLimitConfig type from "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3", which does have validate function.

		// Add more cases for other types as needed
	default:
		// Handle the case when the type is not recognized or supported
	}

	if t.XdsResources == nil {
		t.XdsResources = make(XdsResources)
	}
	if t.XdsResources[rType] == nil {
		t.XdsResources[rType] = make([]types.Resource, 0, 1)
	}

	t.XdsResources[rType] = append(t.XdsResources[rType], xdsResource)
	return nil
}

// AddOrReplaceXdsResource will update an existing resource of rType according to matchFunc or add as a new resource
// if none satisfy the match criteria. It will only update the first match it finds, regardless
// if multiple resources satisfy the match criteria.
func (t *ResourceVersionTable) AddOrReplaceXdsResource(rType resourcev3.Type, resource types.Resource, matchFunc func(existing types.Resource, new types.Resource) bool) error {
	if t.XdsResources == nil || t.XdsResources[rType] == nil {
		if err := t.AddXdsResource(rType, resource); err != nil {
			return err
		} else {
			return nil
		}
	}

	var found bool
	for i, r := range t.XdsResources[rType] {
		if matchFunc(r, resource) {
			t.XdsResources[rType][i] = resource
			found = true
			break
		}
	}
	if !found {
		if err := t.AddXdsResource(rType, resource); err != nil {
			return err
		} else {
			return nil
		}
	}
	return nil
}

// SetResources will update an entire entry of the XdsResources for a certain type to the provided resources
func (t *ResourceVersionTable) SetResources(rType resourcev3.Type, xdsResources []types.Resource) {
	if t.XdsResources == nil {
		t.XdsResources = make(XdsResources)
	}

	t.XdsResources[rType] = xdsResources
}

// Merge combines the resources from another ResourceVersionTable into this one
func (t *ResourceVersionTable) Merge(other *ResourceVersionTable) {
	if t.XdsResources == nil {
		t.XdsResources = make(XdsResources)
	}

	for rType, resources := range other.XdsResources {
		if t.XdsResources[rType] == nil {
			t.XdsResources[rType] = make([]types.Resource, 0)
		}
		t.XdsResources[rType] = append(t.XdsResources[rType], resources...)
	}

	t.EnvoyPatchPolicyStatuses = append(t.EnvoyPatchPolicyStatuses, other.EnvoyPatchPolicyStatuses...)
}
