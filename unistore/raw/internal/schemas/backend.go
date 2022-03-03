package schemas

import (
	"fmt"
	"strings"

	"unistore/schemas/common"
)

// MarshalDescriptor is a convenience formatter for marshaling a common.Node as a string structured
// as nested function invocations.
func MarshalDescriptor(node *common.Node) string {
	var parameters []string

	childrenConstants := make(map[string]string)
	childrenNodes := make(map[string]*common.Node)

	for param, value := range node.Children {
		switch c := value.Child.(type) {
		case *common.Node_Value_Value:
			childrenConstants[param] = c.Value
		case *common.Node_Value_Node:
			childrenNodes[param] = c.Node
		}
	}

	// Node parameters (recursive definitions) are always represented at the end of the
	// parameters list.

	for param, child := range childrenConstants {
		parameters = append(parameters, fmt.Sprintf("%s: %s", param, child))
	}

	for param, child := range childrenNodes {
		parameters = append(parameters, fmt.Sprintf("%s: %s", param, MarshalDescriptor(child)))
	}

	return fmt.Sprintf("%s(%s)", node.Name, strings.Join(parameters, ", "))
}
