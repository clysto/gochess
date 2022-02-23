package chess

import (
	_ "embed"
	"github.com/holiman/uint256"
)

var MovesTable struct {
	BbKnightMasks       [256]uint256.Int
	BbKnightAttacks     [256]map[uint256.Int]uint256.Int
	BbBishopMasks       [256]uint256.Int
	BbBishopAttacks     [256]map[uint256.Int]uint256.Int
	BbPawnAttacks       map[bool][256]uint256.Int
	BbKingAttacks       [256]uint256.Int
	BbAdvisorAttacks    [256]uint256.Int
	BbRankMasks         [256]uint256.Int
	BbRankAttacks       [256]map[uint256.Int]uint256.Int
	BbFileMasks         [256]uint256.Int
	BbFileAttacks       [256]map[uint256.Int]uint256.Int
	BbCannonRankMasks   [256]uint256.Int
	BbCannonRankAttacks [256]map[uint256.Int]uint256.Int
	BbCannonFileMasks   [256]uint256.Int
	BbCannonFileAttacks [256]map[uint256.Int]uint256.Int
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

func jumpAttacks(square uint8, occupied *uint256.Int, deltas []int) *uint256.Int {
	attacks := BbEmpty.Clone()
	for _, delta := range deltas {
		hops := 0
		sq := int(square)
		for {
			sq += delta
			if !(sq >= 0 && sq < 256) || SquareDistance(uint8(sq), uint8(sq-delta)) > 2 {
				break
			}

			if !And(occupied, &BbSquares[sq]).IsZero() {
				if hops == 1 {
					attacks.Or(attacks, &BbSquares[sq])
					break
				} else {
					hops++
				}
			}

		}
	}
	return attacks
}

func carryRippler(mask *uint256.Int, f func(uint256.Int)) {
	// Carry-Rippler trick to iterate subsets of mask.
	subset := BbEmpty
	for {
		f(subset)
		subset = *And(subset.Sub(&subset, mask), mask)
		if subset.IsZero() {
			break
		}
	}
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

func genPawnAttacks() {
	MovesTable.BbPawnAttacks = map[bool][256]uint256.Int{}
	var attacks [256]uint256.Int
	for _, sq := range Squares {
		if sq > I4 {
			attacks[sq] = *stepAttacks(sq, []int{-1, 16, 1})
		} else {
			attacks[sq] = *stepAttacks(sq, []int{16})
		}
	}
	MovesTable.BbPawnAttacks[Red] = attacks
	for _, sq := range Squares {
		if sq < A5 {
			attacks[sq] = *stepAttacks(sq, []int{-1, -16, 1})
		} else {
			attacks[sq] = *stepAttacks(sq, []int{-16})
		}
	}
	MovesTable.BbPawnAttacks[Black] = attacks
}

func genKingAttacks() {
	for _, sq := range Squares {
		if !And(&BbSquares[sq], &BbInPalace).IsZero() {
			MovesTable.BbKingAttacks[sq] = *And(stepAttacks(sq, []int{-16, 16, 1, -1}), &BbInPalace)
		}
	}
}

func genAdvisorAttacks() {
	for _, sq := range Squares {
		if !And(&BbSquares[sq], &BbSquaresAdvisor).IsZero() {
			MovesTable.BbAdvisorAttacks[sq] = *And(stepAttacks(sq, []int{15, 17, -15, -17}), &BbInPalace)
		}
	}
}

func attackTable(deltas []int, masks *[256]uint256.Int, attacks *[256]map[uint256.Int]uint256.Int, jump bool) {
	for _, sq := range Squares {
		mask := And(slidingAttacks(sq, &BbEmpty, deltas), &BbInBoard)
		if !jump {
			mask.And(mask, Not(Edges(sq)))
		}
		attacks[sq] = map[uint256.Int]uint256.Int{}
		carryRippler(mask, func(subset uint256.Int) {
			if jump {
				attacks[sq][subset] = *Or(
					jumpAttacks(sq, &subset, deltas),
					And(slidingAttacks(sq, &subset, deltas), Not(&subset)),
				)
			} else {
				attacks[sq][subset] = *slidingAttacks(sq, &subset, deltas)
			}
		})
		masks[sq] = *mask
	}
}

func init() {
	genKnightAttacks()
	genBishopAttacks()
	genPawnAttacks()
	genKingAttacks()
	genAdvisorAttacks()
	attackTable([]int{-1, 1}, &MovesTable.BbRankMasks, &MovesTable.BbRankAttacks, false)
	attackTable([]int{-16, 16}, &MovesTable.BbFileMasks, &MovesTable.BbFileAttacks, false)
	attackTable([]int{-1, 1}, &MovesTable.BbCannonRankMasks, &MovesTable.BbCannonRankAttacks, true)
	attackTable([]int{-16, 16}, &MovesTable.BbCannonFileMasks, &MovesTable.BbCannonFileAttacks, true)
}
