package protocol

import (
	"fmt"
)

func EncodeSimple(s string) []byte {
	return fmt.Appendf(nil, "+%s\r\n", s)
}

func EncodeBulk(s string) []byte {
	if len(s) == 0 {
		return []byte("$-1\r\n")
	}

	return fmt.Appendf(nil, "$%d\r\n%s\r\n", len(s), s)
}

func EncodeMultiBulk(s []string) []byte {
	if s == nil {
		return []byte("*-1\r\n")
	}

	if len(s) == 0 {
		return []byte("*0\r\n")
	}

	b := fmt.Appendf(nil, "*%d\r\n", len(s))

	for _, v := range s {
		b = append(b, EncodeBulk(v)...)
	}

	return b
}

func EncodeArray(arr []any) []byte {
	var buf []byte

	buf = append(buf, []byte(fmt.Sprintf("*%d\r\n", len(arr)))...)

	for _, v := range arr {
		switch x := v.(type) {

		case nil:
			buf = append(buf, []byte("$-1\r\n")...)

		case string:
			buf = append(buf, EncodeBulk(x)...)

		case int:
			buf = append(buf, EncodeInteger(x)...)

		case int64:
			buf = append(buf, EncodeInteger(int(x))...)

		case []any:
			buf = append(buf, EncodeArray(x)...)

		case []string:
			buf = append(buf, EncodeMultiBulk(x)...)

		default:
			panic(fmt.Sprintf("unsupported RESP type: %T", v))
		}
	}

	return buf
}

func EncodeError(s string) []byte {
	return fmt.Appendf(nil, "-%s\r\n", s)
}

func EncodeInteger(i int) []byte {
	return fmt.Appendf(nil, ":%d\r\n", i)
}
