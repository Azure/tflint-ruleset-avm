package rules_test

import (
	"fmt"
	"testing"

	"github.com/Azure/tflint-ruleset-avm/rules"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	called := false
	cases := []struct {
		c                            func() rules.Checker
		secondCallbackShouldBeCalled bool
		description                  string
	}{
		{
			c: func() rules.Checker {
				return rules.NewChecker().Check(func() (bool, error) {
					return true, nil
				}).Check(func() (bool, error) {
					called = true
					return true, nil
				})
			},
			secondCallbackShouldBeCalled: true,
			description:                  "happy_path_second_callback_should_be_called",
		},
		{
			c: func() rules.Checker {
				return rules.NewChecker().Check(func() (bool, error) {
					return false, nil
				}).Check(func() (bool, error) {
					called = true
					return true, nil
				})
			},
			secondCallbackShouldBeCalled: false,
			description:                  "false_path_second_callback_should_not_be_called",
		},
		{
			c: func() rules.Checker {
				return rules.NewChecker().Check(func() (bool, error) {
					return true, fmt.Errorf("error")
				}).Check(func() (bool, error) {
					called = true
					return true, nil
				})
			},
			secondCallbackShouldBeCalled: false,
			description:                  "error_path_second_callback_should_be_called",
		},
	}
	for _, cc := range cases {
		tc := cc
		t.Run(tc.description, func(t *testing.T) {
			called = false
			_ = cc.c()
			assert.Equal(t, cc.secondCallbackShouldBeCalled, called)
		})
	}
}

func TestCheckWithReturnValue(t *testing.T) {
	called := false
	cases := []struct {
		c                      func() rules.Checker
		callbackShouldBeCalled bool
		description            string
	}{
		{
			c: func() rules.Checker {
				checker := rules.NewChecker()
				_, checker = rules.CheckWithReturnValue(checker, func() (int, bool, error) {
					return 0, true, nil
				})
				rules.CheckWithReturnValue(checker, func() (string, bool, error) {
					called = true
					return "", true, nil
				})
				return checker
			},
			callbackShouldBeCalled: true,
			description:            "happy_path_second_callback_should_be_called",
		},
		{
			c: func() rules.Checker {
				checker := rules.NewChecker()
				_, checker = rules.CheckWithReturnValue(checker, func() (int, bool, error) {
					return 0, false, nil
				})
				rules.CheckWithReturnValue(checker, func() (string, bool, error) {
					called = true
					return "", true, nil
				})
				return checker
			},
			callbackShouldBeCalled: false,
			description:            "false_path_second_callback_should_not_be_called",
		},
		{
			c: func() rules.Checker {
				checker := rules.NewChecker()
				_, checker = rules.CheckWithReturnValue(checker, func() (int, bool, error) {
					return 0, true, fmt.Errorf("error")
				})
				rules.CheckWithReturnValue(checker, func() (string, bool, error) {
					called = true
					return "", true, nil
				})
				return checker
			},
			callbackShouldBeCalled: false,
			description:            "error_path_second_callback_should_be_called",
		},
	}
	for _, cc := range cases {
		tc := cc
		t.Run(tc.description, func(t *testing.T) {
			called = false
			_ = cc.c()
			assert.Equal(t, cc.callbackShouldBeCalled, called)
		})
	}
}
