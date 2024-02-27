package initials

import (
	"bytes"
	"embed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type InitialsAvatarGenerator struct {
	Config   *InitialsAvatarOptions
	FontFace *font.Face
	FontData *truetype.Font
	Cache    *expirable.LRU[string, []byte]
}

type InitialsAvatarOptions struct {
	AvaratWidth  int
	AvatarHeight int
	FontSize     int
	FontPath     string
	FontFS       *embed.FS
	Caching      *bool
}

func NewInitialsAvatarGenerator(config *InitialsAvatarOptions) (*InitialsAvatarGenerator, error) {
	var fontBytes []byte
	if config.FontFS != nil {
		b, err := config.FontFS.ReadFile(config.FontPath)
		if err != nil {
			return nil, err
		}
		fontBytes = b
	} else {
		b, err := os.ReadFile(config.FontPath)
		if err != nil {
			return nil, err
		}
		fontBytes = b
	}

	fontData, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	opts := truetype.Options{
		Size: float64(config.FontSize),
	}
	fontFace := truetype.NewFace(fontData, &opts)

	var cache *expirable.LRU[string, []byte] = nil
	if config.Caching != nil && *config.Caching {
		cache = expirable.NewLRU[string, []byte](0, nil, 0)
	}

	return &InitialsAvatarGenerator{
		Config:   config,
		FontFace: &fontFace,
		FontData: fontData,
		Cache:    cache,
	}, nil
}

func (i *InitialsAvatarGenerator) GenerateInitialsAvatar(initials string) ([]byte, error) {

	if i.Cache != nil {
		if data, ok := i.Cache.Get(initials); ok {
			return data, nil
		}
	}

	foregroundColor := image.White
	rect := image.Rect(0, 0, i.Config.AvaratWidth, i.Config.AvatarHeight)
	backgroundImage := image.NewRGBA(rect)
	backgroundColor := generateAvatarColor(initials)
	draw.Draw(backgroundImage, backgroundImage.Bounds(), &image.Uniform{backgroundColor}, image.Point{}, draw.Src)

	rgba := image.NewRGBA(image.Rect(0, 0, i.Config.AvaratWidth, i.Config.AvatarHeight))
	draw.Draw(rgba, rgba.Bounds(), backgroundImage, image.Point{}, draw.Src)
	freetypeContext := freetype.NewContext()
	freetypeContext.SetFont(i.FontData)
	freetypeContext.SetClip(rgba.Bounds())
	freetypeContext.SetDst(rgba)
	freetypeContext.SetSrc(foregroundColor)
	freetypeContext.SetFontSize(float64(i.Config.FontSize))

	textWidth := font.MeasureString(*i.FontFace, initials).Ceil()
	textHeight := (*i.FontFace).Metrics().Ascent.Ceil() + (*i.FontFace).Metrics().Descent.Ceil()
	x := (i.Config.AvaratWidth - textWidth) / 2
	y := (i.Config.AvatarHeight / 2) + (textHeight / 4)
	pt := freetype.Pt(x, y)
	_, err := freetypeContext.DrawString(initials, pt)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(nil)
	err = png.Encode(buff, rgba)
	if err != nil {
		return nil, err
	}
	bgBlob := buff.Bytes()

	if i.Cache != nil {
		i.Cache.Add(initials, bgBlob)
	}

	return bgBlob, nil
}

func (i *InitialsAvatarGenerator) Preload() error {
	initialsList := GetAllInitials()
	for _, initials := range initialsList {
		_, err := i.GenerateInitialsAvatar(initials)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *InitialsAvatarGenerator) GetCacheSize() int {
	if i.Cache == nil {
		return 0
	}
	return i.Cache.Len()
}

func generateAvatarColor(initials string) color.Color {
	sum := 0
	for _, char := range initials {
		sum += int(char)
	}

	r := 20 + ((sum * 123) % 160)
	g := 20 + ((sum * 456) % 160)
	b := 20 + ((sum * 789) % 160)

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

func GetAllInitials() []string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var initials []string
	for i := range letterRunes {
		for j := range letterRunes {
			initials = append(initials, string(letterRunes[i])+string(letterRunes[j]))
		}
	}
	return initials
}

func GetRandomInitials() string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 2)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
