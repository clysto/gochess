package chess

import (
	"fmt"
	"testing"
)

func TestConstants(t *testing.T) {
	fmt.Println(BbSquares[A0])
	fmt.Println(BbSquares[B0])
	BbPrint(&BbA0)
}

func TestBbKnightMasks(t *testing.T) {
	BbPrint(&MovesTable.BbBishopMasks[C0])
	fmt.Println("-----------------")
	attacks := MovesTable.BbBishopAttacks[C0][BbD1]
	BbPrint(&attacks)
}

func TestAttacksMask(t *testing.T) {
	b := NewBoard()
	mask := b.AttacksMask(C0)
	BbPrint(mask)
}

func TestPseudoLegalMoves(t *testing.T) {
	b := NewBoard()
	moves := b.PseudoLegalMoves(&BbB0, &BbInBoard)
	for _, move := range moves {
		fmt.Println(move)
	}
}
