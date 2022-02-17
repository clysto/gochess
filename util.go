package main

import (
	"fmt"
	"strings"

	"github.com/holiman/uint256"
)

func Uint256(hex string) *uint256.Int {
	hex = strings.ReplaceAll(hex, "_", "")
	if len(hex) > 3 {
		end := 0
		for i := 2; i < len(hex); i++ {
			if hex[i] == '0' {
				end++
			} else {
				break
			}
		}
		hex = hex[:2] + hex[2+end:]
	}
	n, err := uint256.FromHex(hex)
	if err != nil {
		panic(err)
	}
	return n
}

func Zero() *uint256.Int {
	return uint256.NewInt(0)
}

func Msb(x *uint256.Int) int {
	return x.BitLen() - 1
}

func ScanReversed(bb *uint256.Int) []uint8 {
	l := make([]uint8, 0)
	for !bb.IsZero() {
		r := bb.BitLen() - 1
		l = append(l, uint8(r))
		bb.Xor(bb, &BbSquares[r])
	}
	return l
}

func Lsh(x *uint256.Int, n uint) *uint256.Int {
	return Zero().Lsh(x, n)
}

func And(nums ...*uint256.Int) *uint256.Int {
	z := BbAll
	// TODO:下标有可能溢出
	for _, n := range nums {
		z.And(&z, n)
	}
	return &z
}

func Or(nums ...*uint256.Int) *uint256.Int {
	z := Zero()
	for _, n := range nums {
		z.Or(z, n)
	}
	return z
}

func BbPrint(bb *uint256.Int) {
	builder := strings.Builder{}
	for _, sq := range Squares180 {
		mask := &BbSquares[sq]
		if And(mask, &BbInBoard).IsZero() {
			continue
		}
		if And(bb, mask).IsZero() {
			fmt.Fprint(&builder, ".")
		} else {
			fmt.Fprint(&builder, "1")
		}
		if !And(mask, &FileI).IsZero() {
			fmt.Fprint(&builder, "\n")
		} else if sq != I0 {
			fmt.Fprint(&builder, " ")
		}
	}
	print(builder.String())
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
