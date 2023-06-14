package theme

import (
	"bytes"
	"errors"
	"image/color"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/dweymouth/supersonic/backend"
	"github.com/dweymouth/supersonic/res"
	"github.com/dweymouth/supersonic/sharedutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const ColorNamePageBackground fyne.ThemeColorName = "PageBackground"

type AppearanceMode string

const (
	AppearanceLight AppearanceMode = "Light"
	AppearanceDark  AppearanceMode = "Dark"
	AppearanceAuto  AppearanceMode = "Auto"

	DefaultAppearance AppearanceMode = AppearanceDark
)

var (
	normalFont fyne.Resource
	boldFont   fyne.Resource
)

type MyTheme struct {
	NormalFont   string
	BoldFont     string
	config       *backend.ThemeConfig
	themeFileDir string

	loadedThemeFilename string
	loadedThemeFile     *ThemeFile
	defaultThemeFile    *ThemeFile
}

var _ fyne.Theme = (*MyTheme)(nil)

func NewMyTheme(config *backend.ThemeConfig) *MyTheme {
	m := &MyTheme{config: config}
	m.defaultThemeFile, _ = DecodeThemeFile(bytes.NewReader(res.ResDefaultToml.StaticContent))
	m.createThemeIcons()
	return m
}

func (m *MyTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	// load theme file if necessary
	if m.loadedThemeFile == nil || m.config.ThemeFile != m.loadedThemeFilename {
		t, err := ReadThemeFile(path.Join(m.themeFileDir, m.config.ThemeFile))
		if err == nil {
			m.loadedThemeFile = t
		} else {
			log.Printf("failed to load theme file %q: %s", m.config.ThemeFile, err.Error())
			m.loadedThemeFile = m.defaultThemeFile
		}
		m.loadedThemeFilename = m.config.ThemeFile
	}

	variant := m.getVariant()
	thFile := m.loadedThemeFile
	if !thFile.SupportsVariant(variant) {
		thFile = m.defaultThemeFile
	}
	colors := thFile.DarkColors
	if variant == theme.VariantLight {
		colors = thFile.LightColors
	}
	switch name {
	case ColorNamePageBackground:
		return colorOrDefault(colors.PageBackground, name, variant)
	case theme.ColorNameBackground:
		return colorOrDefault(colors.Background, name, variant)
	case theme.ColorNameButton:
		return colorOrDefault(colors.Button, name, variant)
	case theme.ColorNameDisabled:
		return colorOrDefault(colors.Disabled, name, variant)
	case theme.ColorNameDisabledButton:
		return colorOrDefault(colors.DisabledButton, name, variant)
	case theme.ColorNameError:
		return colorOrDefault(colors.Error, name, variant)
	case theme.ColorNameFocus:
		return colorOrDefault(colors.Focus, name, variant)
	case theme.ColorNameForeground:
		return colorOrDefault(colors.Foreground, name, variant)
	case theme.ColorNameHover:
		return colorOrDefault(colors.Hover, name, variant)
	case theme.ColorNameInputBackground:
		return colorOrDefault(colors.InputBackground, name, variant)
	case theme.ColorNameInputBorder:
		return colorOrDefault(colors.InputBorder, name, variant)
	case theme.ColorNameMenuBackground:
		return colorOrDefault(colors.MenuBackground, name, variant)
	case theme.ColorNameOverlayBackground:
		return colorOrDefault(colors.OverlayBackground, name, variant)
	case theme.ColorNamePlaceHolder:
		return colorOrDefault(colors.Placeholder, name, variant)
	case theme.ColorNamePressed:
		return colorOrDefault(colors.Pressed, name, variant)
	case theme.ColorNamePrimary:
		return colorOrDefault(colors.Primary, name, variant)
	case theme.ColorNameScrollBar:
		return colorOrDefault(colors.ScrollBar, name, variant)
	case theme.ColorNameSelection:
		return colorOrDefault(colors.Selection, name, variant)
	case theme.ColorNameSeparator:
		return colorOrDefault(colors.Separator, name, variant)
	case theme.ColorNameShadow:
		return colorOrDefault(colors.Shadow, name, variant)
	case theme.ColorNameSuccess:
		return colorOrDefault(colors.Success, name, variant)
	case theme.ColorNameWarning:
		return colorOrDefault(colors.Warning, name, variant)
	default:
		return colorOrDefault("", name, variant)
	}
}

func colorOrDefault(colorStr string, name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if c, err := ColorStringToColor(colorStr); err == nil {
		return c
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m *MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Returns a map [themeFileName] -> displayName
func (m *MyTheme) ListThemeFiles() map[string]string {
	files, _ := filepath.Glob(m.themeFileDir + "/*.toml")
	result := make(map[string]string)
	for _, filename := range files {
		if themeFile, err := ReadThemeFile(path.Join(m.themeFileDir, filename)); err == nil {
			result[filename] = themeFile.SupersonicTheme.Name
		}
	}
	return result
}

type myThemedResource struct {
	myTheme      MyTheme
	darkVariant  *fyne.StaticResource
	lightVariant *fyne.StaticResource
}

var _ fyne.Resource = myThemedResource{}

func (p myThemedResource) Content() []byte {
	if p.myTheme.getVariant() == theme.VariantDark {
		return p.darkVariant.StaticContent
	}
	return p.lightVariant.StaticContent
}

func (p myThemedResource) Name() string {
	if p.myTheme.getVariant() == theme.VariantDark {
		return p.darkVariant.StaticName
	}
	return p.lightVariant.StaticName
}

var (
	AlbumIcon       fyne.Resource
	ArtistIcon      fyne.Resource
	FavoriteIcon    fyne.Resource
	NotFavoriteIcon fyne.Resource
	GenreIcon       fyne.Resource
	NowPlayingIcon  fyne.Resource
	PlaylistIcon    fyne.Resource
	ShuffleIcon     fyne.Resource
	TracksIcon      fyne.Resource
	FilterIcon      fyne.Resource = theme.NewThemedResource(res.ResFilterSvg)
)

// MUST be called at startup!
func (m MyTheme) createThemeIcons() {
	AlbumIcon = myThemedResource{myTheme: m, darkVariant: res.ResDiscInvertPng, lightVariant: res.ResDiscPng}
	ArtistIcon = myThemedResource{myTheme: m, darkVariant: res.ResPeopleInvertPng, lightVariant: res.ResPeoplePng}
	FavoriteIcon = myThemedResource{myTheme: m, darkVariant: res.ResHeartFilledInvertPng, lightVariant: res.ResHeartFilledPng}
	NotFavoriteIcon = myThemedResource{myTheme: m, darkVariant: res.ResHeartOutlineInvertPng, lightVariant: res.ResHeartOutlinePng}
	GenreIcon = myThemedResource{myTheme: m, darkVariant: res.ResTheatermasksInvertPng, lightVariant: res.ResTheatermasksPng}
	NowPlayingIcon = myThemedResource{myTheme: m, darkVariant: res.ResHeadphonesInvertPng, lightVariant: res.ResHeadphonesPng}
	PlaylistIcon = myThemedResource{myTheme: m, darkVariant: res.ResPlaylistInvertPng, lightVariant: res.ResPlaylistPng}
	ShuffleIcon = myThemedResource{myTheme: m, darkVariant: res.ResShuffleInvertSvg, lightVariant: res.ResShuffleSvg}
	TracksIcon = myThemedResource{myTheme: m, darkVariant: res.ResMusicnotesInvertPng, lightVariant: res.ResMusicnotesPng}
}

func (m *MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	switch style {
	case fyne.TextStyle{}:
		if m.NormalFont != "" && normalFont == nil {
			if content, err := readTTFFile(m.NormalFont); err != nil {
				m.NormalFont = ""
				m.BoldFont = ""
			} else {
				normalFont = fyne.NewStaticResource("normalFont", content)
			}
		}
		if normalFont != nil {
			return normalFont
		}
	case fyne.TextStyle{Bold: true}:
		if m.BoldFont != "" && boldFont == nil {
			if content, err := ioutil.ReadFile(m.BoldFont); err != nil {
				m.BoldFont = ""
			} else {
				normalFont = fyne.NewStaticResource("boldFont", content)
			}
		}
		if boldFont != nil {
			return boldFont
		}
		if normalFont != nil {
			return normalFont
		}
	}
	return theme.DefaultTheme().Font(style)
}

func (m *MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m *MyTheme) getVariant() fyne.ThemeVariant {
	v := DefaultAppearance // default if config has invalid or missing setting
	if sharedutil.SliceContains(
		[]string{string(AppearanceLight), string(AppearanceDark), string(AppearanceAuto)},
		m.config.Appearance) {
		v = AppearanceMode(m.config.Appearance)
	}

	if AppearanceMode(v) == AppearanceDark {
		return theme.VariantDark
	} else if AppearanceMode(v) == AppearanceLight {
		return theme.VariantLight
	}
	return fyne.CurrentApp().Settings().ThemeVariant()
}

func readTTFFile(filepath string) ([]byte, error) {
	if !strings.HasSuffix(filepath, ".ttf") {
		err := errors.New("only .ttf fonts are supported")
		log.Printf("error loading custom font %q: %s", filepath, err.Error())
		return nil, err
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("error loading custom font %q: %s", filepath, err.Error())
	}
	return content, err
}
