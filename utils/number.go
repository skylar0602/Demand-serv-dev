package utils

import (
	"fmt"
	"math/big"
	"strings"
)

func String2BigFloat(s string) (*big.Float, error) {
	f, _, err := big.ParseFloat(s, 10, 0, big.ToNearestEven)
	if err != nil {
		fmt.Printf("Error converting string to big.Float: %v\n", err)
		return nil, err
	}

	return f, nil
}

func Float2String(f float64) string {
	// 将浮点数转换为字符串（不使用科学记数法）
	s := fmt.Sprintf("%f", f)
	// 去除字符串末尾的零
	s = strings.TrimRight(s, "0")
	// 如果字符串以小数点结尾，则也去除小数点
	s = strings.TrimRight(s, ".")
	return s
}
