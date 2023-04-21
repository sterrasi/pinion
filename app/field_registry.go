package app

// ConfigurationBuilderFn function used by an Application to build its configuration field definitions
type ConfigurationBuilderFn func(registry *FieldRegistry) Error

// FieldRegistry is used to create and store Fields
type FieldRegistry struct {
	fields map[string]*Field
}

// CreateIntField Builds a new Integer type FieldBuilder
func (r *FieldRegistry) CreateIntField(name string) *FieldBuilder[int] {
	return &FieldBuilder[int]{
		name:      name,
		registry:  r,
		valueType: Int}
}

// CreateUintField Builds a new uint type FieldBuilder
func (r *FieldRegistry) CreateUintField(name string) *FieldBuilder[uint] {
	return &FieldBuilder[uint]{
		name:      name,
		registry:  r,
		valueType: Uint}
}

// CreateFloatField Builds a new float64 type FieldBuilder
func (r *FieldRegistry) CreateFloatField(name string) *FieldBuilder[float64] {
	return &FieldBuilder[float64]{
		name:      name,
		registry:  r,
		valueType: Float}
}

// CreateBooleanField Builds a new bool type FieldBuilder
func (r *FieldRegistry) CreateBooleanField(name string) *FieldBuilder[bool] {
	return &FieldBuilder[bool]{
		name:      name,
		registry:  r,
		valueType: Bool}
}

// CreateStringField Builds a new string type FieldBuilder
func (r *FieldRegistry) CreateStringField(name string) *FieldBuilder[string] {
	return &FieldBuilder[string]{
		name:      name,
		registry:  r,
		valueType: String}
}
