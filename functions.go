package trigger

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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
	"golang.org/x/crypto/sha3"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"html/template"
	"io"
	"net/url"
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
			decls.NewInstanceOverload(
				"sha1_overload",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "sha1_overload",
			Function: defaultFuncMap["sha1"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["sha1"](value)
			},
		},
	},
	"sha256": {
		decl: decls.NewFunction("sha256",
			decls.NewInstanceOverload(
				"sha256_overload",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "sha256_overload",
			Function: defaultFuncMap["sha256"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["sha256"](value)
			},
		},
	},
	"sha3": {
		decl: decls.NewFunction("sha3",
			decls.NewInstanceOverload(
				"sha3_overload",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "sha3_overload",
			Function: defaultFuncMap["sha3"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["sha3"](value)
			},
		},
	},
	"base64Encode": {
		decl: decls.NewFunction("base64Encode",
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
	"replace": {
		decl: decls.NewFunction("replace",
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
	"geoDistance": {
		decl: decls.NewFunction("geoDistance",
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
			decls.NewInstanceOverload(
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
	"parseHeader": {
		decl: decls.NewFunction("parseHeader",
			decls.NewInstanceOverload(
				"parseHeader",
				[]*expr.Type{decls.String},
				strMap,
			),
		),
		overload: &functions.Overload{
			Operator: "parseHeader",
			Function: defaultFuncMap["parseHeader"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["parseHeader"](value)
			},
		},
	},
	"parseSignature": {
		decl: decls.NewFunction("parseSignature",
			decls.NewInstanceOverload(
				"parseSignature",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "parseSignature",
			Function: defaultFuncMap["parseSignature"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["parseSignature"](value)
			},
		},
	},
	"typeOf": {
		decl: decls.NewFunction("typeOf",
			decls.NewInstanceOverload(
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
	"encrypt": {
		decl: decls.NewFunction("encrypt",
			decls.NewInstanceOverload(
				"encrypt",
				[]*expr.Type{decls.String, decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "encrypt",
			Function: defaultFuncMap["encrypt"],
			Binary: func(value ref.Val, value2 ref.Val) ref.Val {
				return defaultFuncMap["encrypt"](value, value2)
			},
		},
	},
	"decrypt": {
		decl: decls.NewFunction("decrypt",
			decls.NewInstanceOverload(
				"decrypt",
				[]*expr.Type{decls.String, decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "decrypt",
			Function: defaultFuncMap["decrypt"],
			Binary: func(value ref.Val, value2 ref.Val) ref.Val {
				return defaultFuncMap["decrypt"](value, value2)
			},
		},
	},
	"parseHost": {
		decl: decls.NewFunction("parseHost",
			decls.NewInstanceOverload(
				"parseHost",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "parseHost",
			Function: defaultFuncMap["parseHost"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["parseHost"](value)
			},
		},
	},
	"parseQuery": {
		decl: decls.NewFunction("parseQuery",
			decls.NewInstanceOverload(
				"parseQuery",
				[]*expr.Type{decls.String},
				strMap,
			),
		),
		overload: &functions.Overload{
			Operator: "parseQuery",
			Function: defaultFuncMap["parseQuery"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["parseQuery"](value)
			},
		},
	},
	"parsePath": {
		decl: decls.NewFunction("parsePath",
			decls.NewInstanceOverload(
				"parsePath",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "parsePath",
			Function: defaultFuncMap["parsePath"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["parsePath"](value)
			},
		},
	},
	"parseScheme": {
		decl: decls.NewFunction("parseScheme",
			decls.NewInstanceOverload(
				"parseScheme",
				[]*expr.Type{decls.String},
				decls.String,
			),
		),
		overload: &functions.Overload{
			Operator: "parseScheme",
			Function: defaultFuncMap["parseScheme"],
			Unary: func(value ref.Val) ref.Val {
				return defaultFuncMap["parseScheme"](value)
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
	"parseHeader": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("parseHeader", "expected one params")
		}
		header, err := parseHeader(cast.ToString(vals[0].Value()))
		if err != nil {
			return errFunction("parseClaims", err.Error())
		}
		return types.NewStringInterfaceMap(types.DefaultTypeAdapter, header)
	},
	"parseSignature": func(vals ...ref.Val) ref.Val {
		if len(vals) != 1 {
			return errFunction("parseSignature", "expected one params")
		}
		sig, err := parseSignature(cast.ToString(vals[0].Value()))
		if err != nil {
			return errFunction("parseSignature", err.Error())
		}
		return types.String(sig)
	},
	"typeOf": func(vals ...ref.Val) ref.Val {
		return types.String(reflect.TypeOf(vals[0].Value()).String())
	},
	"encrypt": func(vals ...ref.Val) ref.Val {
		encrypted, err := encrypt([]byte(cast.ToString(vals[0].Value())), cast.ToString(vals[1].Value()))
		if err != nil {
			return errFunction("encrypt", err.Error())
		}
		return types.String(encrypted)
	},
	"decrypt": func(vals ...ref.Val) ref.Val {
		decrypted, err := decrypt([]byte(cast.ToString(vals[0].Value())), cast.ToString(vals[1].Value()))
		if err != nil {
			return errFunction("decrypt", err.Error())
		}
		return types.String(decrypted)
	},
	"parseHost": func(vals ...ref.Val) ref.Val {
		u, err := url.Parse(cast.ToString(vals[0].Value()))
		if err != nil {
			return errFunction("parseHost", err.Error())
		}
		return types.String(u.Host)
	},
	"parsePath": func(vals ...ref.Val) ref.Val {
		u, err := url.Parse(cast.ToString(vals[0].Value()))
		if err != nil {
			return errFunction("parsePath", err.Error())
		}
		return types.String(u.Path)
	},
	"parseScheme": func(vals ...ref.Val) ref.Val {
		u, err := url.Parse(cast.ToString(vals[0].Value()))
		if err != nil {
			return errFunction("parseScheme", err.Error())
		}
		return types.String(u.Scheme)
	},
	"parseQuery": func(vals ...ref.Val) ref.Val {
		u, err := url.Parse(cast.ToString(vals[0].Value()))
		if err != nil {
			return errFunction("parseQuery", err.Error())
		}
		data := map[string]interface{}{}
		for k, vals := range u.Query() {
			if len(vals) > 0 {
				data[k] = vals[0]
			}
		}
		return types.NewStringInterfaceMap(types.DefaultTypeAdapter, data)
	},
}

func parseJWT(token string) ([]string, error) {
	token = strings.ReplaceAll(token, "Bearer ", "")
	token = strings.ReplaceAll(token, "bearer ", "")
	split := strings.Split(token, ".")
	if len(split) != 3 {
		return nil, errors.New("expected 3 jwt segments")
	}
	return split, nil
}

func parseHeader(token string) (map[string]interface{}, error) {
	split, err := parseJWT(token)
	if err != nil {
		return nil, err
	}
	payload := []byte(split[0])
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

func parseSignature(token string) (string, error) {
	split, err := parseJWT(token)
	if err != nil {
		return "", err
	}
	return split[2], nil
}

func parseClaims(token string) (map[string]interface{}, error) {
	split, err := parseJWT(token)
	if err != nil {
		return nil, err
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

func encrypt(key []byte, message string) (string, error) {
	plainText := []byte(message)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func decrypt(key []byte, securemess string) (string, error) {
	cipherText, err := base64.URLEncoding.DecodeString(securemess)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipher text length invalid")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}
