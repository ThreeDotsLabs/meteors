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
var PlayerSprite = MustLoadImage("img/Ships/ship4.png")
var ShieldSprite = MustLoadImage("img/Ships/shield.png")
var MeteorSprites = mustLoadImages("img/Meteors/*.png")
var Heal = MustLoadImage("img/Parts/heal1.png")
var HighSpeedFollowPlayerEnemySprite = MustLoadImage("img/Ships/enemy5.png")
var LowSpeedEnemyLightMissile = MustLoadImage("img/Ships/enemy4.png")
var LowSpeedEnemyAutoLightMissile = MustLoadImage("img/Ships/enemy3.png")

// items
var ItemMissileSprite = MustLoadImage("img/Items/missile_item.png")
var ItemDoubleMissileSprite = MustLoadImage("img/Items/double_missile_item.png")
var ItemLaserCanonSprite = MustLoadImage("img/Items/laser_cannon_item.png")
var ItemDoubleLaserCanonSprite = MustLoadImage("img/Items/double_laser_cannon_item.png")
var ItemMachineGun = MustLoadImage("img/Items/machine_gun_item.png")
var ItemDoubleMachineGun = MustLoadImage("img/Items/double_machine_gun_item.png")
var ItemPlasmaGun = MustLoadImage("img/Items/plasma_gun_item.png")
var ItemDoublePlasmaGun = MustLoadImage("img/Items/double_plasma_gun_item.png")
var ItemClusterMines = MustLoadImage("img/Items/cluster_mines_item.png")

// enemies
var Enemy1 = MustLoadImage("img/Ships/enemy1.png")
var Enemy2 = MustLoadImage("img/Ships/enemy2.png")
var Enemy3 = MustLoadImage("img/Ships/enemy3.png")
var Enemy4 = MustLoadImage("img/Ships/enemy4.png")
var Enemy5 = MustLoadImage("img/Ships/enemy5.png")
var Enemy6 = MustLoadImage("img/Ships/enemy6.png")
var Enemy7 = MustLoadImage("img/Ships/enemy7.png")
var Enemy8 = MustLoadImage("img/Ships/enemy8.png")
var Enemy9 = MustLoadImage("img/Ships/enemy9.png")

// weapons projectiles
var MissileSprite = MustLoadImage("img/Missiles/player_missile_projectile3.png")
var DoubleMissileSprite = MustLoadImage("img/Missiles/player_missile_projectile3.png")
var LaserCanon = MustLoadImage("img/Items/laser_cannon_item.png")
var DoubleLaserCanon = MustLoadImage("img/Items/double_laser_cannon_item.png")
var MachineGun = MustLoadImage("img/Missiles/player_bullet_projectile1.png")
var DoubleMachineGun = MustLoadImage("img/Missiles/player_bullet_projectile1.png")
var PlasmaGun = MustLoadImage("img/Effects/plasma_gun2.png")
var ClusterMines = MustLoadImage("img/Missiles/cluster_mines_projectile.png")
var BigBomb = MustLoadImage("img/Building/spaceBuilding_004.png")
var PentaLaser = MustLoadImage("img/Effects/spaceEffects_002.png")

// enemy weapons projectiles
var EnemyLightMissile = MustLoadImage("img/Missiles/missile_projectile1.png")
var EnemyAutoLightMissile = MustLoadImage("img/Missiles/missile_projectile2.png")
var EnemyHeavyMissile = MustLoadImage("img/Missiles/missile_projectile4.png")
var EnemyGunProjevtile = MustLoadImage("img/Missiles/missile_projectile5.png")
var EnemyBullet = MustLoadImage("img/Missiles/bullet_projectile1.png")
var EnemyMidMissile = MustLoadImage("img/Missiles/missile_projectile3.png")

// backgrounds
var Backgrounds = mustLoadImages("img/Backgrounds/*.png")

// spritesheets
var EnemyBlowSpriteSheet = MustLoadImage("img/Effects/blow.png")
var BigBlowSpriteSheet = MustLoadImage("img/Effects/bigblow.png")
var PlayerFireburstSpriteSheet = MustLoadImage("img/Effects/fire1.png")
var ShieldSpriteSheet = MustLoadImage("img/Ships/shieldSpriteSheet1.png")
var ProjectileBlowSpriteSheet = MustLoadImage("img/Effects/projectile_blow.png")
var LightMissileBlowSpriteSheet = MustLoadImage("img/Effects/light_missile_blow2.png")
var ClusterMinesBlowSpriteSheet = MustLoadImage("img/Effects/cluster_mines_blow.png")
var PlasmaGunProjectileSpriteSheet = MustLoadImage("img/Effects/plasma_gun1.png")

// fonts
var ScoreFont1024x768 = mustLoadFont("fonts/Kenney Pixel.ttf", 48)
var InfoFont1024x768 = mustLoadFont("fonts/Kenney Future.ttf", 18)
var SmallFont1024x768 = mustLoadFont("fonts/Kenney Mini.ttf", 18)
var SmallUIFont1024x768 = mustLoadFont("fonts/Kenney Future.ttf", 8)
var ProfileFont1024x768 = mustLoadFont("fonts/Kenney Future.ttf", 12)
var ProfileBigFont1024x768 = mustLoadFont("fonts/Kenney Future.ttf", 24)

var ScoreFont1920x1080 = mustLoadFont("fonts/Kenney Pixel.ttf", 96)
var InfoFont1920x1080 = mustLoadFont("fonts/Kenney Pixel.ttf", 64)
var SmallFont1920x1080 = mustLoadFont("fonts/Kenney Mini.ttf", 36)
var ProfileFont1920x1080 = mustLoadFont("fonts/Kenney Future.ttf", 20)
var ProfileBigFont1920x1080 = mustLoadFont("fonts/Kenney Future.ttf", 42)

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
