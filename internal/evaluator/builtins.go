package evaluator

import (
	"fmt"

	"github.com/myselfBZ/interpreter/internal/object"
)





var builtIns = map[string]*object.BuiltIn{
	"puts":&object.BuiltIn{
		Fn: PutsBuiltin,
	},
}

func PutsBuiltin(objs ...object.Object) object.Object {
	for _, obj := range objs {
		switch t := obj.(type) {
		case *object.Integer:
			fmt.Print(t.Value)
		case *object.Boolean:
			fmt.Print(t.Value)
		case *object.Null:
			fmt.Print("NULL")
		default:
			return newError("non-printable object: %s", obj.Type())
		}
		fmt.Print(" ")
	}
	fmt.Println()
	return NULL
}
