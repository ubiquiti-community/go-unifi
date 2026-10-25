package unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
)

// DiffOptions adjusts Diff for keys that a per-key comparison cannot express.
type DiffOptions struct {
	// Clear lists keys that are sent as "" or [] when prior holds a non-empty
	// value and target omits the key or holds null. Diff returns an error for a
	// cleared key whose prior value is not a string or a list.
	Clear []string
	// Groups lists sets of keys that are sent together when any one of them is
	// sent, because the controller validates them against each other.
	Groups [][]string
}

// Diff returns the body of a sparse PUT that turns prior into target.
//
// The controller merges a PUT into the stored document one top-level key at a
// time and replaces nested objects and lists whole. Diff therefore compares the
// JSON encodings of prior and target per top-level key and returns the full
// target value of every key that differs. Keys that target omits are left out,
// except those listed in opts.Clear. _id and site_id are never returned.
func Diff[T any](prior, target *T, opts DiffOptions) (map[string]any, error) {
	before, err := jsonObject(prior)
	if err != nil {
		return nil, fmt.Errorf("unable to encode prior: %w", err)
	}
	after, err := jsonObject(target)
	if err != nil {
		return nil, fmt.Errorf("unable to encode target: %w", err)
	}
	delete(after, "_id")
	delete(after, "site_id")

	patch := map[string]any{}
	for k, v := range after {
		if old, ok := before[k]; !ok || !deepEqualJSON(old, v) {
			patch[k] = v
		}
	}

	for _, k := range opts.Clear {
		if v, ok := after[k]; ok && v != nil {
			continue
		}
		delete(patch, k)
		switch old := before[k].(type) {
		case nil:
		case string:
			if old != "" {
				patch[k] = ""
			}
		case []any:
			if len(old) > 0 {
				patch[k] = []any{}
			}
		default:
			return nil, fmt.Errorf("unable to clear %q: only strings and lists have a known empty value", k)
		}
	}

	for added := true; added; {
		added = false
		for _, group := range opts.Groups {
			if !slices.ContainsFunc(group, func(k string) bool { _, ok := patch[k]; return ok }) {
				continue
			}
			for _, k := range group {
				if _, sent := patch[k]; sent {
					continue
				}
				if v, ok := after[k]; ok {
					patch[k] = v
					added = true
				}
			}
		}
	}

	return patch, nil
}

// jsonObject encodes v and decodes it into a map, keeping numbers exact.
func jsonObject(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var m map[string]any
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("value encodes to null")
	}
	return m, nil
}

// UpdateWLANPatch PUTs patch, usually built by Diff, to the WLAN with the given
// id and returns the stored WLAN. An empty patch sends no PUT. A controller
// validation error is returned unwrapped with its api.err code.
func (c *ApiClient) UpdateWLANPatch(
	ctx context.Context,
	site, id string,
	patch map[string]any,
) (*WLAN, error) {
	return updatePatch(ctx, c, fmt.Sprintf("api/s/%s/rest/wlanconf/%s", site, id), patch,
		func() (*WLAN, error) { return c.getWLAN(ctx, site, id) })
}

// UpdateNetworkPatch is UpdateWLANPatch for the network with the given id.
func (c *ApiClient) UpdateNetworkPatch(
	ctx context.Context,
	site, id string,
	patch map[string]any,
) (*Network, error) {
	return updatePatch(ctx, c, fmt.Sprintf("api/s/%s/rest/networkconf/%s", site, id), patch,
		func() (*Network, error) { return c.getNetwork(ctx, site, id) })
}

func updatePatch[T any](
	ctx context.Context,
	c *ApiClient,
	path string,
	patch map[string]any,
	get func() (*T, error),
) (*T, error) {
	if len(patch) == 0 {
		return get()
	}

	var respBody struct {
		Meta meta `json:"meta"`
		Data []T  `json:"data"`
	}
	if err := c.do(ctx, http.MethodPut, path, patch, &respBody); err != nil {
		return nil, err
	}

	// The controller answers a PUT that changes nothing with an empty data array.
	if len(respBody.Data) == 0 {
		return get()
	}
	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}
	return &respBody.Data[0], nil
}
