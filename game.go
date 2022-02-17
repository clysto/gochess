package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	currentPiece uint8
	board        *Board
}

func NewGame() *Game {
	return &Game{
		board: NewBoard(),
	}
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		sq := g.GetClickSquare(x, y)
		if g.board.PieceAt(sq) != nil {
			g.currentPiece = sq
		}
	}
	return nil
}

func (g *Game) GetClickSquare(x, y int) uint8 {
	file := (x-40)/160 + 3
	rank := (y-40)/160 + 3
	if file >= len(Files) || rank >= len(Ranks) {
		return 0
	}
	if file < 3 || rank < 3 {
		return 0
	}
	square := Msb(And(&Files[file], &Ranks[rank]))
	return Squares180[square]
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	screen.DrawImage(BackgroundImage, nil)
	screen.DrawImage(BoardImage, nil)

	// draw pieces
	for _, sq := range Squares {
		if sq != 0 {
			piece := g.board.PieceAt(sq)
			if piece != nil {
				g.DrawPieceAt(screen, PieceImage(piece), sq)
			}
		}
	}

	// draw boxes
	if g.currentPiece > 0 {
		g.DrawPieceAt(screen, BlueBoxImage, g.currentPiece)
		for _, m := range g.board.PseudoLegalMoves(&BbSquares[g.currentPiece], &BbInBoard) {
			g.DrawPieceAt(screen, RedBoxImage, m.ToSquare)
		}
	}
}

func (g *Game) DrawImageAt(screen *ebiten.Image, image *ebiten.Image, x float64, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(image, op)
}

func (g *Game) DrawPieceAt(screen *ebiten.Image, image *ebiten.Image, sq uint8) {
	x := 40 + (SquareFile(Squares180[sq])-3)*160
	y := 40 + (SquareRank(Squares180[sq])-3)*160
	g.DrawImageAt(screen, image, float64(x), float64(y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1520, 1680
}
