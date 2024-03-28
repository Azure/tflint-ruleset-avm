package interfaces_test

import (
	"fmt"
	"github.com/Azure/tflint-ruleset-avm/interfaces"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	called := false
	cases := []struct {
		c                            func() interfaces.Checker
		secondCallbackShouldBeCalled bool
		description                  string
	}{
		{
			c: func() interfaces.Checker {
				return interfaces.NewChecker().Check(func() (bool, error) {
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
			c: func() interfaces.Checker {
				return interfaces.NewChecker().Check(func() (bool, error) {
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
			c: func() interfaces.Checker {
				return interfaces.NewChecker().Check(func() (bool, error) {
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
		c                      func() interfaces.Checker
		callbackShouldBeCalled bool
		description            string
	}{
		{
			c: func() interfaces.Checker {
				checker := interfaces.NewChecker()
				_, checker = interfaces.CheckWithReturnValue(checker, func() (int, bool, error) {
					return 0, true, nil
				})
				interfaces.CheckWithReturnValue(checker, func() (string, bool, error) {
					called = true
					return "", true, nil
				})
				return checker
			},
			callbackShouldBeCalled: true,
			description:            "happy_path_second_callback_should_be_called",
		},
		{
			c: func() interfaces.Checker {
				checker := interfaces.NewChecker()
				_, checker = interfaces.CheckWithReturnValue(checker, func() (int, bool, error) {
					return 0, false, nil
				})
				interfaces.CheckWithReturnValue(checker, func() (string, bool, error) {
					called = true
					return "", true, nil
				})
				return checker
			},
			callbackShouldBeCalled: false,
			description:            "false_path_second_callback_should_not_be_called",
		},
		{
			c: func() interfaces.Checker {
				checker := interfaces.NewChecker()
				_, checker = interfaces.CheckWithReturnValue(checker, func() (int, bool, error) {
					return 0, true, fmt.Errorf("error")
				})
				interfaces.CheckWithReturnValue(checker, func() (string, bool, error) {
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
