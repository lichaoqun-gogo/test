package openapi

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// md5Hex 计算字符串的 MD5 并以 32 位小写 hex 输出
func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// SecretSign 生成请求头 auth-secret：hex(md5(secret + timestamp))
func SecretSign(secret, timestamp string) string {
	return md5Hex(secret + timestamp)
}

// buildSignString 拼接 auth-sign 的待签字符串（业界通用规则，与 IAM 网关约定一致）：
//  1. 所有请求参数按 key 字典序升序排列
//  2. 每个 key-value：数组值内部排序后逐元素拼 key=value&；单值拼 key=value&
//  3. 追加 auth-appid={appid}&
//  4. 追加 auth-timestamp={timestamp}&
//  5. 追加 auth-secret={authSecret}（authSecret 为签名值，即 SecretSign(secret, timestamp) 的结果）
//
// 待签字符串整体做 hex(md5(...)) 即得 auth-sign。
// 校验方持有 secret：先按同一规则复算签名值，再拼串 md5 与 auth-sign 比对，可完整复算。
func buildSignString(appid, authSecret, timestamp string, params map[string]any) string {
	var sb strings.Builder

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		switch v := params[k].(type) {
		case []string:
			vals := append([]string(nil), v...)
			sort.Strings(vals)
			for _, item := range vals {
				sb.WriteString(k + "=" + item + "&")
			}
		case []any:
			vals := make([]string, 0, len(v))
			for _, item := range v {
				vals = append(vals, fmt.Sprint(item))
			}
			sort.Strings(vals)
			for _, item := range vals {
				sb.WriteString(k + "=" + item + "&")
			}
		default:
			sb.WriteString(k + "=" + fmt.Sprint(v) + "&")
		}
	}

	sb.WriteString("auth-appid=" + appid + "&")
	sb.WriteString("auth-timestamp=" + timestamp + "&")
	sb.WriteString("auth-secret=" + authSecret)

	return sb.String()
}

// MD5Sign 生成请求头 auth-sign：hex(md5(待签字符串))。
// 待签字符串中的 auth-secret 取签名值（SecretSign 结果），与请求头 auth-secret 同源，
// 保证 IAM 校验方可复算；不拼明文密钥（原需求文档"末尾再追加密钥"经核对为误解，已按业界通用规则修正）。
func MD5Sign(appid, secret, timestamp string, params map[string]any) string {
	authSecret := SecretSign(secret, timestamp)
	return md5Hex(buildSignString(appid, authSecret, timestamp, params))
}
