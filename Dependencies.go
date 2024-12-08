package Binance

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.Config{
	EscapeHTML:              false, // Avoids escaping HTML (faster for most JSON use cases)
	SortMapKeys:             false, // Disables map key sorting (reduces overhead)
	MarshalFloatWith6Digits: true,  // Optimizes float marshaling
	// ValidateJsonRawMessage:  false, // Skips validation of raw messages
	// OnlyTaggedField: true, // Ignores untagged struct fields (if you don't need them)
	// ObjectFieldMustBeSimpleString: true,  // Optimizes string field parsing
	// CaseSensitive:                 false, // Allows case-insensitive matching
}.Froze()
