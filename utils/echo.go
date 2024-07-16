package utils

import "net/url"

func PopQueryParam[T string | []string](qp url.Values, key string) T {
	var result T
	switch any(result).(type) {
	case string:
		result = any(qp.Get(key)).(T)
	case []string:
		result = any(qp[key]).(T)
	}
	qp.Del(key)
	return result
}
