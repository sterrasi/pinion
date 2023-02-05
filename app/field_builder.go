package app

type FieldType interface {
	int | uint | float64 | bool | string
}

// FieldBuilder builder for a Field
type FieldBuilder[T FieldType] struct {
	shortDescription  string
	longDescription   string
	name              string
	argName           string
	envVar            string
	configSectionName string
	configFieldName   string
	defaultValue      T
	required          bool
	valueType         ValueType
	config            *Configuration
}

func (b *FieldBuilder[T]) Descriptions(shortDesc string, longDesc string) *FieldBuilder[T] {
	b.shortDescription = shortDesc
	b.longDescription = longDesc
	return b
}

func (b *FieldBuilder[T]) ShortDesc(shortDesc string) *FieldBuilder[T] {
	b.shortDescription = shortDesc
	return b
}

func (b *FieldBuilder[T]) ArgName(argName string) *FieldBuilder[T] {
	b.argName = argName
	return b
}
func (b *FieldBuilder[T]) EnvVar(envVar string) *FieldBuilder[T] {
	b.envVar = envVar
	return b
}
func (b *FieldBuilder[T]) ConfigName(configSectionName string, configFieldName string) *FieldBuilder[T] {
	b.configSectionName = configSectionName
	b.configFieldName = configFieldName
	return b
}
func (b *FieldBuilder[T]) Default(defaultValue T) *FieldBuilder[T] {
	b.defaultValue = defaultValue
	return b
}
func (b *FieldBuilder[T]) Required() *FieldBuilder[T] {
	b.required = true
	return b
}
func (b *FieldBuilder[T]) Register() *Field {
	f := &Field{
		ShortDescription:  b.shortDescription,
		LongDescription:   b.longDescription,
		Name:              b.name,
		ArgName:           b.argName,
		EnvVar:            b.envVar,
		ConfigSectionName: b.configSectionName,
		ConfigFieldName:   b.configFieldName,
		DefaultValue:      b.defaultValue,
		Required:          b.required,
		Type:              b.valueType,
	}
	if b.config.fields == nil {
		b.config.fields = make(map[string]*Field)
	}
	b.config.fields[f.Name] = f
	return f
}
