package eval

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"time"
)

var defaultFuncMap = map[string]func(...ref.Val) ref.Val{
	"now": func(val ...ref.Val) ref.Val {
		return types.Int(time.Now().Unix())
	},
	"sha1": func(vals ...ref.Val) ref.Val {
		hash := sha1.New()
		for _, val := range vals {
			switch val.Type() {
			case types.StringType:
				_, err := hash.Write([]byte(val.Value().(string)))
				if err != nil {
					return types.NewErr("failed to sha1 hash: %s", err.Error())
				}
			}
		}

		return types.String(hex.EncodeToString(hash.Sum(nil)))
	},
}
