package cartridge

import (
	"embed"
	"fmt"
	"image"
	"math"
	"math/rand"
	"time"

	"github.com/TheMightyGit/marv/marvlib"
	"github.com/TheMightyGit/marv/marvtypes"
)

//go:embed "resources/*"
var Resources embed.FS

const (
	SpriteBGDots1 = iota
	SpriteBGDots2
	SpriteText
	SpriteBaddieStart
	SpriteBaddieEnd = iota + 120
	SpritePlayer
)

const (
	MapBankGfx = iota
	MapBankText
)

const (
	GfxBankFont = iota
)

type Element struct {
	SpriteId int
	Pos      image.Point
}

type Baddie struct {
	Element
}

func (b *Baddie) Setup() {
	marvlib.API.SpritesGet(b.SpriteId).Show(GfxBankFont, area)
	marvlib.API.SpritesGet(b.SpriteId).ChangePos(image.Rectangle{
		Min: b.Pos,
		Max: image.Point{X: 6, Y: 8},
	})
	marvlib.API.SpritesGet(b.SpriteId).ChangeViewport(image.Point{X: 3 * 6, Y: 15 * 8})
	go b.MoveBrain() // yes, background 'threads' (actually coroutines)
}

func (b *Baddie) IsHit(rect image.Rectangle) bool {
	baddieRect := image.Rectangle{
		Min: b.Pos.Sub(image.Point{X: 2, Y: 2}),
		Max: b.Pos.Add(image.Point{X: 2, Y: 2}),
	}
	return rect.Overlaps(baddieRect)
}

func (b *Baddie) MoveBrain() {
	frame := 0
	timeout := time.NewTicker(16 * time.Millisecond)
	for {
		select {
		case <-marvlib.API.ConsoleResetChan():
			fmt.Println(b.SpriteId, "CANCELLED!!!!!!!!!!!!!!!!!!")
			return
		case <-timeout.C:
			if frame < 40 {
				// wobble on the spot
				b.Pos.X += rand.Intn(3) - 1
				b.Pos.Y += rand.Intn(3) - 1
			} else if frame < 60 {
				// run at player
				d := 1
				if player.Dead { // if player dead then lose interest and run away
					d = -2
				}
				// but also generally move towards the player
				if player.Pos.X < b.Pos.X {
					b.Pos.X -= d
				} else if player.Pos.X > b.Pos.X {
					b.Pos.X += d
				}
				if player.Pos.Y < b.Pos.Y {
					b.Pos.Y -= d
				} else if player.Pos.Y > b.Pos.Y {
					b.Pos.Y += d
				}
			} else {
				frame = 0
			}
			marvlib.API.SpritesGet(b.SpriteId).ChangePos(image.Rectangle{
				Min: b.Pos.Sub(image.Point{X: 3, Y: 4}),
				Max: image.Point{X: 6, Y: 8},
			})
			frame++
			timeout.Reset(16 * time.Millisecond)
		}
	}
}

type Player struct {
	Dead bool
	Element
}

func (p *Player) Setup() {
	marvlib.API.SpritesGet(p.SpriteId).Show(GfxBankFont, area)
	marvlib.API.SpritesGet(p.SpriteId).ChangePos(image.Rectangle{
		Min: p.Pos.Sub(image.Point{X: 3, Y: 4}),
		Max: image.Point{X: 6, Y: 8},
	})
	marvlib.API.SpritesGet(p.SpriteId).ChangeViewport(image.Point{X: 10 * 6, Y: 14 * 8})
}

const (
	MAPPED_GAMEPAD_BIT_DPAD_UP    = uint16(1 << 0)
	MAPPED_GAMEPAD_BIT_DPAD_DOWN  = uint16(1 << 1)
	MAPPED_GAMEPAD_BIT_DPAD_LEFT  = uint16(1 << 2)
	MAPPED_GAMEPAD_BIT_DPAD_RIGHT = uint16(1 << 3)

	MAPPED_GAMEPAD_BIT_BUTTON_BOTTOM = uint16(1 << 4)
	MAPPED_GAMEPAD_BIT_BUTTON_LEFT   = uint16(1 << 5)
	MAPPED_GAMEPAD_BIT_BUTTON_TOP    = uint16(1 << 6)
	MAPPED_GAMEPAD_BIT_BUTTON_RIGHT  = uint16(1 << 7)

	MAPPED_GAMEPAD_BIT_SHOULDER_LEFT  = uint16(1 << 8)
	MAPPED_GAMEPAD_BIT_SHOULDER_RIGHT = uint16(1 << 9)

	MAPPED_GAMEPAD_BIT_SELECT = uint16(1 << 10)
	MAPPED_GAMEPAD_BIT_START  = uint16(1 << 11)
)

func isMoveLeft() bool {
	return marvlib.API.InputIsKeyDown(marvlib.KeyA) ||
		marvlib.API.InputIsKeyDown(marvlib.KeyArrowLeft) // ||
	// (marv.Input.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_LEFT != 0) ||
	// (marv.Input.GamepadButtonStates[1].Mapped&MAPPED_GAMEPAD_BIT_DPAD_LEFT != 0) ||
	// (marv.Input.GamepadButtonStates[2].Mapped&MAPPED_GAMEPAD_BIT_DPAD_LEFT != 0) ||
	// (marv.Input.GamepadButtonStates[3].Mapped&MAPPED_GAMEPAD_BIT_DPAD_LEFT != 0)
}
func isMoveRight() bool {
	return marvlib.API.InputIsKeyDown(marvlib.KeyD) ||
		marvlib.API.InputIsKeyDown(marvlib.KeyArrowRight) // ||
	// (marv.Input.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_RIGHT != 0) ||
	// (marv.Input.GamepadButtonStates[1].Mapped&MAPPED_GAMEPAD_BIT_DPAD_RIGHT != 0) ||
	// (marv.Input.GamepadButtonStates[2].Mapped&MAPPED_GAMEPAD_BIT_DPAD_RIGHT != 0) ||
	// (marv.Input.GamepadButtonStates[3].Mapped&MAPPED_GAMEPAD_BIT_DPAD_RIGHT != 0)
}
func isMoveUp() bool {
	return marvlib.API.InputIsKeyDown(marvlib.KeyW) ||
		marvlib.API.InputIsKeyDown(marvlib.KeyArrowUp) // ||
	// (marv.Input.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_UP != 0) ||
	// (marv.Input.GamepadButtonStates[1].Mapped&MAPPED_GAMEPAD_BIT_DPAD_UP != 0) ||
	// (marv.Input.GamepadButtonStates[2].Mapped&MAPPED_GAMEPAD_BIT_DPAD_UP != 0) ||
	// (marv.Input.GamepadButtonStates[3].Mapped&MAPPED_GAMEPAD_BIT_DPAD_UP != 0)
}
func isMoveDown() bool {
	return marvlib.API.InputIsKeyDown(marvlib.KeyS) ||
		marvlib.API.InputIsKeyDown(marvlib.KeyArrowDown) // ||
	// (marv.Input.GamepadButtonStates[0].Mapped&MAPPED_GAMEPAD_BIT_DPAD_DOWN != 0) ||
	// (marv.Input.GamepadButtonStates[1].Mapped&MAPPED_GAMEPAD_BIT_DPAD_DOWN != 0) ||
	// (marv.Input.GamepadButtonStates[2].Mapped&MAPPED_GAMEPAD_BIT_DPAD_DOWN != 0) ||
	// (marv.Input.GamepadButtonStates[3].Mapped&MAPPED_GAMEPAD_BIT_DPAD_DOWN != 0)
}

func (p *Player) Update() {
	if !p.Dead {
		playerRect := image.Rectangle{
			Min: p.Pos.Sub(image.Point{X: 3, Y: 4}),
			Max: p.Pos.Add(image.Point{X: 6, Y: 8}),
		}
		for i := range baddies {
			if baddies[i].IsHit(playerRect) {
				player.Dead = true
				notify("Oh damn!")
				time.AfterFunc(5*time.Second, func() {
					marvlib.API.ConsoleReset()
				})
			}
		}
	}

	if isMoveLeft() {
		p.Pos.X--
	}
	if isMoveRight() {
		p.Pos.X++
	}
	if isMoveUp() {
		p.Pos.Y--
	}
	if isMoveDown() {
		p.Pos.Y++
	}

	/*
		// move towards mouse
		if p.Pos.X < marv.Input.MousePos.X {
			p.Pos.X++
		} else if player.Pos.X > marv.Input.MousePos.X {
			p.Pos.X--
		}
		if player.Pos.Y < marv.Input.MousePos.Y {
			p.Pos.Y++
		} else if player.Pos.Y > marv.Input.MousePos.Y {
			p.Pos.Y--
		}
	*/

	if p.Dead {
		// we're a ghost now!! wooooo!
		marvlib.API.SpritesGet(p.SpriteId).ChangePos(image.Rectangle{
			Min: p.Pos.Sub(image.Point{X: 3, Y: 4}),
			Max: image.Point{X: 6, Y: 8},
		})
		marvlib.API.SpritesGet(p.SpriteId).ChangeViewport(image.Point{X: 3 * 6, Y: 15 * 8}) // ghost graphic
	} else {
		marvlib.API.SpritesGet(p.SpriteId).ChangePos(image.Rectangle{
			Min: p.Pos.Sub(image.Point{X: 3, Y: 4}),
			Max: image.Point{X: 6, Y: 8},
		}) // player graphic
	}
}

var (
	area       marvtypes.MapBankArea
	areaDots   marvtypes.MapBankArea
	areaNotify marvtypes.MapBankArea
	baddies    []*Baddie
	player     *Player
)

func addBaddiesOverTime() {
	timeout := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-marvlib.API.ConsoleResetChan():
			fmt.Println("CANCELLED!!!!!!!!!!!!!!!!!!")
			return
		case <-timeout.C:
			if len(baddies) < SpriteBaddieEnd-SpriteBaddieStart { // any baddies left to add?
				addBaddie()
			}
			timeout.Reset(1 * time.Second)
		}
	}
}

var space = "                                                         "

func notify(text string) {
	areaNotify.StringToMap(image.Point{}, 1, 16, space[:(40-len(text))/2]+text+space)
}

func setupNotifyText() {
	areaNotify = marvlib.API.MapBanksGet(MapBankText).AllocArea(image.Point{X: 40, Y: 1})
	notify("Better run for it!")

	marvlib.API.SpritesGet(SpriteText).Show(GfxBankFont, areaNotify)
	marvlib.API.SpritesGet(SpriteText).ChangePos(image.Rectangle{
		Min: image.Point{X: 0, Y: 120},
		Max: image.Point{X: 240, Y: 8},
	})
}

func Start() {
	/*
		go func() {
			timeout := time.NewTicker(1 * time.Second)
			for {
				select {
				case <-marv.ctx.Done():
					fmt.Println("CANCELLED!!!!!!!!!!!!!!!!!!")
					return
				case <-timeout.C:
					fmt.Println("I AM IN THE BACKGROUND!!!!!")
					timeout.Reset(1 * time.Second)
				}
			}
		}()
	*/

	// marv.ModBanks[0].Play()
	marvlib.API.SfxBanksGet(0).PlayLooped()

	setupArea()
	setupBackground()
	setupBaddies()
	setupPlayer()
	setupNotifyText()
	go addBaddiesOverTime()
}

func Update() {
	updateBackground()
	player.Update()
}

func setupArea() {
	// mimic the sprite sheet layout directly into the area
	area = marvlib.API.MapBanksGet(MapBankGfx).AllocArea(image.Point{X: 32, Y: 32})
	p := image.Point{}
	for p.Y = 0; p.Y < 32; p.Y++ {
		for p.X = 0; p.X < 32; p.X++ {
			if p.Y == 15 && p.X == 3 { // ghosty
				area.Set(p, uint8(p.X), uint8(p.Y), 10, 16)
			} else {
				area.Set(p, uint8(p.X), uint8(p.Y), 14, 16)
			}
		}
	}
}

var bgCount float64

func updateBackground() {
	bgCount += 0.025
	marvlib.API.SpritesGet(SpriteBGDots1).ChangeViewport(image.Point{
		X: 64 + int(math.Cos(bgCount)*30),
		Y: 64 + int(math.Sin(bgCount)*30),
	})
	marvlib.API.SpritesGet(SpriteBGDots2).ChangeViewport(image.Point{
		X: 64 + int(math.Cos(bgCount)*15),
		Y: 64 + int(math.Sin(bgCount)*15),
	})
}

func setupBackground() {
	areaDots = marvlib.API.MapBanksGet(MapBankGfx).AllocArea(image.Point{X: 60, Y: 50}) // at 8x8
	p := image.Point{}
	for p.Y = 0; p.Y < 50; p.Y++ {
		for p.X = 0; p.X < 60; p.X++ {
			areaDots.Set(p, 14, 2, 0, 16)
		}
	}

	marvlib.API.SpritesGet(SpriteBGDots1).Show(GfxBankFont, areaDots)
	marvlib.API.SpritesGet(SpriteBGDots2).Show(GfxBankFont, areaDots)

	marvlib.API.SpritesGet(SpriteBGDots1).ChangePos(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: 240, Y: 240},
	})
	marvlib.API.SpritesGet(SpriteBGDots2).ChangePos(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: 240, Y: 240},
	})
}

func setupBaddies() {
	baddies = make([]*Baddie, 0)
	for i := 0; i < 10; i++ {
		addBaddie()
	}
}

func baddieSpawnPoint() image.Point {
	pos := rand.Intn(240)
	p := image.Point{}

	switch rand.Intn(4) {
	case 0: // north
		p.X = pos
		p.Y = -20
	case 1: // east
		p.X = 240 + 20
		p.Y = pos
	case 2: // south
		p.X = pos
		p.Y = 240 + 20
	case 3: // west
		p.X = -20
		p.Y = pos
	default:
		panic("HUH?")
	}
	return p
}

func addBaddie() {
	b := &Baddie{
		Element: Element{
			SpriteId: SpriteBaddieStart + len(baddies),
			Pos:      baddieSpawnPoint(),
		},
	}
	b.Setup()
	baddies = append(baddies, b)
}

func setupPlayer() {
	player = &Player{
		Element: Element{
			SpriteId: SpritePlayer,
			Pos: image.Point{
				X: 120,
				Y: 120,
			},
		},
	}
	player.Setup()
}
