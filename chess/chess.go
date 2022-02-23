package chess

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
	BbInPalace   = *Or(&BbD0, &BbE0, &BbF0, &BbD1, &BbE1, &BbF1, &BbD2, &BbE2, &BbF2,
		&BbD7, &BbE7, &BbF7, &BbD8, &BbE8, &BbF8, &BbD9, &BbE9, &BbF9)
	BbSquaresAdvisor = *Or(&BbD0, &BbF0, &BbE1, &BbD2, &BbF2, &BbD7, &BbF7, &BbE8, &BbD9, &BbF9)
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

func Edges(square uint8) *uint256.Int {
	return Or(
		And(Or(&Rank0, &Rank9), Not(&Ranks[SquareRank(square)])),
		And(Or(&FileA, &FileI), Not(&Files[SquareFile(square)])),
	)
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

func (b *Board) Turn() bool {
	return b.turn
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

func (b *Board) PieceTypeAt(square uint8) uint8 {
	mask := &BbSquares[square]
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
	} else if !And(bbSquare, b.pawns).IsZero() {
		color := !And(bbSquare, b.occupiedColor[Red]).IsZero()
		mask = MovesTable.BbPawnAttacks[color][square]
	} else if !And(bbSquare, b.kings).IsZero() {
		mask = MovesTable.BbKingAttacks[square]
	} else if !And(bbSquare, b.advisors).IsZero() {
		mask = MovesTable.BbAdvisorAttacks[square]
	} else if !And(bbSquare, b.rooks).IsZero() {
		rankMask := MovesTable.BbFileAttacks[square][*And(&MovesTable.BbFileMasks[square], b.occupied)]
		fileMask := MovesTable.BbRankAttacks[square][*And(&MovesTable.BbRankMasks[square], b.occupied)]
		mask = *Or(&rankMask, &fileMask)
	} else if !And(bbSquare, b.cannons).IsZero() {
		rankMask := MovesTable.BbCannonFileAttacks[square][*And(&MovesTable.BbCannonFileMasks[square], b.occupied)]
		fileMask := MovesTable.BbCannonRankAttacks[square][*And(&MovesTable.BbCannonRankMasks[square], b.occupied)]
		mask = *Or(&rankMask, &fileMask)
	}

	return &mask
}

func (b *Board) PseudoLegalMoves(fromMask *uint256.Int, toMask *uint256.Int) []*Move {
	ourPieces := b.occupiedColor[b.turn]
	fromSquares := And(ourPieces, fromMask)
	moves := make([]*Move, 0)
	for _, fromSquare := range ScanReversed(fromSquares) {
		moveMask := b.AttacksMask(fromSquare)
		moveMask.And(moveMask, Not(ourPieces))
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

func (b *Board) removePieceAt(square uint8) uint8 {
	pieceType := b.PieceTypeAt(square)
	mask := &BbSquares[square]

	switch pieceType {
	case Pawn:
		b.pawns.Xor(b.pawns, mask)
	case Knight:
		b.knights.Xor(b.knights, mask)
	case Bishop:
		b.bishops.Xor(b.bishops, mask)
	case Rook:
		b.rooks.Xor(b.rooks, mask)
	case Cannon:
		b.cannons.Xor(b.cannons, mask)
	case King:
		b.kings.Xor(b.kings, mask)
	case Advisor:
		b.advisors.Xor(b.advisors, mask)
	default:
		return 0
	}
	b.occupied.Xor(b.occupied, mask)
	b.occupiedColor[Red].And(b.occupiedColor[Red], Not(mask))
	b.occupiedColor[Black].And(b.occupiedColor[Black], Not(mask))
	return pieceType
}

func (b *Board) setPieceAt(square uint8, pieceType uint8, color bool) {
	b.removePieceAt(square)
	mask := &BbSquares[square]

	switch pieceType {
	case Pawn:
		b.pawns.Or(b.pawns, mask)
	case Knight:
		b.knights.Or(b.knights, mask)
	case Bishop:
		b.bishops.Or(b.bishops, mask)
	case Rook:
		b.rooks.Or(b.rooks, mask)
	case Cannon:
		b.cannons.Or(b.cannons, mask)
	case King:
		b.kings.Or(b.kings, mask)
	case Advisor:
		b.advisors.Or(b.advisors, mask)
	default:
		return
	}

	b.occupied.Xor(b.occupied, mask)
	b.occupiedColor[color].Xor(b.occupiedColor[color], mask)
}

func (b *Board) Push(move *Move) {
	pieceType := b.removePieceAt(move.FromSquare)
	b.setPieceAt(move.ToSquare, pieceType, b.turn)
	b.turn = !b.turn
}

func (b *Board) IsPseudoLegal(move *Move) bool {
	if move == nil {
		return false
	}
	piece := b.PieceAt(move.FromSquare)
	if piece == nil {
		return false
	}
	fromMask := &BbSquares[move.FromSquare]
	toMask := &BbSquares[move.ToSquare]
	// 是否是自己的棋子
	if And(b.occupiedColor[b.turn], fromMask).IsZero() {
		return false
	}
	// 目标格子不能有自己的棋子
	if !And(b.occupiedColor[b.turn], toMask).IsZero() {
		return false
	}

	return !And(b.AttacksMask(move.FromSquare), toMask).IsZero()
}
