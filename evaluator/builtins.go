package evaluator

import "github.com/hollykbuck/muskmelon/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			return NULL
		},
	},
}
