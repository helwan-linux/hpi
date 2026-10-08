package packageinfo

import "strings"

type Dependency struct {
	RawName  string
	Name     string
	Operator string
	Version  string
}

func ParseDependency(raw string) Dependency {
	raw = strings.TrimSpace(raw)

	dependency := Dependency{
		RawName: raw,
	}

	if raw == "" {
		return dependency
	}

	operators := []string{
		">=",
		"<=",
		"=",
		">",
		"<",
	}

	for _, operator := range operators {
		index := strings.Index(raw, operator)

		if index < 0 {
			continue
		}

		dependency.Name = strings.TrimSpace(
			raw[:index],
		)

		dependency.Operator = operator

		dependency.Version = strings.TrimSpace(
			raw[index+len(operator):],
		)

		return dependency
	}

	dependency.Name = raw

	return dependency
}

func (d Dependency) DisplayName() string {
	if d.Name == "" {
		return d.RawName
	}

	if d.Operator == "" {
		return d.Name
	}

	return d.Name + " " + d.Operator + " " + d.Version
}

func (p *Package) Dependencies() []Dependency {
	if p == nil || p.Metadata == nil {
		return nil
	}

	result := make(
		[]Dependency,
		0,
		len(p.Metadata.Depends),
	)

	for _, raw := range p.Metadata.Depends {
		dependency := ParseDependency(raw)

		if dependency.Name == "" {
			continue
		}

		result = append(
			result,
			dependency,
		)
	}

	return result
}

func (p *Package) OptionalDependencies() []Dependency {
	if p == nil || p.Metadata == nil {
		return nil
	}

	result := make(
		[]Dependency,
		0,
		len(p.Metadata.OptDepends),
	)

	for _, raw := range p.Metadata.OptDepends {
		dependency := ParseDependency(raw)

		if dependency.Name == "" {
			continue
		}

		result = append(
			result,
			dependency,
		)
	}

	return result
}
