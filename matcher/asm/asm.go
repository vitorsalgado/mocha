package asm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher"
)

// matcher.Matcher keys
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
	_mHaveKey          = "havekey"
	_mHasPrefix        = "hasprefix"
	_mStartsWith       = "startswith"
	_mHasSuffix        = "hassuffix"
	_mEqualsIgnoreCase = "equalsignorecase"
	_mEqualFold        = "equalfold"
	_mJSONPath         = "jsonpath"
	_mLen              = "len"
	_mLowerCase        = "lowercase"
	_mRegex            = "regex"
	_mNot              = "not"
	_mPresent          = "present"
	_mSplit            = "split"
	_mTrim             = "trim"
	_mUpperCase        = "uppercase"
	_mURLPath          = "urlpath"
	_mXOR              = "xor"
)

func BuildMatcher(v any) (m matcher.Matcher, err error) {
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
		return matcher.EqualIgnoreCase(v.(string)), nil
	case reflect.Slice:
		val := reflect.ValueOf(v)
		if val.Len() == 0 {
			return nil, fmt.Errorf("array must equal or greather than 1")
		}

		mk, ok := val.Index(0).Interface().(string)
		if !ok {
			return nil, fmt.Errorf("first index must be the matcher name")
		}

		return discoverAndBuild(strings.ToLower(mk), val.Slice(1, val.Len()).Interface())
	default:
		return matcher.Equal(v), nil
	}
}

func extractMultipleMatchers(v any) ([]matcher.Matcher, error) {
	a, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("only arrays")
	}

	matchers := make([]matcher.Matcher, len(a))

	for _, entry := range a {
		mat, err := BuildMatcher(entry)
		if err != nil {
			return nil, err
		}

		matchers = append(matchers, mat)
	}

	return matchers, nil
}

func discoverAndBuild(key string, args any) (m matcher.Matcher, err error) {
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

	case _mAll, _mAllOf:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil, err
		}

		return matcher.AllOf(matchers...), nil

	case _mAny, _mAnyOf:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil, err
		}

		return matcher.AnyOf(matchers...), nil

	case _mContain, _mContains:
		return matcher.Contain(args), nil

	case _mBoth:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil, err
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("required 2")
		}

		return matcher.Both(matchers[0], matchers[1]), nil

	case _mEach:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.Each(m), nil

	case _mEither:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil, err
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("required 2")
		}

		return matcher.Either(matchers[0], matchers[1]), nil

	case _mEmpty:
		return matcher.Empty(), nil

	case _mEqual, _mEqualTo:
		return matcher.Equal(args), nil

	case _mEqualIgnoreCase, _mEqualsIgnoreCase, _mEqualFold:
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string required")
		}

		return matcher.EqualIgnoreCase(str), nil

	case _mEqualJSON:
		return matcher.EqualJSON(args), nil

	case _mHasKey, _mHaveKey:
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string required")
		}

		return matcher.HaveKey(str), nil

	case _mHasPrefix, _mStartsWith:
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string required")
		}

		return matcher.HasPrefix(str), nil

	case _mHasSuffix, _mEndsWith:
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string required")
		}

		return matcher.HasSuffix(str), nil

	case _mJSONPath:
		a, ok := args.([]any)
		if !ok {
			return nil, fmt.Errorf("array")
		}

		if len(a) != 2 {
			return nil, fmt.Errorf("")
		}

		chain, ok := a[0].(string)
		if !ok {
			return nil, fmt.Errorf("path string")
		}

		m, err := BuildMatcher(a[1])
		if err != nil {
			return nil, err
		}

		return matcher.JSONPath(chain, m), nil

	case _mLen:
		num, ok := args.(float64)
		if !ok {
			return nil, fmt.Errorf("number required")
		}

		return matcher.HaveLen(int(num)), nil

	case _mLowerCase:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.ToLower(m), nil

	case _mRegex:
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("string")
		}

		return matcher.Matches(str), nil

	case _mNot:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.Not(m), nil

	case _mPresent:
		return matcher.Present(), nil

	case _mSplit:
		a, ok := args.([]any)
		if !ok {
			return nil, fmt.Errorf("array")
		}

		if len(a) != 2 {
			return nil, fmt.Errorf("")
		}

		separator, ok := a[0].(string)
		if !ok {
			return nil, fmt.Errorf("separator string")
		}

		m, err := BuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.Split(separator, m), nil

	case _mTrim:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.Trim(m), nil

	case _mUpperCase:
		m, err := BuildMatcher(args)
		if err != nil {
			return nil, err
		}

		return matcher.ToUpper(m), nil

	case _mURLPath:
		str, ok := args.(string)
		if !ok {
			return nil, fmt.Errorf("")
		}

		return matcher.URLPath(str), nil

	case _mXOR:
		matchers, err := extractMultipleMatchers(args)
		if err != nil {
			return nil, err
		}

		if len(matchers) != 2 {
			return nil, fmt.Errorf("required 2")
		}

		return matcher.XOR(matchers[0], matchers[1]), nil

	default:
		return nil, fmt.Errorf("unknown matcher key=%s", key)
	}
}
