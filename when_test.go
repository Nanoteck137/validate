// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package validate

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func abcValidation(val string) bool {
	return val == "abc"
}

func TestWhen(t *testing.T) {
	abcRule := NewStringRule(abcValidation, "wrong_abc")
	validateMeRule := NewStringRule(validateMe, "wrong_me")

	tests := []struct {
		tag       string
		condition bool
		value     interface{}
		rules     []Rule
		elseRules []Rule
		err       string
	}{
		// True condition
		{"t1.1", true, nil, []Rule{}, []Rule{}, ""},
		{"t1.2", true, "", []Rule{}, []Rule{}, ""},
		{"t1.3", true, "", []Rule{abcRule}, []Rule{}, ""},
		{"t1.4", true, 12, []Rule{Required}, []Rule{}, ""},
		{"t1.5", true, nil, []Rule{Required}, []Rule{}, "cannot be blank"},
		{"t1.6", true, "123", []Rule{abcRule}, []Rule{}, "wrong_abc"},
		{"t1.7", true, "abc", []Rule{abcRule}, []Rule{}, ""},
		{"t1.8", true, "abc", []Rule{abcRule, abcRule}, []Rule{}, ""},
		{"t1.9", true, "abc", []Rule{abcRule, validateMeRule}, []Rule{}, "wrong_me"},
		{"t1.10", true, "me", []Rule{abcRule, validateMeRule}, []Rule{}, "wrong_abc"},
		{"t1.11", true, "me", []Rule{}, []Rule{abcRule}, ""},

		// False condition
		{"t2.1", false, "", []Rule{}, []Rule{}, ""},
		{"t2.2", false, "", []Rule{abcRule}, []Rule{}, ""},
		{"t2.3", false, "abc", []Rule{abcRule}, []Rule{}, ""},
		{"t2.4", false, "abc", []Rule{abcRule, abcRule}, []Rule{}, ""},
		{"t2.5", false, "abc", []Rule{abcRule, validateMeRule}, []Rule{}, ""},
		{"t2.6", false, "me", []Rule{abcRule, validateMeRule}, []Rule{}, ""},
		{"t2.7", false, "", []Rule{abcRule, validateMeRule}, []Rule{}, ""},
		{"t2.8", false, "me", []Rule{}, []Rule{abcRule, validateMeRule}, "wrong_abc"},
	}

	for _, test := range tests {
		err := Validate(test.value, When(test.condition, test.rules...).Else(test.elseRules...))
		assertError(t, test.err, err, test.tag)
	}
}

type ctxKey int

const (
	contains ctxKey = iota
)

func TestWhenWithContext(t *testing.T) {
	rule := WithContext(func(ctx context.Context, value interface{}) error {
		if !strings.Contains(value.(string), ctx.Value(contains).(string)) {
			return errors.New("unexpected value")
		}
		return nil
	})
	ctx1 := context.WithValue(context.Background(), contains, "abc")
	ctx2 := context.WithValue(context.Background(), contains, "xyz")

	tests := []struct {
		tag       string
		condition bool
		value     interface{}
		ctx       context.Context
		err       string
	}{
		// True condition
		{"t1.1", true, "abc", ctx1, ""},
		{"t1.2", true, "abc", ctx2, "unexpected value"},
		{"t1.3", true, "xyz", ctx1, "unexpected value"},
		{"t1.4", true, "xyz", ctx2, ""},

		// False condition
		{"t2.1", false, "abc", ctx1, ""},
		{"t2.2", false, "abc", ctx2, "unexpected value"},
		{"t2.3", false, "xyz", ctx1, "unexpected value"},
		{"t2.4", false, "xyz", ctx2, ""},
	}

	for _, test := range tests {
		err := ValidateWithContext(test.ctx, test.value, When(test.condition, rule).Else(rule))
		assertError(t, test.err, err, test.tag)
	}
}
