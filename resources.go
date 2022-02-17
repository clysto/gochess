package main

import (
	"embed"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed assets
var f embed.FS

var (
	RedRookImage      *ebiten.Image
	RedKnightImage    *ebiten.Image
	RedBishopImage    *ebiten.Image
	RedAdvisorImage   *ebiten.Image
	RedKingImage      *ebiten.Image
	RedPawnImage      *ebiten.Image
	RedCannonImage    *ebiten.Image
	BlackRookImage    *ebiten.Image
	BlackKnightImage  *ebiten.Image
	BlackBishopImage  *ebiten.Image
	BlackAdvisorImage *ebiten.Image
	BlackKingImage    *ebiten.Image
	BlackPawnImage    *ebiten.Image
	BlackCannonImage  *ebiten.Image
	BackgroundImage   *ebiten.Image
	BoardImage        *ebiten.Image
	RedBoxImage       *ebiten.Image
	BlueBoxImage      *ebiten.Image
)

func init() {
	backgroundFile, err := f.Open("assets/bg.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer backgroundFile.Close()
	BackgroundImage, _, err = ebitenutil.NewImageFromReader(backgroundFile)
	if err != nil {
		log.Fatal(err)
	}

	boardFile, err := f.Open("assets/board.png")
	if err != nil {
		log.Fatal(err)
	}
	defer boardFile.Close()
	BoardImage, _, err = ebitenutil.NewImageFromReader(boardFile)
	if err != nil {
		log.Fatal(err)
	}

	pieceFile, err := f.Open("assets/piece.png")
	if err != nil {
		log.Fatal(err)
	}
	defer pieceFile.Close()
	allPieces, _, err := ebitenutil.NewImageFromReader(pieceFile)
	if err != nil {
		log.Fatal(err)
	}
	RedRookImage = allPieces.SubImage(image.Rect(0, 0, 160, 160)).(*ebiten.Image)
	RedKnightImage = allPieces.SubImage(image.Rect(0, 160*1, 160, 160*2)).(*ebiten.Image)
	RedBishopImage = allPieces.SubImage(image.Rect(0, 160*2, 160, 160*3)).(*ebiten.Image)
	RedAdvisorImage = allPieces.SubImage(image.Rect(0, 160*3, 160, 160*4)).(*ebiten.Image)
	RedKingImage = allPieces.SubImage(image.Rect(0, 160*4, 160, 160*5)).(*ebiten.Image)
	RedCannonImage = allPieces.SubImage(image.Rect(0, 160*5, 160, 160*6)).(*ebiten.Image)
	RedPawnImage = allPieces.SubImage(image.Rect(0, 160*6, 160, 160*7)).(*ebiten.Image)
	BlackRookImage = allPieces.SubImage(image.Rect(0, 160*7, 160, 160*8)).(*ebiten.Image)
	BlackKnightImage = allPieces.SubImage(image.Rect(0, 160*8, 160, 160*9)).(*ebiten.Image)
	BlackBishopImage = allPieces.SubImage(image.Rect(0, 160*9, 160, 160*10)).(*ebiten.Image)
	BlackAdvisorImage = allPieces.SubImage(image.Rect(0, 160*10, 160, 160*11)).(*ebiten.Image)
	BlackKingImage = allPieces.SubImage(image.Rect(0, 160*11, 160, 160*12)).(*ebiten.Image)
	BlackCannonImage = allPieces.SubImage(image.Rect(0, 160*12, 160, 160*13)).(*ebiten.Image)
	BlackPawnImage = allPieces.SubImage(image.Rect(0, 160*13, 160, 160*14)).(*ebiten.Image)
	RedBoxImage = allPieces.SubImage(image.Rect(0, 160*14, 160, 160*15)).(*ebiten.Image)
	BlueBoxImage = allPieces.SubImage(image.Rect(0, 160*15, 160, 160*16)).(*ebiten.Image)
}

func PieceImage(piece *Piece) *ebiten.Image {
	if piece.Color {
		switch piece.PieceType {
		case Pawn:
			return RedPawnImage
		case Cannon:
			return RedCannonImage
		case Rook:
			return RedRookImage
		case Knight:
			return RedKnightImage
		case Bishop:
			return RedBishopImage
		case Advisor:
			return RedAdvisorImage
		case King:
			return RedKingImage
		}
	} else {
		switch piece.PieceType {
		case Pawn:
			return BlackPawnImage
		case Cannon:
			return BlackCannonImage
		case Rook:
			return BlackRookImage
		case Knight:
			return BlackKnightImage
		case Bishop:
			return BlackBishopImage
		case Advisor:
			return BlackAdvisorImage
		case King:
			return BlackKingImage
		}
	}
	return nil
}
