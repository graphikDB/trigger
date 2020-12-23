package trigger

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/uuid"
	"time"
)

func errNoParams(fn string) ref.Val {
	return types.NewErr("eval: [%s] expected at least one paramater", fn)
}

func errFunction(fn string, msg string) ref.Val {
	return types.NewErr("eval: [%s] %s", fn, msg)
}

func errBadType(fn string, expected string) ref.Val {
	return types.NewErr("eval: function %s expected type %s", fn, expected)
}

var defaultFuncMap = map[string]func(...ref.Val) ref.Val{
	"now": func(val ...ref.Val) ref.Val {
		return types.Int(time.Now().Unix())
	},
	"uuid": func(val ...ref.Val) ref.Val {
		return types.String(uuid.New().String())
	},
	"sha1": func(vals ...ref.Val) ref.Val {
		if len(vals) == 0 {
			return errNoParams("sha1")
		}
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
	"sha256": func(vals ...ref.Val) ref.Val {
		if len(vals) == 0 {
			return errNoParams("sha256")
		}

		hash := sha256.New()
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
	"base64Encode": func(vals ...ref.Val) ref.Val {
		if len(vals) == 0 {
			return errNoParams("base64Encode")
		}
		if _, ok := vals[0].Value().(string); !ok {
			return errBadType("base64Encode", "string")
		}
		return types.String(base64.StdEncoding.EncodeToString([]byte(vals[0].Value().(string))))
	},
	"base64Decode": func(vals ...ref.Val) ref.Val {
		if len(vals) == 0 {
			return errNoParams("base64Decode")
		}
		if _, ok := vals[0].Value().(string); !ok {
			return errBadType("base64Decode", "string")
		}
		decoded, err := base64.StdEncoding.DecodeString(vals[0].Value().(string))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.String(decoded)
	},
	"includes": func(vals ...ref.Val) ref.Val {
		if len(vals) != 2 {
			return errFunction("includes", "expected two params")
		}
		if vals[0].Type() != types.ListType {
			return errFunction("includes", "expected first param to be list")
		}
		target := vals[1].Value()
		for _, val := range vals[0].Value().([]interface{}) {
			if val == target {
				return types.Bool(true)
			}
		}
		return types.Bool(false)
	},
}
