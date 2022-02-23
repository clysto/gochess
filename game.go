package main

import (
	"github.com/clysto/gochess/chess"
	"github.com/clysto/gochess/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Game struct {
	fromSquare uint8
	board      *chess.Board
}

func NewGame() *Game {
	return &Game{
		board: chess.NewBoard(),

	}
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		sq := g.GetClickSquare(x, y)
		piece := g.board.PieceAt(sq)
		if piece != nil && piece.Color == g.board.Turn() {
			g.fromSquare = sq
		} else if (piece == nil || piece.Color != g.board.Turn()) && g.fromSquare > 0 {
			move := &chess.Move{FromSquare: g.fromSquare, ToSquare: sq}
			if g.board.IsPseudoLegal(move) {
				g.board.Push(move)
				g.fromSquare = 0
			}
		}
	}
	return nil
}

func (g *Game) GetClickSquare(x, y int) uint8 {
	file := (x-40)/160 + 3
	rank := (y-40)/160 + 3
	if file >= 12 || rank >= 13 {
		return 0
	}
	if file < 3 || rank < 3 {
		return 0
	}
	square := chess.Msb(chess.And(&chess.Files[file], &chess.Ranks[rank]))
	return chess.Squares180[square]
}

func (g *Game) Draw(screen *ebiten.Image) {
	boardImage := ebiten.NewImage(1520, 1680)
	boardImage.Fill(color.White)
	boardImage.DrawImage(resources.BackgroundImage, nil)
	boardImage.DrawImage(resources.BoardImage, nil)

	// draw pieces
	for _, sq := range chess.Squares {
		if sq != 0 {
			piece := g.board.PieceAt(sq)
			if piece != nil {
				g.DrawPieceAt(boardImage, resources.PieceImage(piece), sq)
			}
		}
	}

	// draw boxes
	if g.fromSquare > 0 {
		g.DrawPieceAt(boardImage, resources.BlueBoxImage, g.fromSquare)
		for _, m := range g.board.PseudoLegalMoves(&chess.BbSquares[g.fromSquare], &chess.BbInBoard) {
			g.DrawPieceAt(boardImage, resources.RedBoxImage, m.ToSquare)
		}
	}

	screen.DrawImage(boardImage, nil)
	bottomBarImage := ebiten.NewImage(1520, 150)
	bottomBarImage.Fill(color.White)
	g.DrawImageAt(screen, bottomBarImage, 0, 1680)
}

func (g *Game) DrawImageAt(screen *ebiten.Image, image *ebiten.Image, x float64, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(image, op)
}

func (g *Game) DrawPieceAt(screen *ebiten.Image, image *ebiten.Image, sq uint8) {
	x := 40 + (chess.SquareFile(chess.Squares180[sq])-3)*160
	y := 40 + (chess.SquareRank(chess.Squares180[sq])-3)*160
	g.DrawImageAt(screen, image, float64(x), float64(y))
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 1520, 1680 + 150
}
