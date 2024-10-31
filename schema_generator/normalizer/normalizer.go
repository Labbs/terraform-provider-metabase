package normalizer

import (
	"strings"
)

type SchemaNameMapper struct {
	originalToNormalized map[string]string
	normalizedToOriginal map[string]string
}

func NewSchemaNameMapper() *SchemaNameMapper {
	return &SchemaNameMapper{
		originalToNormalized: make(map[string]string),
		normalizedToOriginal: make(map[string]string),
	}
}

func (m *SchemaNameMapper) Add(original, normalized string) {
	m.originalToNormalized[original] = normalized
	m.normalizedToOriginal[normalized] = original
}

func (m *SchemaNameMapper) GetNormalized(original string) string {
	if normalized, exists := m.originalToNormalized[original]; exists {
		return normalized
	}
	return original
}

// Function to decode special characters in a reference
func decodeRefString(ref string) string {
	// Replace ~1 with / (according to JSON Pointer specification)
	ref = strings.ReplaceAll(ref, "~1", "/")
	// Replace ~0 with ~ (according to JSON Pointer specification)
	ref = strings.ReplaceAll(ref, "~0", "~")
	return ref
}

// Function to normalize a schema name
func normalizeSchemaName(name string) string {
	// First decode special characters
	name = decodeRefString(name)

	// Replace all special characters with underscores
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, "~", "_")

	// Clean up multiple underscores
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}

	// Remove underscores at the beginning and end
	name = strings.Trim(name, "_")

	return name
}

// Function to convert a reference using the mapper
func normalizeSchemaRef(ref string, mapper *SchemaNameMapper) string {
	if strings.HasPrefix(ref, "#/components/schemas/") {
		schemaName := strings.TrimPrefix(ref, "#/components/schemas/")
		decodedName := decodeRefString(schemaName)

		// First check in the mapper
		if normalized := mapper.GetNormalized(decodedName); normalized != decodedName {
			return "#/components/schemas/" + normalized
		}

		// Otherwise, normalize the name
		normalizedName := normalizeSchemaName(schemaName)
		return "#/components/schemas/" + normalizedName
	}
	return ref
}

// Recursive function to normalize all references in an object
func NormalizeRefs(obj interface{}, mapper *SchemaNameMapper) interface{} {
	switch v := obj.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			if key == "$ref" {
				if refStr, ok := value.(string); ok {
					result[key] = normalizeSchemaRef(refStr, mapper)
				} else {
					result[key] = value
				}
			} else {
				result[key] = NormalizeRefs(value, mapper)
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = NormalizeRefs(value, mapper)
		}
		return result
	default:
		return obj
	}
}

// Function to normalize schemas and create mapping
func NormalizeSchemas(schemas map[string]interface{}) (map[string]interface{}, *SchemaNameMapper) {
	mapper := NewSchemaNameMapper()
	normalizedSchemas := make(map[string]interface{})

	// First pass: create the name mapping
	for name := range schemas {
		decodedName := decodeRefString(name)
		normalized := normalizeSchemaName(decodedName)
		mapper.Add(name, normalized)
		mapper.Add(decodedName, normalized)
	}

	// Second pass: normalize schemas with mapping
	for name, schema := range schemas {
		normalized := mapper.GetNormalized(name)
		normalizedSchema := NormalizeRefs(schema, mapper)
		normalizedSchemas[normalized] = normalizedSchema
	}

	return normalizedSchemas, mapper
}
