package xml

import (
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
)

var (
	country = specs.Template{
		Message: specs.Message{
			"iso2Code": {
				Name:  "iso2Code",
				Path:  "country.iso2Code",
				Label: labels.Optional,
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.String,
					},
				},
			},
			"name": {
				Name:  "name",
				Path:  "country.name",
				Label: labels.Optional,
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.String,
					},
				},
			},
			"latitude": {
				Name:  "iso2Code",
				Path:  "country.latitude",
				Label: labels.Optional,
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.Float,
					},
				},
			},
			"longitude": {
				Name:  "name",
				Path:  "country.longitude",
				Label: labels.Optional,
				Template: specs.Template{
					Scalar: &specs.Scalar{
						Type: types.Float,
					},
				},
			},
		},
	}

	SchemaEnum = &specs.ParameterMap{
		Property: &specs.Property{
			Name:  "status",
			Path:  "status",
			Label: labels.Required,
			Template: specs.Template{
				Enum: enum,
			},
		},
	}

	SchemaScalar = &specs.ParameterMap{
		Property: &specs.Property{
			Name:  "integer",
			Path:  "integer",
			Label: labels.Required,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type: types.Int32,
				},
			},
		},
	}

	SchemaArray = &specs.Property{
		Name:  "country",
		Path:  "country",
		Label: labels.Optional,
		Template: specs.Template{
			Repeated: specs.Repeated{
				country,
			},
		},
	}

	SchemaObject = &specs.ParameterMap{
		Property: &specs.Property{
			Name:  "root",
			Label: labels.Optional,
			Template: specs.Template{
				Message: specs.Message{
					"bad_label": {
						Position: 1,
						Name:     "bad_label",
						Path:     "bad_label",
						Label:    "unknown",
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
					"numeric": {
						Position: 3,
						Name:     "numeric",
						Path:     "numeric",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "numeric",
							},
							Scalar: &specs.Scalar{
								Type: types.Int32,
							},
						},
					},
					"message": {
						Position: 4,
						Name:     "message",
						Path:     "message",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "message",
							},
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
					"another_message": {
						Position: 5,
						Name:     "another_message",
						Path:     "another_message",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "another_message",
							},
							Scalar: &specs.Scalar{
								Type: types.String,
							},
						},
					},
					"status": {
						Position: 6,
						Name:     "status",
						Path:     "status",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "status",
							},
							Enum: enum,
						},
					},
					"another_status": {
						Position: 7,
						Name:     "another_status",
						Path:     "another_status",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "another_status",
							},
							Enum: enum,
						},
					},
					"nested": {
						Position: 8,
						Name:     "nested",
						Path:     "nested",
						Label:    labels.Optional,
						Template: specs.Template{
							Message: specs.Message{
								"first": {
									Position: 1,
									Name:     "first",
									Path:     "nested.first",
									Label:    labels.Optional,
									Template: specs.Template{
										Reference: &specs.PropertyReference{
											Resource: template.InputResource,
											Path:     "nested.first",
										},
										Scalar: &specs.Scalar{
											Type: types.String,
										},
									},
								},
								"second": {
									Position: 2,
									Name:     "second",
									Path:     "nested.second",
									Label:    labels.Optional,
									Template: specs.Template{
										Reference: &specs.PropertyReference{
											Resource: template.InputResource,
											Path:     "nested.second",
										},
										Scalar: &specs.Scalar{
											Type: types.String,
										},
									},
								},
							},
						},
					},
					"repeating_string": {
						Position: 9,
						Name:     "repeating_string",
						Path:     "repeating_string",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "repeating_string",
							},
							Repeated: specs.Repeated{
								{
									Scalar: &specs.Scalar{
										Type: types.String,
									},
								},
							},
						},
					},
					"repeating_enum": {
						Position: 10,
						Name:     "repeating_enum",
						Path:     "repeating_enum",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "repeating_enum",
							},
							Repeated: specs.Repeated{
								{
									Enum: enum,
								},
							},
						},
					},
					"repeating_numeric": {
						Position: 11,
						Name:     "repeating_numeric",
						Path:     "repeating_numeric",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "repeating_numeric",
							},
							Scalar: &specs.Scalar{
								Type: types.Int32,
							},
						},
					},
					"repeating": {
						Position: 12,
						Name:     "repeating",
						Path:     "repeating",
						Label:    labels.Optional,
						Template: specs.Template{
							Reference: &specs.PropertyReference{
								Resource: template.InputResource,
								Path:     "repeating",
							},
							Repeated: specs.Repeated{
								{
									Message: specs.Message{
										"value": {
											Position: 1,
											Name:     "value",
											Path:     "repeating.value",
											Label:    labels.Optional,
											Template: specs.Template{
												Reference: &specs.PropertyReference{
													Resource: template.InputResource,
													Path:     "repeating.value",
												},
												Scalar: &specs.Scalar{
													Type: types.String,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	enum = &specs.Enum{
		Keys: map[string]*specs.EnumValue{
			"UNKNOWN": {
				Key:      "UNKNOWN",
				Position: 0,
			},
			"PENDING": {
				Key:      "PENDING",
				Position: 1,
			},
		},
		Positions: map[int32]*specs.EnumValue{
			0: {
				Key:      "UNKNOWN",
				Position: 0,
			},
			1: {
				Key:      "PENDING",
				Position: 1,
			},
		},
	}
)
