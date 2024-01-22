package assets

import (
	"embed"
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed *
var assets embed.FS

// objects
var PlayerSprite = mustLoadImage("img/Ships/spaceShips_007.png")
var ShieldSprite = mustLoadImage("img/Ships/shield.png")
var HighSpeedFollowPlayerEnemySprite = mustLoadImage("img/Ships/spaceShips_006.png")
var LowSpeedEnemyLightMissile = mustLoadImage("img/Ships/spaceShips_005.png")
var LowSpeedEnemyAutoLightMissile = mustLoadImage("img/Ships/spaceShips_003.png")
var MeteorSprites = mustLoadImages("img/Meteors/*.png")

// weapons
var MissileSprite = mustLoadImage("img/Missiles/spaceMissiles_016.png")
var DoubleMissileSprite = mustLoadImage("img/Missiles/spaceMissiles_010.png")
var EnemyLightMissile = mustLoadImage("img/Missiles/spaceMissiles_015.png")
var EnemyAutoLightMissile = mustLoadImage("img/Missiles/spaceMissiles_018.png")
var LaserCannon = mustLoadImage("img/Effects/spaceEffects_004.png")
var LaserBeam = mustLoadImage("img/Effects/spaceEffects_006.png")
var ClusterMines = mustLoadImage("img/Building/spaceBuilding_006.png")
var BigBomb = mustLoadImage("img/Building/spaceBuilding_003.png")
var MachineGun = mustLoadImage("img/Missiles/spaceMissiles_038.png")

// backgrounds
var FirstLevelBg = mustLoadImage("img/bg1.png")

// spritesheets
var EnemyBlowSpriteSheet = mustLoadImage("img/Effects/blow.png")
var BigBlowSpriteSheet = mustLoadImage("img/Effects/bigblow.png")
var PlayerFireburstSpriteSheet = mustLoadImage("img/Effects/fire.png")
var ShieldSpriteSheet = mustLoadImage("img/Ships/shieldSpriteSheet1.png")
var ProjectileBlowSpriteSheet = mustLoadImage("img/Effects/projectile_blow1.png")

// fonts
var ScoreFont = mustLoadFont("fonts/Kenney Pixel.ttf", 48)
var InfoFont = mustLoadFont("fonts/Kenney Pixel.ttf", 32)
var SmallFont = mustLoadFont("fonts/Kenney Mini.ttf", 18)

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		panic(err)
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}

	return images
}

func mustLoadFont(name string, size float64) font.Face {
	f, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	tt, err := opentype.Parse(f)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	return face
}
