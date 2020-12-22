package eval

import (
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"time"
)

var defaultFuncMap = map[string]func(...ref.Val) ref.Val{
	"now": func(val ...ref.Val) ref.Val {
		return types.Int(time.Now().Unix())
	},
}
