// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * External DNS Webhook Server
 *
 * Implements the external DNS webhook endpoints.
 *
 * API version: v0.15.0
 */

package externaldnsapi

// Endpoint - This is a DNS record.
type Endpoint struct {
	DnsName string `json:"dnsName,omitempty"`

	// This is the list of targets that this DNS record points to. So for an A record it will be a list of IP addresses.
	Targets []string `json:"targets,omitempty"`

	RecordType string `json:"recordType,omitempty"`

	SetIdentifier string `json:"setIdentifier,omitempty"`

	RecordTTL int64 `json:"recordTTL,omitempty"`

	Labels map[string]string `json:"labels,omitempty"`

	ProviderSpecific []ProviderSpecificProperty `json:"providerSpecific,omitempty"`
}

// AssertEndpointRequired checks if the required fields are not zero-ed
func AssertEndpointRequired(obj Endpoint) error {
	for _, el := range obj.ProviderSpecific {
		if err := AssertProviderSpecificPropertyRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertEndpointConstraints checks if the values respects the defined constraints
func AssertEndpointConstraints(obj Endpoint) error {
	for _, el := range obj.ProviderSpecific {
		if err := AssertProviderSpecificPropertyConstraints(el); err != nil {
			return err
		}
	}
	return nil
}
