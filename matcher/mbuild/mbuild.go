// Package mbuild implements functions to build Matcher instances from external sources.
package mbuild

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher"
)

const (
	_mAll                     = "all"
	_mAny                     = "any"
	_mAnything                = "anything"
	_mBoth                    = "both"
	_mContains                = "contains"
	_mEach                    = "each"
	_mEither                  = "either"
	_mEmpty                   = "empty"
	_mEqualTo                 = "equal"
	_mEqualToAlias            = "eq"
	_mEqualToIgnoreCase       = "equalignorecase"
	_mEqualToIgnoreCaseAlias  = "eqi"
	_mEqualJSON               = "equaljson"
	_mEqualJSONAlias          = "eqj"
	_mEqualToStrict           = "equalstrict"
	_mEqualToStrictAlias      = "eqs"
	_mFalsy                   = "falsy"
	_mGreater                 = "greater"
	_mGreaterAlias            = "gt"
	_mGreaterThanOrEqual      = "greatereq"
	_mGreaterThanOrEqualAlias = "gte"
	_mHasKey                  = "haskey"
	_mHasPrefix               = "hasprefix"
	_mHasSuffix               = "hassuffix"
	_mIsIn                    = "isin"
	_mItem                    = "item"
	_mItemsMatch              = "itemsmatch"
	_mJSONPath                = "jsonpath"
	_mField                   = "field"
	_mLength                  = "length"
	_mLen                     = "len"
	_mLowerCase               = "lowercase"
	_mLessThan                = "less"
	_mLessThanAlias           = "lt"
	_mLessThanOrEqual         = "lesseq"
	_mLessThanOrEqualAlias    = "lte"
	_mNil                     = "nil"
	_mRegex                   = "regex"
	_mSome                    = "some"
	_mNot                     = "not"
	_mPresent                 = "present"
	_mSplit                   = "split"
	_mTrim                    = "trim"
	_mTruthy                  = "truthy"
	_mUpperCase               = "uppercase"
	_mURLPath                 = "urlpath"
	_mXOR                     = "xor"
)

// TryBuildMatcher builds a matcher.Matcher if the given parameter is a slice following the matchers convention: [\"<MATCHER_NAME>\", ARG_1, ARG_2...]
// If it is not slice, it will return an equal matcher by default.
func TryBuildMatcher(possibleMatcher any) (m matcher.Matcher, err error) {
	val := reflect.ValueOf(possibleMatcher)
	if possibleMatcher == nil || !val.IsValid() {
		return nil, fmt.Errorf("matcher: definition must be a string or an array in the format: [\"<MATCHER_NAME>\", ARG_1, ARG_2...]")
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		return buildMatcherFromArray(possibleMatcher)
	default:
		return matcher.Equal(possibleMatcher), nil
	}
}

// BuildMatcher always builds a matcher.Matcher from a text or a slice matcher specification.
func BuildMatcher(possibleMatcher any) (m matcher.Matcher, err error) {
	val := reflect.ValueOf(possibleMatcher)
	if possibleMatcher == nil || !val.IsValid() {
		return nil, fmt.Errorf("matcher: definition must be a string or an array in the format: [\"<MATCHER_NAME>\", ARG_1, ARG_2...]")
	}

	switch val.Kind() {
	case reflect.String:
		return discoverAndBuild(val.String(), nil)
	default:
		return buildMatcherFromArray(possibleMatcher)
	}
}

func buildMatcherFromArray(possibleMatcher any) (matcher.Matcher, error) {
	val := reflect.ValueOf(possibleMatcher)
	if val.Len() == 0 {
		return nil, fmt.Errorf("matcher: definition must be a string or an array in the format: [\"<MATCHER_NAME>\", ARG_1, ARG_2...]")
	}

	mk, ok := val.Index(0).Interface().(string)
	if !ok {
		return nil, fmt.Errorf(
			"matcher: first index of a matcher definition must be the matcher name. eg.: [\"<MATCHER_NAME>\", ARGUMENTS...]. got: %v",
			val.Index(0).Interface())
	}

	if val.Len() == 1 {
		return discoverAndBuild(mk, nil)
	} else if val.Len() == 2 {
		return discoverAndBuild(mk, val.Index(1).Interface())
	}

	return discoverAndBuild(mk, val.Slice(1, val.Len()).Interface())
}

func extractMultipleMatchers(v any) ([]matcher.Matcher, error) {
	a, ok := v.([]any)
	if !ok {
		return nil,
			fmt.Errorf("matcher: attempt to build multiple matchers using non-array type. got=%v", reflect.TypeOf(v))
	}

	matchers := make([]matcher.Matcher, len(a))

	for i, entry := range a {
		var mat matcher.Matcher
		var err error

		eType := reflect.TypeOf(entry)
		switch eType.Kind() {
		case reflect.Slice, reflect.Array:
			mat, err = buildMatcherFromArray(entry)
		case reflect.String:
			mat, err = discoverAndBuild(entry.(string), nil)
		}

		if err != nil {
			return nil,
				fmt.Errorf("matcher: error building multiple matchers at index [%d].\n%w", i, err)
		}

		matchers[i] = mat
	}

	return matchers, nil
}

func discoverAndBuild(key string, args any) (ma matcher.Matcher, err error) {
	defer func() {
		if recovery := recover(); recovery != nil {
			err = fmt.Errorf(
				"panic: parsing matcher=%s with args=%v. reason=%v",
				key,
				args,
				recovery,
			)

			return
		}
	}()

	switch strings.ToLower(key) {

	case _mAll:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building matcher list.\n%w", _mAll, err)
		}

		return matcher.All(matchers...), nil

	case _mAny:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building matcher list.\n%w", _mAny, err)
		}

		return matcher.Any(matchers...), nil

	case _mAnything:
		return matcher.Anything(), nil

	case _mContains:
		return matcher.Contain(args), nil

	case _mBoth:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] building error.\n%w", _mBoth, err)
		}

		if len(matchers) != 2 {
			return nil,
				fmt.Errorf("[%s] expects 2 arguments. got=%d", _mBoth, len(matchers))
		}

		return matcher.Both(matchers[0], matchers[1]), nil

	case _mEach:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, fmt.Errorf("[%s] building error.\n%w", _mEach, err)
		}

		return matcher.Each(m), nil

	case _mEither:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil, fmt.Errorf("[%s] error building parameters.\n%w", _mEither, err)
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("[%s] expects 2 arguments. got=%d", _mEither, len(matchers))
		}

		return matcher.Either(matchers[0], matchers[1]), nil

	case _mEmpty:
		return matcher.Empty(), nil

	case _mEqualTo, _mEqualToAlias:
		return matcher.Equal(args), nil

	case _mEqualToIgnoreCase, _mEqualToIgnoreCaseAlias:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s, %s] expects a string argument. got=%v", _mEqualToIgnoreCase, _mEqualToIgnoreCaseAlias, args)
		}

		return matcher.EqualIgnoreCase(str), nil

	case _mEqualJSON, _mEqualJSONAlias:
		return matcher.EqualJSON(args), nil

	case _mEqualToStrict, _mEqualToStrictAlias:
		return matcher.StrictEqual(args), nil

	case _mFalsy:
		return matcher.Falsy(), nil

	case _mGreater, _mGreaterAlias:
		num, err := getFloat64(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s, %s] expects an numeric argument. got=%d", _mGreater, _mGreaterAlias, args)
		}

		return matcher.GreaterThan(num), nil

	case _mGreaterThanOrEqual, _mGreaterThanOrEqualAlias:
		num, err := getFloat64(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s, %s] expects an numeric argument. got=%d", _mGreaterThanOrEqual, _mGreaterThanOrEqualAlias, args)
		}

		return matcher.GreaterThanOrEqual(num), nil

	case _mHasKey:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf(
					"[%s] expects a string argument describing the field path. got=%v",
					_mHasKey,
					args,
				)
		}

		return matcher.HasKey(str), nil

	case _mHasPrefix:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s] expects a string argument. got=%v", _mHasPrefix, args)
		}

		return matcher.HasPrefix(str), nil

	case _mHasSuffix:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s] expects a string argument. got=%v", _mHasSuffix, args)
		}

		return matcher.HasSuffix(str), nil

	case _mIsIn:
		a, ok := args.([]any)
		if !ok {
			return nil,
				fmt.Errorf("[%s] expects an array argument. got=%v", _mIsIn, args)
		}

		return matcher.IsIn(a), nil

	case _mItem:
		a, ok := args.([]any)
		if !ok {
			return nil,
				fmt.Errorf("[%s] expects an array argument. got=%v", _mItem, args)
		}

		if len(a) != 2 {
			return nil,
				fmt.Errorf(
					"[%s] expects at least 2 arguments, 1: item index, 2: Matcher to be applied on the array item. got=%v",
					_mItem,
					args,
				)
		}

		idx, err := getInt(a[0])
		if err != nil {
			return nil,
				fmt.Errorf("[%s] field path must be a number. got=%v", _mItem, a[0])
		}

		m, err := BuildMatcher(a[1])
		if err != nil {
			return nil, fmt.Errorf("[%s] building error.\n%w", _mItem, err)
		}

		return matcher.Item(int(idx), m), nil

	case _mItemsMatch:
		a, ok := args.([]any)
		if !ok {
			return nil,
				fmt.Errorf("[%s] expects an array argument. got=%v", _mItemsMatch, args)
		}

		return matcher.ItemsMatch(a), nil

	case _mJSONPath, _mField:
		a, ok := args.([]any)
		if !ok {
			return nil,
				fmt.Errorf("[%s, %s] expects an array argument. got=%v", _mJSONPath, _mField, args)
		}

		if len(a) != 2 {
			return nil,
				fmt.Errorf(
					"[%s, %s] expects at least 2 arguments, 1: JSON field path, 2: Matcher to be applied on JSON field. got=%v",
					_mJSONPath,
					_mField,
					args,
				)
		}

		chain, ok := a[0].(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s, %s] field path must be a string. got=%v", _mJSONPath, _mField, a[0])
		}

		m, err := BuildMatcher(a[1])
		if err != nil {
			return nil, fmt.Errorf("[%s, %s] building error.\n%w", _mJSONPath, _mField, err)
		}

		return matcher.Field(chain, m), nil

	case _mLength, _mLen:
		num, err := getInt(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s, %s] expects an integer argument. got=%d", _mLen, _mLength, args)
		}

		return matcher.Len(int(num)), nil

	case _mLowerCase:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building.\n%w", _mLowerCase, err)
		}

		return matcher.ToLower(m), nil

	case _mLessThan, _mLessThanAlias:
		num, err := getFloat64(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s, %s] expects an numeric argument. got=%d", _mGreater, _mLessThanAlias, args)
		}

		return matcher.LessThan(num), nil

	case _mLessThanOrEqual, _mLessThanOrEqualAlias:
		num, err := getFloat64(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s, %s] expects an numeric argument. got=%d", _mGreaterThanOrEqual, _mLessThanOrEqualAlias, args)
		}

		return matcher.LessThanOrEqual(num), nil

	case _mNil:
		return matcher.IsNil(), nil

	case _mNot:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building.\n%w", _mNot, err)
		}

		return matcher.Not(m), nil

	case _mPresent:
		return matcher.Present(), nil

	case _mRegex:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s] expects a string argument. got=%v", _mRegex, args)
		}

		return matcher.Matches(str), nil

	case _mSome:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, fmt.Errorf("[%s] error building.\n%w", _mSome, err)
		}

		return matcher.Some(m), nil

	case _mSplit:
		a, ok := args.([]any)
		if !ok {
			return nil, fmt.Errorf("[%s] expects an argument of type array. got=%v", _mSplit, args)
		}

		if len(a) != 2 {
			return nil,
				fmt.Errorf("[%s] expects two arguments. 1: Matcher, 2: Separator. got=%d", _mSplit, len(a))
		}

		separator, ok := a[0].(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s] second parameter must be a string. got=%v", _mSplit, a[1])
		}

		m, err := BuildMatcher(a[1])
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building.\n%w", _mSplit, err)
		}

		return matcher.Split(separator, m), nil

	case _mTrim:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building.\n%w", _mTrim, err)
		}

		return matcher.Trim(m), nil

	case _mTruthy:
		return matcher.Truthy(), nil

	case _mUpperCase:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, fmt.Errorf("[%s] building error.\n%w", _mUpperCase, err)
		}

		return matcher.ToUpper(m), nil

	case _mURLPath:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s] matcher expects a string argument. got=%v", _mURLPath, args)
		}

		return matcher.URLPath(str), nil

	case _mXOR:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil, fmt.Errorf("[%s] building error.\n%w", _mXOR, err)
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("[%s] expects two parameters. got=%d", _mXOR, len(matchers))
		}

		return matcher.XOR(matchers[0], matchers[1]), nil

	default:
		return nil, fmt.Errorf("[%s] unknown matcher key=%s", key, key)
	}
}

func getInt(v any) (int64, error) {
	vv := reflect.ValueOf(v)
	k := vv.Kind()

	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return vv.Int(), nil
	case reflect.Float64, reflect.Float32:
		return int64(vv.Float()), nil
	}

	return 0, fmt.Errorf("invalid integer value %v", vv.Interface())
}

func getFloat64(v any) (float64, error) {
	vv := reflect.ValueOf(v)
	k := vv.Kind()

	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(vv.Int()), nil
	case reflect.Float64, reflect.Float32:
		return vv.Float(), nil
	}

	return 0, fmt.Errorf("invalid float64 value %v", vv.Interface())
}
