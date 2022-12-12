package matcher

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	_mAll              = "all"
	_mAllOf            = "allof"
	_mAny              = "any"
	_mAnyOf            = "anyof"
	_mContain          = "contain"
	_mContains         = "contains"
	_mBoth             = "both"
	_mEach             = "each"
	_mEither           = "either"
	_mEmpty            = "empty"
	_mEndsWith         = "endswith"
	_mEqual            = "equal"
	_mEqualTo          = "equalTo"
	_mEqualIgnoreCase  = "equalignorecase"
	_mEqualJSON        = "equaljson"
	_mHasKey           = "haskey"
	_mHasPrefix        = "hasprefix"
	_mStartsWith       = "startswith"
	_mHasSuffix        = "hassuffix"
	_mEqualsIgnoreCase = "equalsignorecase"
	_mEqualFold        = "equalfold"
	_mJSONPath         = "jsonpath"
	_mField            = "field"
	_mLen              = "len"
	_mLowerCase        = "lowercase"
	_mRegex            = "regex"
	_mSome             = "some"
	_mNot              = "not"
	_mPresent          = "present"
	_mSplit            = "split"
	_mTrim             = "trim"
	_mUpperCase        = "uppercase"
	_mURLPath          = "urlpath"
	_mXOR              = "xor"
)

func BuildMatcher(v any) (m Matcher, err error) {
	defer func() {
		if r := recover(); r != nil {
			m = nil
			err = fmt.Errorf("panic=%v", r)
			return
		}
	}()

	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.String:
		return EqualIgnoreCase(v.(string)), nil
	case reflect.Slice, reflect.Array:
		val := reflect.ValueOf(v)
		if val.Len() == 0 {
			return nil,
				fmt.Errorf("matcher definition must be a string or an array in the format: [\"equal\", \"test\"]")
		}

		mk, ok := val.Index(0).Interface().(string)
		if !ok {
			return nil,
				fmt.Errorf(
					"first index of a matcher definition must be the matcher name. eg.: [\"equal\", \"test\"]. got: %v",
					val.Index(0).Interface(),
				)
		}

		return discoverAndBuild(strings.ToLower(mk), val.Slice(1, val.Len()).Interface())
	default:
		return Equal(v), nil
	}
}

func extractMultipleMatchers(v any) ([]Matcher, error) {
	a, ok := v.([]any)
	if !ok {
		return nil,
			fmt.Errorf("attempt to build multiple matchers using non-array type. got=%v", v)
	}

	matchers := make([]Matcher, len(a))

	for i, entry := range a {
		mat, err := BuildMatcher(entry)
		if err != nil {
			return nil,
				fmt.Errorf("error building multiple matchers at index [%d]. %w", i, err)
		}

		matchers = append(matchers, mat)
	}

	return matchers, nil
}

func discoverAndBuild(key string, args any) (m Matcher, err error) {
	defer func() {
		if recovery := recover(); recovery != nil {
			m = nil
			err = fmt.Errorf(
				"panic parsing matcher=%s with args=%v. reason=%v",
				key,
				args,
				recovery,
			)

			return
		}
	}()

	switch strings.ToLower(key) {

	// TODO: use one
	case _mAll, _mAllOf:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s, %s] error building matcher list. %w", _mAll, _mAllOf, err)
		}

		return AllOf(matchers...), nil

		// TODO: use one
	case _mAny, _mAnyOf:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s, %s] error building matcher list. %w", _mAny, _mAnyOf, err)
		}

		return AnyOf(matchers...), nil

		// TODO: use one
	case _mContain, _mContains:
		return Contain(args), nil

	case _mBoth:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] building error. %w", _mBoth, err)
		}

		if len(matchers) != 2 {
			return nil,
				fmt.Errorf("[%s] expects 2 arguments. got=%d", _mBoth, len(matchers))
		}

		m1, err := BuildMatcher(matchers[0])
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building first matcher. %w", _mBoth, err)
		}

		m2, err := BuildMatcher(matchers[1])
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building second matcher. %w", _mBoth, err)
		}

		return Both(m1, m2), nil

	case _mEach:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, fmt.Errorf("[%s] building error. %w", _mEach, err)
		}

		return Each(m), nil

	case _mEither:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil, fmt.Errorf("[%s] error building parameters. %w", _mEither, err)
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("[%s] expects 2 arguments. got=%d", _mEither, len(matchers))
		}

		m1, err := BuildMatcher(matchers[0])
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building first matcher. %w", _mEither, err)
		}

		m2, err := BuildMatcher(matchers[1])
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building second matcher. %w", _mEither, err)
		}

		return Either(m1, m2), nil

	case _mEmpty:
		return Empty(), nil

		// TODO: use one
	case _mEqual, _mEqualTo:
		return Equal(args), nil

		// TODO: use one
	case _mEqualIgnoreCase, _mEqualsIgnoreCase, _mEqualFold:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf(
					"[%s, %s, %s] expects a string argument. got=%v",
					_mEqualIgnoreCase, _mEqualsIgnoreCase, _mEqualFold,
					args,
				)
		}

		return EqualIgnoreCase(str), nil

	case _mEqualJSON:
		return EqualJSON(args), nil

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

		return HaveKey(str), nil

		// TODO: use one
	case _mHasPrefix, _mStartsWith:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf(
					"[%s, %s] expects a string argument. got=%v",
					_mHasPrefix,
					_mStartsWith,
					args,
				)
		}

		return HasPrefix(str), nil

		// TODO: use one
	case _mHasSuffix, _mEndsWith:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf(
					"[%s, %s] expects a string argument. got=%v",
					_mHasSuffix,
					_mEndsWith,
					args,
				)
		}

		return HasSuffix(str), nil

		// TODO: use one
	case _mJSONPath, _mField:
		a, ok := args.([]any)
		if !ok {
			return nil,
				fmt.Errorf("[%s] expects an array argument. got=%v", _mJSONPath, args)
		}

		if len(a) != 2 {
			return nil,
				fmt.Errorf(
					"[%s] expects at least 2 arguments, 1: JSON field path, 2: Matcher to be applied on JSON field. got=%v",
					_mJSONPath,
					args,
				)
		}

		chain, ok := a[0].(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s] field path must be a string. got=%v", _mJSONPath, a[0])
		}

		m, err := BuildMatcher(a[1])
		if err != nil {
			return nil, fmt.Errorf("[%s] building error. %w", _mJSONPath, err)
		}

		return JSONPath(chain, m), nil

	case _mLen:
		num, ok := args.(float64)
		if !ok {
			return nil,
				fmt.Errorf("[%s] expects an integer argument. got=%d", _mLen, args)
		}

		return HaveLen(int(num)), nil

	case _mLowerCase:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building. %w", _mLowerCase, err)
		}

		return ToLower(m), nil

	case _mNot:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building. %w", _mNot, err)
		}

		return Not(m), nil

	case _mPresent:
		return Present(), nil

	case _mRegex:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s] expects a string argument. got=%v", _mRegex, args)
		}

		return Matches(str), nil

	case _mSome:
		return Some(args.([]any)), nil

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
				fmt.Errorf("[%s] error building. %w", _mSplit, err)
		}

		return Split(separator, m), nil

	case _mTrim:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building. %w", _mTrim, err)
		}

		return Trim(m), nil

	case _mUpperCase:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, fmt.Errorf("[%s] building error. %w", _mUpperCase, err)
		}

		return ToUpper(m), nil

	case _mURLPath:
		str, ok := args.(string)
		if !ok {
			return nil,
				fmt.Errorf("[%s] matcher expects a string argument. got=%v", _mURLPath, args)
		}

		return URLPath(str), nil

	case _mXOR:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil, fmt.Errorf("[%s] building error. %w", _mXOR, err)
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("[%s] expects two parameters. got=%d", _mXOR, len(matchers))
		}

		m1, err := BuildMatcher(matchers[0])
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building first conditon. %w", _mXOR, err)
		}

		m2, err := BuildMatcher(matchers[1])
		if err != nil {
			return nil,
				fmt.Errorf("[%s] error building second conditon. %w", _mXOR, err)
		}

		return XOR(m1, m2), nil

	default:
		return nil, fmt.Errorf("unknown matcher key=%s", key)
	}
}
