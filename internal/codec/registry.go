// Package codec owns native format registration, detection, text decoding, and
// model-driven decode and render capability lookup.
package codec

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/source"
)

// Format is one canonical schema format identity.
type Format string

const (
	FormatAuto   Format = "auto"
	FormatSubRip Format = "subrip"
	FormatWebVTT Format = "webvtt"
)

// Capability names one independently optional codec operation.
type Capability string

const (
	CapabilityDecode Capability = "decode"
	CapabilityRender Capability = "render"
)

// DecodeOptions controls format-native model derivation.
type DecodeOptions struct {
	Encoding                string
	DisableSpeakerDetection bool
}

// RenderOptions controls deterministic model-driven serialization.
type RenderOptions struct {
	Strict bool
}

// RenderResult contains rendered bytes and ordered non-fatal diagnostics.
type RenderResult struct {
	Bytes       []byte
	Diagnostics []model.Diagnostic
}

// Decoder derives one Cue JSON document from an already captured source.
type Decoder func(context.Context, source.Captured, DecodeOptions) (model.Document, error)

// Renderer serializes a validated Cue JSON document without restoring source bytes.
type Renderer func(context.Context, model.Document, RenderOptions) (RenderResult, error)

// Detector returns bounded content evidence for one format.
type Detector func([]byte) Evidence

// Registration describes one canonical format and its optional native capabilities.
type Registration struct {
	Format     Format
	Aliases    []string
	Extensions []string
	Detect     Detector
	Decode     Decoder
	Render     Renderer
}

// Registry is an immutable deterministic collection of native format registrations.
type Registry struct {
	registrations []Registration
	selectors     map[string]int
	extensions    map[string]int
}

// NewRegistry validates and freezes the supplied registrations.
func NewRegistry(registrations ...Registration) (*Registry, error) {
	prepared := make([]Registration, len(registrations))
	for index := range registrations {
		prepared[index] = cloneRegistration(registrations[index])
		prepared[index].Format = Format(normalizeToken(string(prepared[index].Format)))
		if prepared[index].Format == "" || prepared[index].Format == FormatAuto {
			return nil, fmt.Errorf("codec registration %d has invalid canonical format %q", index, registrations[index].Format)
		}
		for aliasIndex := range prepared[index].Aliases {
			prepared[index].Aliases[aliasIndex] = normalizeToken(prepared[index].Aliases[aliasIndex])
			if prepared[index].Aliases[aliasIndex] == "" || prepared[index].Aliases[aliasIndex] == string(FormatAuto) {
				return nil, fmt.Errorf("codec %q has invalid alias %q", prepared[index].Format, registrations[index].Aliases[aliasIndex])
			}
		}
		for extensionIndex := range prepared[index].Extensions {
			prepared[index].Extensions[extensionIndex] = normalizeExtension(prepared[index].Extensions[extensionIndex])
			if prepared[index].Extensions[extensionIndex] == "" {
				return nil, fmt.Errorf("codec %q has an empty extension", prepared[index].Format)
			}
		}
	}
	sort.Slice(prepared, func(left, right int) bool { return prepared[left].Format < prepared[right].Format })

	registry := &Registry{registrations: prepared, selectors: make(map[string]int), extensions: make(map[string]int)}
	for index := range prepared {
		selectors := append([]string{string(prepared[index].Format)}, prepared[index].Aliases...)
		for _, selector := range selectors {
			if prior, exists := registry.selectors[selector]; exists {
				return nil, fmt.Errorf("codec selector %q is claimed by both %q and %q", selector, prepared[prior].Format, prepared[index].Format)
			}
			registry.selectors[selector] = index
		}
		for _, extension := range prepared[index].Extensions {
			if prior, exists := registry.extensions[extension]; exists {
				return nil, fmt.Errorf("codec extension %q is claimed by both %q and %q", extension, prepared[prior].Format, prepared[index].Format)
			}
			registry.extensions[extension] = index
		}
	}
	return registry, nil
}

// Lookup resolves a canonical format or registered alias.
func (registry *Registry) Lookup(selector string) (Registration, bool) {
	if registry == nil {
		return Registration{}, false
	}
	index, exists := registry.selectors[normalizeToken(selector)]
	if !exists {
		return Registration{}, false
	}
	return cloneRegistration(registry.registrations[index]), true
}

// RequireDecoder returns the installed decoder or a typed capability failure.
func (registry *Registry) RequireDecoder(format Format) (Decoder, error) {
	registration, exists := registry.Lookup(string(format))
	if !exists {
		return nil, &UnknownFormatError{Selector: string(format)}
	}
	if registration.Decode == nil {
		return nil, &MissingCapabilityError{Format: registration.Format, Capability: CapabilityDecode}
	}
	return registration.Decode, nil
}

// RequireRenderer returns the installed renderer or a typed capability failure.
func (registry *Registry) RequireRenderer(format Format) (Renderer, error) {
	registration, exists := registry.Lookup(string(format))
	if !exists {
		return nil, &UnknownFormatError{Selector: string(format)}
	}
	if registration.Render == nil {
		return nil, &MissingCapabilityError{Format: registration.Format, Capability: CapabilityRender}
	}
	return registration.Render, nil
}

// MissingCapabilityError distinguishes a recognized format from an installed operation.
type MissingCapabilityError struct {
	Format     Format
	Capability Capability
}

func (err *MissingCapabilityError) Error() string {
	return fmt.Sprintf("format %q is recognized but its native %s capability is unavailable", err.Format, err.Capability)
}

// UnknownFormatError identifies an unrecognized explicit selector or source.
type UnknownFormatError struct {
	Selector string
}

func (err *UnknownFormatError) Error() string {
	if err.Selector == "" {
		return "source format is unknown"
	}
	return fmt.Sprintf("format %q is unknown", err.Selector)
}

func cloneRegistration(registration Registration) Registration {
	registration.Aliases = append([]string(nil), registration.Aliases...)
	registration.Extensions = append([]string(nil), registration.Extensions...)
	return registration
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func normalizeExtension(value string) string {
	return strings.TrimPrefix(normalizeToken(value), ".")
}
