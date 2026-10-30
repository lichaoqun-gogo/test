package openapi

import "testing"

// TestSecretSign 验证 auth-secret 算法：hex(md5(secret + timestamp))，32 位小写
func TestSecretSign(t *testing.T) {
	// md5("ab"+"c") = md5("abc") = 900150983cd24fb0d6963f7d28e17f72（公开已知值）
	if got := SecretSign("ab", "c"); got != "900150983cd24fb0d6963f7d28e17f72" {
		t.Errorf("SecretSign(ab, c) = %q, want 900150983cd24fb0d6963f7d28e17f72", got)
	}
	if got := SecretSign("08iUN@diT0PVidgR4c%Y00eTjjX0bcc", "1697000000000"); len(got) != 32 {
		t.Errorf("SecretSign 输出应为 32 位 hex，实际 %d 位: %s", len(got), got)
	}
}

// TestBuildSignString 验证 auth-sign 待签字符串拼接规则：
// 参数按 key 字典序，数组内部排序后逐元素 key=value&，
// 末尾追加 appid/timestamp 与 auth-secret={签名值}（不拼明文密钥）
func TestBuildSignString(t *testing.T) {
	params := map[string]any{
		"b": "2",
		"a": "1",
	}
	// authSecret 参数即 SecretSign(secret, timestamp) 的签名值
	const authSecret = "signed-secret-value"
	want := "a=1&b=2&auth-appid=osm&auth-timestamp=1697000000000&auth-secret=" + authSecret
	if got := buildSignString("osm", authSecret, "1697000000000", params); got != want {
		t.Errorf("buildSignString = %q, want %q", got, want)
	}
}

// TestBuildSignStringArray 验证数组参数：内部排序后逐元素 key=value&
func TestBuildSignStringArray(t *testing.T) {
	params := map[string]any{
		"ids": []string{"2", "1"},
		"x":   "y",
	}
	const authSecret = "signed-secret-value"
	want := "ids=1&ids=2&x=y&auth-appid=osm&auth-timestamp=1697000000000&auth-secret=" + authSecret
	if got := buildSignString("osm", authSecret, "1697000000000", params); got != want {
		t.Errorf("buildSignString = %q, want %q", got, want)
	}
}

// TestMD5Sign 验证 auth-sign = hex(md5(待签字符串))，且待签字符串中的
// auth-secret 取 SecretSign(secret, timestamp) 的签名值（与请求头 auth-secret 同源，校验方可复算）
func TestMD5Sign(t *testing.T) {
	params := map[string]any{"a": "1"}
	authSecret := SecretSign("SECRET", "1697000000000")
	got := MD5Sign("osm", "SECRET", "1697000000000", params)
	want := md5Hex(buildSignString("osm", authSecret, "1697000000000", params))
	if got != want {
		t.Errorf("MD5Sign 与 buildSignString 不一致: %q vs %q", got, want)
	}
	if len(got) != 32 {
		t.Errorf("MD5Sign 输出应为 32 位 hex，实际 %d 位", len(got))
	}
}
