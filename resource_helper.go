package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// resourcePropertyHelper provides commonly-used functionality for working with Terraform's schema.ResourceData.
type resourcePropertyHelper struct {
	data *schema.ResourceData
}

func propertyHelper(data *schema.ResourceData) resourcePropertyHelper {
	return resourcePropertyHelper{data}
}

func (helper resourcePropertyHelper) GetStringList(key string) (elements []string) {
	value, ok := helper.data.GetOk(key)
	if !ok {
		return
	}

	untypedElements := value.([]interface{})
	elements = make([]string, len(untypedElements))
	for index, untypedElement := range untypedElements {
		elements[index] = untypedElement.(string)
	}

	return
}

func (helper resourcePropertyHelper) SetStringList(key string, elements []string) {
	untypedElements := make([]interface{}, len(elements))
	for index, element := range elements {
		var untypedElement interface{}
		untypedElement = element

		untypedElements[index] = untypedElement
	}

	helper.data.Set(key, untypedElements)
}

func (helper resourcePropertyHelper) GetOptionalString(key string, allowEmpty bool) *string {
	value := helper.data.Get(key)
	switch typedValue := value.(type) {
	case string:
		if len(typedValue) > 0 || allowEmpty {
			return &typedValue
		}
	}

	return nil
}

func (helper resourcePropertyHelper) GetOptionalInt(key string, allowZero bool) *int {
	value := helper.data.Get(key)
	switch typedValue := value.(type) {
	case int:
		if typedValue != 0 || allowZero {
			return &typedValue
		}
	}

	return nil
}

func (helper resourcePropertyHelper) GetOptionalBool(key string) *bool {
	value := helper.data.Get(key)
	switch typedValue := value.(type) {
	case bool:
		return &typedValue
	default:
		return nil
	}
}
