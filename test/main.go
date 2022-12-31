package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"

	"github.com/unsafe-risk/utilx/typex/convx"
)

type HexString string

func main() {
	fmt.Println(convx.IntoOr(int64(123456), ""))
	fmt.Println(convx.IntoOr(true, ""))
	fmt.Println(convx.IntoOr("5678", int32(0)))

	convx.Register(func(b []byte) HexString {
		return HexString(hex.EncodeToString(b))
	})

	buf := [32]byte{}
	rand.Read(buf[:])
	fmt.Println(convx.IntoOr(buf[:], HexString("")))

	v, ok := convx.Into[int64, bool](1)
	fmt.Println(v, ok)
}
