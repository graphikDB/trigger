package trigger

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
	"golang.org/x/crypto/sha3"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"html/template"
	"reflect"
	"strings"
	"time"
)

var strMap = decls.NewMapType(decls.String, decls.Any)

type Function struct {
	decl     *expr.Decl
	overload *functions.Overload
}

func NewFunction(decl *expr.Decl, overload *functions.Overload) *Function {
	return &Function{
		decl:     decl,
		overload: overload,
	}
}

type FuncMap map[string]*Function

var Functions = FuncMap{
	"now": {
		decl: decls.NewFunction("now",
			decls.NewOverload(
				"now",
				[]*expr.Type{},
				decls.Int,
			),
		),
		overload: &functions.Overload{
			Operator: "now",
			Function: defaultFuncMap["now"],
		},
	},
	"uuid": {
		decl: decls.NewFunction("uuid",
			decls.NewOverload(
				"uuid",
				[]*expr.Type{},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "uuid",
			Function: defaultFuncMap["uuid"],
		},
	},
	"sha1": {
		decl: decls.NewFunction("sha1",
			decls.NewOverload(
				"sha1_string",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "sha1_string",
			Function: defaultFuncMap["sha1"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["sha1"](value)
			},
		},
	},
	"sha256": {
		decl: decls.NewFunction("sha256",
			decls.NewOverload(
				"sha256_string",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "sha256_string",
			Function: defaultFuncMap["sha256"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["sha256"](value)
			},
		},
	},
	"sha3": {
		decl: decls.NewFunction("sha3",
			decls.NewOverload(
				"sha3_string",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "sha3_string",
			Function: defaultFuncMap["sha3"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["sha3"](value)
			},
		},
	},
	"base64Encode": {
		decl: decls.NewFunction("base64Encode",
			decls.NewOverload(
				"base64Encode_string",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "base64Encode_string",
			Function: defaultFuncMap["base64Encode"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["base64Encode"](value)
			},
		},
	},
	"base64Decode": {
		decl: decls.NewFunction("base64Decode",
			decls.NewOverload(
				"base64Decode_string",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "base64Decode_string",
			Function: defaultFuncMap["base64Decode"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["base64Decode"](value)
			},
		},
	},
	"jsonEncode": {
		decl: decls.NewFunction("jsonEncode",
			decls.NewOverload(
				"jsonEncode_string",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "jsonEncode_string",
			Function: defaultFuncMap["jsonEncode"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["jsonEncode"](value)
			},
		},
	},
	"jsonDecode": {
		decl: decls.NewFunction("jsonDecode",
			decls.NewOverload(
				"jsonDecode_string",
				[]*expr.Type{decls.String},
				strMap,
			),
		),
		overload: &functions.Overload{
			Operator: "jsonDecode_string",
			Function: defaultFuncMap["jsonDecode"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["jsonDecode"](value)
			},
		},
	},
	"includes": {
		decl: decls.NewFunction("includes",
			decls.NewOverload(
				"includes_list",
				[]*expr.Type{decls.NewListType(decls.Any), decls.Any},
				decls.Bool,
			),
		),
		overload: &functions.Overload{
			Operator: "includes_list",
			Function: defaultFuncMap["includes"],
			Binary: func(value ref.Val, value2 ref.Val) ref.Val {
				return defaultFuncMap["includes"](value, value2)
			},
		},
	},
	"replace": {
		decl: decls.NewFunction("replace",
			decls.NewOverload(
				"replace",
				[]*expr.Type{decls.String, decls.String, decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "replace",
			Function: defaultFuncMap["replace"],
			//Binary: func(value ref.Val, value2 ref.Val) ref.Val {
			//	return defaultFuncMap["replace"](value, value2)
			//},
		},
	},
	"join": {
		decl: decls.NewFunction("join",
			decls.NewOverload(
				"join",
				[]*expr.Type{decls.NewListType(decls.String), decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "join",
			Function: defaultFuncMap["join"],
			Binary: func(value ref.Val, value2 ref.Val) ref.Val {
				return defaultFuncMap["join"](value, value2)
			},
		},
	},
	"titleCase": {
		decl: decls.NewFunction("titleCase",
			decls.NewOverload(
				"titleCase",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "titleCase",
			Function: defaultFuncMap["titleCase"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["titleCase"](value)
			},
		},
	},
	"lowerCase": {
		decl: decls.NewFunction("lowerCase",
			decls.NewOverload(
				"lowerCase",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "lowerCase",
			Function: defaultFuncMap["lowerCase"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["lowerCase"](value)
			},
		},
	},
	"upperCase": {
		decl: decls.NewFunction("upperCase",
			decls.NewOverload(
				"upperCase",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "upperCase",
			Function: defaultFuncMap["upperCase"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["upperCase"](value)
			},
		},
	},
	"trimSpace": {
		decl: decls.NewFunction("trimSpace",
			decls.NewOverload(
				"trimSpace",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "trimSpace",
			Function: defaultFuncMap["trimSpace"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["trimSpace"](value)
			},
		},
	},
	"trimPrefix": {
		decl: decls.NewFunction("trimPrefix",
			decls.NewOverload(
				"trimPrefix",
				[]*expr.Type{decls.String, decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "trimPrefix",
			Function: defaultFuncMap["trimPrefix"],
			Binary: func(value ref.Val, value2 ref.Val) ref.Val {
				return defaultFuncMap["trimPrefix"](value, value2)
			},
		},
	},
	"trimSuffix": {
		decl: decls.NewFunction("trimSuffix",
			decls.NewOverload(
				"trimSuffix",
				[]*expr.Type{decls.String, decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "trimSuffix",
			Function: defaultFuncMap["trimSuffix"],
			Binary: func(value ref.Val, value2 ref.Val) ref.Val {
				return defaultFuncMap["trimSuffix"](value)
			},
		},
	},
	"split": {
		decl: decls.NewFunction("split",
			decls.NewOverload(
				"split",
				[]*expr.Type{decls.String, decls.String},
				decls.NewListType(decls.String),
			),
		),
		overload: &functions.Overload{
			Operator: "split",
			Function: defaultFuncMap["split"],
			Binary: func(value ref.Val, value2 ref.Val) ref.Val {
				return defaultFuncMap["split"](value, value2)
			},
		},
	},
	"geoDistance": {
		decl: decls.NewFunction("geoDistance",
			decls.NewOverload(
				"geoDistance",
				[]*expr.Type{decls.NewListType(decls.Double), decls.NewListType(decls.Double)},
				decls.Double,
			),
		),
		overload: &functions.Overload{
			Operator: "geoDistance",
			Function: defaultFuncMap["geoDistance"],
			Binary: func(value ref.Val, value2 ref.Val) ref.Val {
				return defaultFuncMap["geoDistance"](value, value2)
			},
		},
	},
	"render": {
		decl: decls.NewFunction("render",
			decls.NewOverload(
				"render",
				[]*expr.Type{decls.String, strMap},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "render",
			Function: defaultFuncMap["render"],
			Binary: func(value ref.Val, value2 ref.Val) ref.Val {
				return defaultFuncMap["render"](value, value2)
			},
		},
	},
	"parseClaims": {
		decl: decls.NewFunction("parseClaims",
			decls.NewOverload(
				"parseClaims",
				[]*expr.Type{decls.String},
				strMap,
			),
		),
		overload: &functions.Overload{
			Operator: "parseClaims",
			Function: defaultFuncMap["parseClaims"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["parseClaims"](value)
			},
		},
	},
	"typeOf": {
		decl: decls.NewFunction("typeOf",
			decls.NewOverload(
				"typeOf",
				[]*expr.Type{decls.Any},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "typeOf",
			Function: defaultFuncMap["typeOf"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["typeOf"](value)
			},
		},
	},
}

func errFunction(fn string, msg string) ref.Val {
	return types.NewErr("trigger: [%s] %s", fn, msg)
}

var defaultFuncMap = map[string]func(...ref.Val) ref.Val{
	"now": func(val ...ref.Val) ref.Val {
		return types.Int(time.Now().Unix())
	},
	"uuid": func(val ...ref.Val) ref.Val {
		return types.String(uuid.New().String())
	},
	"sha1": func(vals ...ref.Val) ref.Val {
		hash := sha1.New()
		for _, val := range vals {
			switch val.Type() {
			case types.StringType:
				_, err := hash.Write([]byte(cast.ToString(val.Value())))
				if err != nil {
					return types.NewErr("failed to sha1 hash: %s", err.Error())
				}
			}
		}
		return types.String(hex.EncodeToString(hash.Sum(nil)))
	},
	"sha256": func(vals ...ref.Val) ref.Val {
		hash := sha256.New()
		for _, val := range vals {
			switch val.Type() {
			case types.StringType:
				_, err := hash.Write([]byte(cast.ToString(val.Value())))
				if err != nil {
					return types.NewErr("failed to sha1 hash: %s", err.Error())
				}
			}
		}
		return types.String(hex.EncodeToString(hash.Sum(nil)))
	},
	"sha3": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("sha3", "expected one param")
		}

		hash := sha3.New512()
		for _, val := range vals {
			switch val.Type() {
			case types.StringType:
				_, err := hash.Write([]byte(cast.ToString(val.Value())))
				if err != nil {
					return types.NewErr("failed to sha1 hash: %s", err.Error())
				}
			}
		}
		return types.String(hex.EncodeToString(hash.Sum(nil)))
	},
	"base64Encode": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("base64Encode", "expected one param")
		}
		return types.String(base64.StdEncoding.EncodeToString([]byte(cast.ToString(vals[0].Value()))))
	},
	"base64Decode": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("base64Decode", "expected one param")
		}
		decoded, err := base64.StdEncoding.DecodeString(cast.ToString(vals[0].Value()))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.String(decoded)
	},
	"jsonEncode": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("jsonEncode", "expected one param")
		}
		bits, _ := json.Marshal(vals[0].Value())
		return types.String(bits)
	},
	"jsonDecode": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("jsonDecode", "expected one param")
		}
		data := map[string]interface{}{}
		json.Unmarshal([]byte(cast.ToString(vals[0].Value())), &data)
		return types.NewStringInterfaceMap(types.DefaultTypeAdapter, data)
	},
	"includes": func(vals ...ref.Val) ref.Val {
		if len(vals) != 2 {
			return errFunction("includes", "expected two params")
		}
		if vals[0].Type() != types.ListType {
			return errFunction("includes", "expected first param to be list")
		}

		return types.Bool(funk.Contains(vals[0].Value(), vals[1].Value()))
	},
	"replaceN": func(vals ...ref.Val) ref.Val {
		if len(vals) != 4 {
			return errFunction("replace", "expected 4 params")
		}
		return types.String(strings.Replace(
			cast.ToString(vals[0].Value()),
			cast.ToString(vals[1].Value()),
			cast.ToString(vals[2].Value()),
			cast.ToInt(vals[3].Value()),
		))
	},
	"replace": func(vals ...ref.Val) ref.Val {
		if len(vals) != 3 {
			return errFunction("replace", "expected three params")
		}
		return types.String(strings.Replace(cast.ToString(vals[0].Value()), cast.ToString(vals[1].Value()), cast.ToString(vals[2].Value()), -1))
	},
	"join": func(vals ...ref.Val) ref.Val {
		if len(vals) != 2 {
			return errFunction("join", "expected two params")
		}
		if vals[0].Type() != types.ListType {
			return errFunction("join", "expected first param to be list")
		}
		elems, err := cast.ToStringSliceE(vals[0].Value())
		if err != nil {
			return errFunction("join", err.Error())
		}
		sep := cast.ToString(vals[1].Value())

		return types.String(strings.Join(elems, sep))
	},
	"titleCase": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("titleCase", "expected one params")
		}

		return types.String(strings.Title(cast.ToString(vals[0].Value())))
	},
	"lowerCase": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("lowerCase", "expected one params")
		}

		return types.String(strings.ToLower(cast.ToString(vals[0].Value())))
	},
	"upperCase": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("upperCase", "expected one params")
		}

		return types.String(strings.ToUpper(cast.ToString(vals[0].Value())))
	},
	"trimSpace": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("trimSpace", "expected one params")
		}
		return types.String(strings.TrimSpace(cast.ToString(vals[0].Value())))
	},
	"trimPrefix": func(vals ...ref.Val) ref.Val {
		if len(vals) != 2 {
			return errFunction("trimPrefix", "expected two params")
		}
		return types.String(strings.TrimPrefix(cast.ToString(vals[0].Value()), cast.ToString(vals[1].Value())))
	},
	"trimSuffix": func(vals ...ref.Val) ref.Val {
		if len(vals) != 2 {
			return errFunction("trimSuffix", "expected two params")
		}
		return types.String(strings.TrimSuffix(cast.ToString(vals[0].Value()), cast.ToString(vals[1].Value())))
	},
	"split": func(vals ...ref.Val) ref.Val {
		if len(vals) != 2 {
			return errFunction("split", "expected two params")
		}
		return types.NewStringList(types.DefaultTypeAdapter, strings.Split(cast.ToString(vals[0].Value()), cast.ToString(vals[1].Value())))
	},
	"geoDistance": func(vals ...ref.Val) ref.Val {
		if len(vals) != 2 {
			return errFunction("geoDistance", "expected two params")
		}
		if vals[0].Type() != types.ListType {
			return errFunction("geoDistance", "expected first param to be list")
		}
		if vals[1].Type() != types.ListType {
			return errFunction("geoDistance", "expected second param to be list")
		}
		var from = orb.Point{}
		var to = orb.Point{}

		for i, val := range vals[0].Value().([]float64) {
			from[i] = val
		}
		for i, val := range vals[1].Value().([]float64) {
			to[i] = val
		}
		return types.Double(geo.DistanceHaversine(from, to))
	},
	"render": func(vals ...ref.Val) ref.Val {
		if len(vals) != 2 {
			return errFunction("render", "expected two params")
		}
		if vals[0].Type() != types.StringType {
			return errFunction("render", "expected first param to be string")
		}
		if vals[1].Type() != types.MapType {
			return errFunction("render", "expected second param to be map")
		}

		buf := bytes.NewBuffer(nil)
		data := cast.ToStringMap(vals[1].Value())
		if err := template.Must(template.New("").Parse(cast.ToString(vals[0].Value()))).Execute(buf, data); err != nil {
			return errFunction("render", err.Error())
		}
		return types.String(buf.String())
	},
	"parseClaims": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("parseClaims", "expected one params")
		}
		claims, err := parseClaims(cast.ToString(vals[0].Value()))
		if err != nil {
			return errFunction("parseClaims", err.Error())
		}
		return types.NewStringInterfaceMap(types.DefaultTypeAdapter, claims)
	},
	"typeOf": func(vals ...ref.Val) ref.Val {
		return types.String(reflect.TypeOf(vals[0].Value()).String())
	},
}

func parseClaims(token string) (map[string]interface{}, error) {
	token = strings.ReplaceAll(token, "Bearer ", "")
	token = strings.ReplaceAll(token, "bearer ", "")
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return nil, errors.New("expected 3 jwt segments")
	}
	payload := []byte(split[1])
	bits, err := base64.RawStdEncoding.DecodeString(string(payload))
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal(bits, &data); err != nil {
		return nil, err
	}
	return data, nil
}
