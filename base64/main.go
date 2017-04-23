package main

import (
	"encoding/base64"
	"github.com/gopherjs/gopherjs/js"
	"math/big"
)

func main() {
	js.Global.Set("converter", map[string]interface{}{
		"base64ToBase62": Base64ToBase62,
	})
}

func toBase62()

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func Base64ToBase62(input string) string {
	sDec, _ := base64.StdEncoding.DecodeString(input)
	return ConvertTo62(sDec)
}

func ConvertTo62(raw []byte) string {
	bi := big.Int{}
	bi.SetBytes(raw)
	rem := big.NewInt(0)
	base := big.NewInt(62)
	zero := big.NewInt(0)
	result := ""

	for bi.Cmp(zero) > 0 {
		_, rem = bi.DivMod(&bi, base, rem)
		result += string(alphabet[int(rem.Uint64())])
	}

	for len(result) < 22 {
		result += "0"
	}
	return reverse(result)
}
