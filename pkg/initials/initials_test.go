package initials

import (
	"testing"

	"github.com/mxcd/avatar-generator/pkg/font"
	"github.com/stretchr/testify/assert"
	"github.com/zeebo/blake3"
)

func testInitialsGenerationStability(t *testing.T, caching bool) {
	generator, err := NewInitialsAvatarGenerator(&InitialsAvatarOptions{
		AvaratWidth:  256,
		AvatarHeight: 256,
		FontSize:     160,
		FontPath:     "fonts/theboldfont.ttf",
		FontFS:       &font.Fonts,
		Caching:      &caching,
	})

	assert.Nil(t, err)
	assert.NotNil(t, generator)

	initialHashes := map[string]string{}

	initialsList := GetAllInitials()

	count := 0
	byteCount := 0

	for i := 0; i < 2; i++ {
		for _, initials := range initialsList {
			count++
			if count%100 == 0 {
				t.Logf("Processed %d initials", count)
			}

			data, err := generator.GenerateInitialsAvatar(initials)
			assert.Nil(t, err)
			assert.NotNil(t, data)
			assert.Greater(t, len(data), 0)

			byteCount += len(data)

			hasher := blake3.New()
			hasher.Write(data)
			hash := string(hasher.Sum(nil))

			if _, ok := initialHashes[initials]; ok {
				assert.Equal(t, initialHashes[initials], hash)
			} else {
				initialHashes[initials] = hash
			}
		}
	}

	t.Logf("Processed %d initials", count)
	t.Logf("Total size: %d bytes", byteCount)
}

func TestInitialsGenerationStabilityWithoutCaching(t *testing.T) {
	testInitialsGenerationStability(t, false)
}

func TestInitialsGenerationStabilityWithCaching(t *testing.T) {
	testInitialsGenerationStability(t, true)
}

func TestInitialsGenerationWithPreload(t *testing.T) {
	caching := true
	generator, err := NewInitialsAvatarGenerator(&InitialsAvatarOptions{
		AvaratWidth:  256,
		AvatarHeight: 256,
		FontSize:     160,
		FontPath:     "fonts/theboldfont.ttf",
		FontFS:       &font.Fonts,
		Caching:      &caching,
	})

	assert.Nil(t, err)
	assert.NotNil(t, generator)

	initialsList := GetAllInitials()

	err = generator.Preload()
	assert.Nil(t, err)

	assert.Equal(t, len(initialsList), generator.GetCacheSize())

	for _, initials := range initialsList {
		data, err := generator.GenerateInitialsAvatar(initials)
		assert.Nil(t, err)
		assert.NotNil(t, data)
		assert.Greater(t, len(data), 0)
	}
}

func BenchmarkBulkGeneration(b *testing.B) {
	generator, err := NewInitialsAvatarGenerator(&InitialsAvatarOptions{
		AvaratWidth:  256,
		AvatarHeight: 256,
		FontSize:     160,
		FontPath:     "fonts/theboldfont.ttf",
		FontFS:       &font.Fonts,
	})

	assert.Nil(b, err)
	assert.NotNil(b, generator)

	for i := 0; i < b.N; i++ {
		initials := GetRandomInitials()
		_, err := generator.GenerateInitialsAvatar(initials)
		assert.Nil(b, err)
		assert.Greater(b, len(initials), 0)
	}
}
