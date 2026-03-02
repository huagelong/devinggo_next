// Package system
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package system

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// TestSplitInfoParams 测试 info 参数分割功能
func TestSplitInfoParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "单个参数",
			input:    "user_count",
			expected: []string{"user_count"},
		},
		{
			name:     "多个参数",
			input:    "user_count,subscription_count",
			expected: []string{"user_count", "subscription_count"},
		},
		{
			name:     "带空格的参数",
			input:    " user_count , subscription_count ",
			expected: []string{"user_count", "subscription_count"},
		},
		{
			name:     "空字符串",
			input:    "",
			expected: []string{},
		},
		{
			name:     "只有逗号",
			input:    ",,,",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitInfoParams(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("splitInfoParams(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("splitInfoParams(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestSplitInfo 测试频道查询的 info 参数分割
func TestSplitInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "单个参数",
			input:    "user_count",
			expected: []string{"user_count"},
		},
		{
			name:     "多个参数",
			input:    "user_count,subscription_count",
			expected: []string{"user_count", "subscription_count"},
		},
		{
			name:     "带空格",
			input:    "user_count, subscription_count",
			expected: []string{"user_count", "subscription_count"},
		},
		{
			name:     "空字符串",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitInfo(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("splitInfo(%q) length = %d, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("splitInfo(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestAbs 测试绝对值函数
func TestAbs(t *testing.T) {
	tests := []struct {
		input    int64
		expected int64
	}{
		{input: 10, expected: 10},
		{input: -10, expected: 10},
		{input: 0, expected: 0},
		{input: 100, expected: 100},
		{input: -100, expected: 100},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d, want %d", tt.input, result, tt.expected)
		}
	}
}

// TestVerifyWebhookSignature 测试 Webhook 签名验证
func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test-secret"
	body := []byte(`{"time_ms":123456,"events":[]}`)

	// 生成正确的签名
	signature := generateTestSignature(body, secret)

	// 测试正确的签名
	if !verifyWebhookSignature(body, signature, secret) {
		t.Error("verifyWebhookSignature should return true for valid signature")
	}

	// 测试错误的签名
	if verifyWebhookSignature(body, "invalid-signature", secret) {
		t.Error("verifyWebhookSignature should return false for invalid signature")
	}

	// 测试空签名
	if verifyWebhookSignature(body, "", secret) {
		t.Error("verifyWebhookSignature should return false for empty signature")
	}
}

// generateTestSignature 生成测试签名
func generateTestSignature(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
