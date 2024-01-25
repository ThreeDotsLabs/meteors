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
var PlayerSprite = MustLoadImage("img/Ships/spaceShips_007.png")
var ShieldSprite = MustLoadImage("img/Ships/shield.png")
var MeteorSprites = mustLoadImages("img/Meteors/*.png")
var HighSpeedFollowPlayerEnemySprite = MustLoadImage("img/Ships/spaceShips_006.png")
var LowSpeedEnemyLightMissile = MustLoadImage("img/Ships/spaceShips_005.png")
var LowSpeedEnemyAutoLightMissile = MustLoadImage("img/Ships/spaceShips_003.png")

// weapons
var MissileSprite = MustLoadImage("img/Missiles/spaceMissiles_016.png")
var MissileSpriteImage = LoadAsImage("img/Missiles/spaceMissiles_016.png")
var DoubleMissileSprite = MustLoadImage("img/Missiles/spaceMissiles_010.png")
var EnemyLightMissile = MustLoadImage("img/Missiles/spaceMissiles_015.png")
var EnemyAutoLightMissile = MustLoadImage("img/Missiles/spaceMissiles_018.png")
var LaserCanon = MustLoadImage("img/Effects/spaceEffects_004.png")
var PentaLaser = MustLoadImage("img/Effects/spaceEffects_002.png")
var DoubleLaserCanon = MustLoadImage("img/Effects/spaceEffects_018.png")
var LaserBeam = MustLoadImage("img/Effects/spaceEffects_006.png")
var ClusterMines = MustLoadImage("img/Building/spaceBuilding_008.png")
var BigBomb = MustLoadImage("img/Building/spaceBuilding_004.png")
var MachineGun = MustLoadImage("img/Missiles/spaceMissiles_038.png")
var DoubleMachineGun = MustLoadImage("img/Missiles/spaceMissiles_025.png")
var PlasmaGun = MustLoadImage("img/Effects/plasma_gun2.png")

// backgrounds
var FirstLevelBg = MustLoadImage("img/bg1.png")

// spritesheets
var EnemyBlowSpriteSheet = MustLoadImage("img/Effects/blow.png")
var BigBlowSpriteSheet = MustLoadImage("img/Effects/bigblow.png")
var PlayerFireburstSpriteSheet = MustLoadImage("img/Effects/fire.png")
var ShieldSpriteSheet = MustLoadImage("img/Ships/shieldSpriteSheet1.png")
var ProjectileBlowSpriteSheet = MustLoadImage("img/Effects/projectile_blow.png")
var LightMissileBlowSpriteSheet = MustLoadImage("img/Effects/light_missile_blow2.png")
var ClusterMinesBlowSpriteSheet = MustLoadImage("img/Effects/cluster_mines_blow.png")
var PlasmaGunProjectileSpriteSheet = MustLoadImage("img/Effects/plasma_gun1.png")

// fonts
var ScoreFont1024x768 = mustLoadFont("fonts/Kenney Pixel.ttf", 48)
var InfoFont1024x768 = mustLoadFont("fonts/Kenney Pixel.ttf", 32)
var SmallFont1024x768 = mustLoadFont("fonts/Kenney Mini.ttf", 18)
var ProfileFont1024x768 = mustLoadFont("fonts/Kenney Future.ttf", 12)
var ProfileBigFont1024x768 = mustLoadFont("fonts/Kenney Future.ttf", 24)

var ScoreFont1920x1080 = mustLoadFont("fonts/Kenney Pixel.ttf", 96)
var InfoFont1920x1080 = mustLoadFont("fonts/Kenney Pixel.ttf", 64)
var SmallFont1920x1080 = mustLoadFont("fonts/Kenney Mini.ttf", 36)
var ProfileFont1920x1080 = mustLoadFont("fonts/Kenney Future.ttf", 24)
var ProfileBigFont1920x1080 = mustLoadFont("fonts/Kenney Future.ttf", 48)

func LoadAsImage(name string) image.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return img
}

func MustLoadImage(name string) *ebiten.Image {
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
		images[i] = MustLoadImage(match)
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
