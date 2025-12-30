package captcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/aiflowy/aiflowy-go/pkg/snowflake"
)

// Store holds captcha data with expiration
type Store struct {
	mu      sync.RWMutex
	data    map[string]captchaData
	timeout time.Duration
}

type captchaData struct {
	code      string
	expiresAt time.Time
}

var defaultStore *Store

func init() {
	defaultStore = NewStore(5 * time.Minute)
	// Start cleanup goroutine
	go defaultStore.cleanup()
}

// NewStore creates a new captcha store
func NewStore(timeout time.Duration) *Store {
	return &Store{
		data:    make(map[string]captchaData),
		timeout: timeout,
	}
}

// cleanup removes expired captchas periodically
func (s *Store) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, data := range s.data {
			if now.After(data.expiresAt) {
				delete(s.data, id)
			}
		}
		s.mu.Unlock()
	}
}

// Set stores a captcha
func (s *Store) Set(id, code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = captchaData{
		code:      code,
		expiresAt: time.Now().Add(s.timeout),
	}
}

// Verify verifies a captcha and removes it
func (s *Store) Verify(id, code string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.data[id]
	if !ok {
		return false
	}

	// Delete the captcha regardless of result
	delete(s.data, id)

	// Check if expired
	if time.Now().After(data.expiresAt) {
		return false
	}

	// Case-insensitive comparison
	return equalFold(data.code, code)
}

func equalFold(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		c1, c2 := s1[i], s2[i]
		if c1 >= 'A' && c1 <= 'Z' {
			c1 += 'a' - 'A'
		}
		if c2 >= 'A' && c2 <= 'Z' {
			c2 += 'a' - 'A'
		}
		if c1 != c2 {
			return false
		}
	}
	return true
}

// CaptchaResult contains the captcha data
type CaptchaResult struct {
	ID   string `json:"captchaId"`
	Data string `json:"captchaData"` // Base64 encoded image
}

// Generate generates a new captcha
func Generate() (*CaptchaResult, error) {
	// Generate captcha ID
	id, err := snowflake.GenerateIDString()
	if err != nil {
		id = generateRandomID()
	}

	// Generate random code (4 characters)
	code := generateCode(4)

	// Generate image
	imgData, err := generateImage(code)
	if err != nil {
		return nil, err
	}

	// Store captcha
	defaultStore.Set(id, code)

	return &CaptchaResult{
		ID:   id,
		Data: "data:image/png;base64," + imgData,
	}, nil
}

// Verify verifies a captcha
func Verify(id, code string) bool {
	if id == "" || code == "" {
		return false
	}
	return defaultStore.Verify(id, code)
}

// generateCode generates a random alphanumeric code
func generateCode(length int) string {
	const charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ" // Exclude confusing characters
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// generateRandomID generates a random ID
func generateRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 32)
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}
	return string(id)
}

// generateImage generates a captcha image
func generateImage(code string) (string, error) {
	width := 120
	height := 40

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background with random light color
	bgColor := color.RGBA{
		R: uint8(230 + rand.Intn(25)),
		G: uint8(230 + rand.Intn(25)),
		B: uint8(230 + rand.Intn(25)),
		A: 255,
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw noise dots
	for i := 0; i < 100; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		img.Set(x, y, color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		})
	}

	// Draw noise lines
	for i := 0; i < 3; i++ {
		lineColor := color.RGBA{
			R: uint8(rand.Intn(128)),
			G: uint8(rand.Intn(128)),
			B: uint8(rand.Intn(128)),
			A: 255,
		}
		x1, y1 := rand.Intn(width), rand.Intn(height)
		x2, y2 := rand.Intn(width), rand.Intn(height)
		drawLine(img, x1, y1, x2, y2, lineColor)
	}

	// Draw text
	textColor := color.RGBA{
		R: uint8(rand.Intn(100)),
		G: uint8(rand.Intn(100)),
		B: uint8(rand.Intn(100)),
		A: 255,
	}

	// Use basic font
	face := basicfont.Face7x13
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
	}

	// Calculate starting position
	startX := 20
	charWidth := 22

	for i, c := range code {
		// Add some random vertical offset
		yOffset := rand.Intn(8) - 4
		d.Dot = fixed.P(startX+i*charWidth, height/2+5+yOffset)
		d.DrawString(string(c))
	}

	// Encode to PNG and then to base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// drawLine draws a line using Bresenham's algorithm
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x1, y1, c)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
