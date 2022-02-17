package main

import (
	"github.com/holiman/uint256"
)

var MovesTable struct {
	BbKnightMasks   [256]uint256.Int
	BbKnightAttacks [256]map[uint256.Int]uint256.Int
	BbBishopMasks   [256]uint256.Int
	BbBishopAttacks [256]map[uint256.Int]uint256.Int
}

func stepAttacks(square uint8, deltas []int) *uint256.Int {
	return slidingAttacks(square, &BbAll, deltas)
}

func slidingAttacks(square uint8, occupied *uint256.Int, deltas []int) *uint256.Int {
	attacks := BbEmpty.Clone()
	for _, delta := range deltas {
		sq := int(square)
		for {
			sq += delta
			if !(sq >= 0 && sq < 256) || SquareDistance(uint8(sq), uint8(sq-delta)) > 2 {
				break
			}

			attacks.Or(attacks, &BbSquares[sq])

			if !And(occupied, &BbSquares[sq]).IsZero() {
				break
			}
		}
	}
	return attacks
}

func genKnightAttacks() {
	knightDeltas := []int{33, 31, -14, 18, -33, -31, -18, 14}
	directions := []int{16, 1, -16, -1}
	for k, square := range Squares {
		sq := int(square)
		if !SquareInBoard(square) {
			continue
		}
		MovesTable.BbKnightMasks[k] = BbEmpty
		for _, d := range directions {
			MovesTable.BbKnightMasks[k].Or(&MovesTable.BbKnightMasks[k], &BbSquares[sq+d])
		}
		MovesTable.BbKnightAttacks[k] = map[uint256.Int]uint256.Int{}
		// 马脚位置有16种情况
		for i := 0; i <= 0xf; i++ {
			subset := BbEmpty.Clone()
			var deltas []int
			for j := 0; j < 4; j++ {
				if i>>j&1 != 0 {
					// 别马脚
					subset.Or(subset, &BbSquares[sq+directions[j]])
				} else {
					deltas = append(deltas, knightDeltas[2*j])
					deltas = append(deltas, knightDeltas[2*j+1])
				}
			}
			MovesTable.BbKnightAttacks[k][*subset] = *stepAttacks(square, deltas)
		}
	}
}

func genBishopAttacks() {
	directions := []int{15, 17, -15, -17}
	for k, square := range Squares {
		sq := int(square)
		if !SquareInBoard(square) {
			continue
		}
		var squareSide *uint256.Int
		if And(&BbSquares[square], &BbRedSide).IsZero() {
			squareSide = &BbBlackSide
		} else {
			squareSide = &BbRedSide
		}
		MovesTable.BbBishopMasks[k] = BbEmpty
		for _, d := range directions {
			MovesTable.BbBishopMasks[k].Or(&MovesTable.BbBishopMasks[k], &BbSquares[sq+d])
		}
		MovesTable.BbBishopAttacks[k] = map[uint256.Int]uint256.Int{}
		// 象眼位置有16种情况
		for i := 0; i <= 0xf; i++ {
			subset := BbEmpty.Clone()
			var deltas []int
			for j := 0; j < 4; j++ {
				if i>>j&1 != 0 {
					// 塞象眼
					subset.Or(subset, &BbSquares[sq+directions[j]])
				} else {
					deltas = append(deltas, 2*directions[j])
				}
			}
			MovesTable.BbBishopAttacks[k][*subset] = *And(stepAttacks(square, deltas), squareSide)
		}
	}
}

func init() {
	genKnightAttacks()
	genBishopAttacks()
}
