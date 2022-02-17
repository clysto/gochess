package main

import (
	"github.com/holiman/uint256"
)

type Board struct {
	pawns         *uint256.Int
	knights       *uint256.Int
	bishops       *uint256.Int
	rooks         *uint256.Int
	cannons       *uint256.Int
	advisors      *uint256.Int
	kings         *uint256.Int
	occupied      *uint256.Int
	occupiedColor map[bool]*uint256.Int
	turn          bool
}

type Piece struct {
	Color     bool
	PieceType uint8
}

type Move struct {
	FromSquare uint8
	ToSquare   uint8
}

var (
	BbCorners    = *Or(&BbA0, &BbI0, &BbA9, &BbI9)
	BbRedPawns   = *Or(&BbA3, &BbC3, &BbE3, &BbG3, &BbI3)
	BbBlackPawns = *Or(&BbA6, &BbC6, &BbE6, &BbG6, &BbI6)
)

const (
	Pawn uint8 = iota + 1
	Cannon
	Rook
	Knight
	Bishop
	Advisor
	King
)

const (
	Red   = true
	Black = false
)

func SquareMirror(square uint8) uint8 {
	return square ^ 0xf0
}

func SquareFile(square uint8) int {
	return int(square & 0xf)
}

func SquareRank(square uint8) int {
	return int(square >> 4)
}

func SquareInBoard(square uint8) bool {
	mask := &BbSquares[square]
	if mask == nil {
		return false
	}
	return !And(mask, &BbInBoard).IsZero()
}

func SquareDistance(a uint8, b uint8) int {
	return Max(Abs(SquareFile(a)-SquareFile(b)), Abs(SquareRank(a)-SquareRank(b)))
}

func NewBoard() *Board {
	b := Board{}
	b.pawns = Or(&BbRedPawns, &BbBlackPawns)
	b.knights = Or(&BbB0, &BbH0, &BbB9, &BbH9)
	b.bishops = Or(&BbC0, &BbG0, &BbC9, &BbG9)
	b.rooks = BbCorners.Clone()
	b.cannons = Or(&BbB2, &BbH2, &BbB7, &BbH7)
	b.advisors = Or(&BbD0, &BbF0, &BbD9, &BbF9)
	b.kings = Or(&BbE0, &BbE9)
	b.occupiedColor = map[bool]*uint256.Int{
		Red:   Or(&Rank0, &BbB2, &BbH2, &BbRedPawns),
		Black: Or(&Rank9, &BbB7, &BbH7, &BbBlackPawns),
	}
	b.occupied = Or(b.occupiedColor[Red], b.occupiedColor[Black])
	b.turn = Red
	return &b
}

func (b *Board) PieceAt(sq uint8) *Piece {
	t := b.PieceTypeAt(sq)
	if t > 0 {
		mask := &BbSquares[sq]
		color := !And(b.occupiedColor[Red], mask).IsZero()
		return &Piece{
			PieceType: t,
			Color:     color,
		}
	}
	return nil
}

func (b *Board) PieceTypeAt(sq uint8) uint8 {
	mask := &BbSquares[sq]
	if mask == nil || And(b.occupied, mask).IsZero() {
		return 0
	} else if !And(b.pawns, mask).IsZero() {
		return Pawn
	} else if !And(b.knights, mask).IsZero() {
		return Knight
	} else if !And(b.bishops, mask).IsZero() {
		return Bishop
	} else if !And(b.rooks, mask).IsZero() {
		return Rook
	} else if !And(b.cannons, mask).IsZero() {
		return Cannon
	} else if !And(b.advisors, mask).IsZero() {
		return Advisor
	} else {
		return King
	}
}

func (b *Board) AttacksMask(square uint8) *uint256.Int {
	bbSquare := &BbSquares[square]
	mask := BbEmpty

	if !And(bbSquare, b.knights).IsZero() {
		mask = MovesTable.BbKnightAttacks[square][*And(
			&MovesTable.BbKnightMasks[square], b.occupied,
		)]
	} else if !And(bbSquare, b.bishops).IsZero() {
		mask = MovesTable.BbBishopAttacks[square][*And(
			&MovesTable.BbBishopMasks[square], b.occupied,
		)]
	}

	return &mask
}

func (b *Board) PseudoLegalMoves(fromMask *uint256.Int, toMask *uint256.Int) []*Move {
	ourPieces := b.occupiedColor[b.turn]
	fromSquares := And(ourPieces, fromMask)
	moves := make([]*Move, 0)
	for _, fromSquare := range ScanReversed(fromSquares) {
		moveMask := b.AttacksMask(fromSquare)
		moveMask.And(moveMask, Zero().Not(ourPieces))
		moveMask.And(moveMask, toMask)
		for _, toSquare := range ScanReversed(moveMask) {
			moves = append(moves, &Move{
				FromSquare: fromSquare,
				ToSquare:   toSquare,
			})
		}
	}
	return moves
}
