// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// IgnitionConfigParams ignition config params
//
// swagger:model ignition-config-params
type IgnitionConfigParams struct {

	// config
	Config string `json:"config,omitempty"`
}

// Validate validates this ignition config params
func (m *IgnitionConfigParams) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *IgnitionConfigParams) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IgnitionConfigParams) UnmarshalBinary(b []byte) error {
	var res IgnitionConfigParams
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
