package main

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) Option {
	return func(p *GamePerson) {
		copy(p.name[:], name)
	}
}

func WithCoordinates(x, y, z int) Option {
	return func(p *GamePerson) {
		p.x = int32(x)
		p.y = int32(y)
		p.z = int32(z)
	}
}

func WithGold(gold int) Option {
	return func(p *GamePerson) {
		p.gold = Gold(gold)
	}
}

func boundsCheck(min, max, value int) int {
	if value < min {
		value = min
	} else if value > max {
		value = max
	}

	return value
}

func WithMana(mana int) Option {
	return func(p *GamePerson) {
		mana = boundsCheck(0, 1000, mana)

		raw := uint32(p.manaHealth[0]) |
			uint32(p.manaHealth[1])<<8 |
			uint32(p.manaHealth[2])<<16

		raw &^= 0x3FF
		raw |= uint32(mana) & 0x3FF

		p.manaHealth[0] = byte(raw)
		p.manaHealth[1] = byte(raw >> 8)
		p.manaHealth[2] = byte(raw >> 16)
	}
}

func WithHealth(health int) Option {
	return func(p *GamePerson) {
		health = boundsCheck(0, 1000, health)

		raw := uint32(p.manaHealth[0]) |
			uint32(p.manaHealth[1])<<8 |
			uint32(p.manaHealth[2])<<16

		raw &^= 0x3FF << 10
		raw |= (uint32(health) & 0x3FF) << 10

		p.manaHealth[0] = byte(raw)
		p.manaHealth[1] = byte(raw >> 8)
		p.manaHealth[2] = byte(raw >> 16)
	}
}

func WithRespect(respect int) Option {
	return func(p *GamePerson) {
		respect = boundsCheck(0, 10, respect)

		p.attrs &^= 0xF000
		p.attrs |= RespectStrengthExperienceLevel(respect) << 12
	}
}

func WithStrength(strength int) Option {
	return func(p *GamePerson) {
		strength = boundsCheck(0, 10, strength)

		p.attrs &^= 0x0F00
		p.attrs |= RespectStrengthExperienceLevel(strength) << 8
	}
}

func WithExperience(experience int) Option {
	return func(p *GamePerson) {
		experience = boundsCheck(0, 10, experience)

		p.attrs &^= 0x00F0
		p.attrs |= RespectStrengthExperienceLevel(experience) << 4
	}
}

func WithLevel(level int) Option {
	return func(p *GamePerson) {
		level = boundsCheck(0, 10, level)

		p.attrs &^= 0x000F
		p.attrs |= RespectStrengthExperienceLevel(level)
	}
}

func WithHouse() Option {
	return func(p *GamePerson) {
		p.params |= 0b0100
	}
}

func WithGun() Option {
	return func(p *GamePerson) {
		p.params |= 0b0010
	}
}

func WithFamily() Option {
	return func(p *GamePerson) {
		p.params |= 0b0001
	}
}

func WithType(personType int) Option {
	return func(p *GamePerson) {
		typeVal := (PersonTypeHouseGunFamily(personType) & 0x3) << 4
		p.params &^= 0b00110000
		p.params |= typeVal
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type Coordinates struct {
	x, y, z int32
}

type Gold uint32

type ManaHealth [3]byte

/*
*
15 14 13 12 | 11 10 9 8 | 7 6 5 4 | 3 2 1 0

	Respect  | Strength  | Exp     | Level
*/
type RespectStrengthExperienceLevel uint16

/*
*

	7 6 5 4 | 3 2 1 0

PersonType| Flags
*/
type PersonTypeHouseGunFamily uint8

type Name [42]byte

type GamePerson struct {
	Coordinates
	gold       Gold
	manaHealth ManaHealth
	params     PersonTypeHouseGunFamily
	attrs      RespectStrengthExperienceLevel
	name       Name
}

func NewGamePerson(options ...Option) GamePerson {
	var g = GamePerson{}
	for _, option := range options {
		option(&g)
	}
	return g
}

func (p *GamePerson) Name() string {
	n := 0
	for n < len(p.name) && p.name[n] != 0 {
		n++
	}
	return string(p.name[:n])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	v := uint32(p.manaHealth[0]) | (uint32(p.manaHealth[1]) << 8) | (uint32(p.manaHealth[2]) << 16)
	return int(v & 0x3FF)
}

func (p *GamePerson) Health() int {
	v := uint32(p.manaHealth[0]) | (uint32(p.manaHealth[1]) << 8) | (uint32(p.manaHealth[2]) << 16)
	return int((v >> 10) & 0x3FF)
}

func (p *GamePerson) Respect() int {
	return int(p.attrs & 0xF000 >> 12)
}

func (p *GamePerson) Strength() int {
	return int(p.attrs & 0xF00 >> 8)
}

func (p *GamePerson) Experience() int {
	return int(p.attrs & 0xF0 >> 4)
}

func (p *GamePerson) Level() int {
	return int(p.attrs & 0xF)
}

func (p *GamePerson) HasHouse() bool {
	return (p.params & 0b0100) == 0b0100
}

func (p *GamePerson) HasGun() bool {
	return (p.params & 0b0010) == 0b0010
}

func (p *GamePerson) HasFamilty() bool {
	return (p.params & 0b0001) == 0b0001
}

func (p *GamePerson) Type() int {
	return int(p.params & 0xF0 >> 4)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)

	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}

func TestGamePersonEmpty(t *testing.T) {
	person := NewGamePerson()

	assert.Equal(t, "", person.Name())
	assert.Equal(t, 0, person.X())
	assert.Equal(t, 0, person.Y())
	assert.Equal(t, 0, person.Z())
	assert.Equal(t, 0, person.Gold())
	assert.Equal(t, 0, person.Mana())
	assert.Equal(t, 0, person.Health())
	assert.Equal(t, 0, person.Respect())
	assert.Equal(t, 0, person.Strength())
	assert.Equal(t, 0, person.Experience())
	assert.Equal(t, 0, person.Level())
	assert.False(t, person.HasHouse())
	assert.False(t, person.HasGun())
	assert.False(t, person.HasFamilty())
	assert.Equal(t, 0, person.Type())
}

func TestGamePersonBounds(t *testing.T) {
	person := NewGamePerson(
		WithMana(-100),
		WithHealth(1500),
		WithRespect(-5),
		WithStrength(15),
		WithExperience(-1),
		WithLevel(20),
	)

	assert.Equal(t, 0, person.Mana())
	assert.Equal(t, 1000, person.Health())
	assert.Equal(t, 0, person.Respect())
	assert.Equal(t, 10, person.Strength())
	assert.Equal(t, 0, person.Experience())
	assert.Equal(t, 10, person.Level())
}

func TestGamePersonNameTruncation(t *testing.T) {
	longName := "very_long_name_that_exceeds_the_maximum_length_of_forty_two_characters"
	person := NewGamePerson(WithName(longName))

	assert.Equal(t, longName[:42], person.Name())
}

func TestGamePersonNameWithNullBytes(t *testing.T) {
	nameWithNull := "test\x00name"
	person := NewGamePerson(WithName(nameWithNull))

	assert.Equal(t, "test", person.Name())
}

func TestGamePersonAllTypes(t *testing.T) {
	types := []int{BuilderGamePersonType, BlacksmithGamePersonType, WarriorGamePersonType}

	for _, personType := range types {
		person := NewGamePerson(WithType(personType))
		assert.Equal(t, personType, person.Type())
	}
}

func TestGamePersonAllFlags(t *testing.T) {
	person := NewGamePerson(
		WithHouse(),
		WithGun(),
		WithFamily(),
	)

	assert.True(t, person.HasHouse())
	assert.True(t, person.HasGun())
	assert.True(t, person.HasFamilty())
}

func TestGamePersonPartialFlags(t *testing.T) {
	person1 := NewGamePerson(WithHouse())
	assert.True(t, person1.HasHouse())
	assert.False(t, person1.HasGun())
	assert.False(t, person1.HasFamilty())

	person2 := NewGamePerson(WithGun())
	assert.False(t, person2.HasHouse())
	assert.True(t, person2.HasGun())
	assert.False(t, person2.HasFamilty())

	person3 := NewGamePerson(WithFamily())
	assert.False(t, person3.HasHouse())
	assert.False(t, person3.HasGun())
	assert.True(t, person3.HasFamilty())
}

func TestGamePersonManaHealthCombination(t *testing.T) {
	person := NewGamePerson(
		WithMana(500),
		WithHealth(750),
	)

	assert.Equal(t, 500, person.Mana())
	assert.Equal(t, 750, person.Health())

	person2 := NewGamePerson(
		WithMana(500),
		WithHealth(750),
		WithMana(250),
	)

	assert.Equal(t, 250, person2.Mana())
	assert.Equal(t, 750, person2.Health())

	person3 := NewGamePerson(
		WithMana(500),
		WithHealth(750),
		WithHealth(300),
	)

	assert.Equal(t, 500, person3.Mana())
	assert.Equal(t, 300, person3.Health())
}

func TestGamePersonAttributesCombination(t *testing.T) {
	person := NewGamePerson(
		WithRespect(5),
		WithStrength(7),
		WithExperience(3),
		WithLevel(9),
	)

	assert.Equal(t, 5, person.Respect())
	assert.Equal(t, 7, person.Strength())
	assert.Equal(t, 3, person.Experience())
	assert.Equal(t, 9, person.Level())

	person2 := NewGamePerson(
		WithRespect(5),
		WithStrength(7),
		WithExperience(3),
		WithLevel(9),
		WithRespect(2),
	)

	assert.Equal(t, 2, person2.Respect())
	assert.Equal(t, 7, person2.Strength())
	assert.Equal(t, 3, person2.Experience())
	assert.Equal(t, 9, person2.Level())
}

func TestGamePersonCoordinates(t *testing.T) {
	person := NewGamePerson(
		WithCoordinates(100, -200, 300),
	)

	assert.Equal(t, 100, person.X())
	assert.Equal(t, -200, person.Y())
	assert.Equal(t, 300, person.Z())
}

func TestGamePersonGold(t *testing.T) {
	person := NewGamePerson(WithGold(12345))
	assert.Equal(t, 12345, person.Gold())

	person2 := NewGamePerson(WithGold(0))
	assert.Equal(t, 0, person2.Gold())
}

func TestGamePersonSize(t *testing.T) {
	size := unsafe.Sizeof(GamePerson{})
	assert.LessOrEqual(t, size, uintptr(64))

	t.Logf("GamePerson size: %d bytes", size)
}

func TestGamePersonFieldAlignment(t *testing.T) {
	var person GamePerson

	nameOffset := unsafe.Offsetof(person.name)
	coordsOffset := unsafe.Offsetof(person.Coordinates)
	goldOffset := unsafe.Offsetof(person.gold)
	manaHealthOffset := unsafe.Offsetof(person.manaHealth)
	paramsOffset := unsafe.Offsetof(person.params)
	attrsOffset := unsafe.Offsetof(person.attrs)

	t.Logf("Field offsets: name=%d, coords=%d, gold=%d, manaHealth=%d, params=%d, attrs=%d",
		nameOffset, coordsOffset, goldOffset, manaHealthOffset, paramsOffset, attrsOffset)

	t.Logf("Field sizes: name=%d, coords=%d, gold=%d, manaHealth=%d, params=%d, attrs=%d",
		unsafe.Sizeof(person.name), unsafe.Sizeof(person.Coordinates), unsafe.Sizeof(person.gold),
		unsafe.Sizeof(person.manaHealth), unsafe.Sizeof(person.params), unsafe.Sizeof(person.attrs))

	totalSize := unsafe.Sizeof(person)
	t.Logf("Total structure size: %d bytes", totalSize)
	assert.LessOrEqual(t, totalSize, uintptr(64))

	assert.Greater(t, totalSize, uintptr(0))
}

func TestGamePersonMultipleOptions(t *testing.T) {
	person := NewGamePerson(
		WithName("Alice"),
		WithCoordinates(10, 20, 30),
		WithGold(100),
		WithMana(50),
		WithHealth(75),
		WithRespect(3),
		WithStrength(5),
		WithExperience(7),
		WithLevel(2),
		WithHouse(),
		WithGun(),
		WithType(WarriorGamePersonType),
	)

	assert.Equal(t, "Alice", person.Name())
	assert.Equal(t, 10, person.X())
	assert.Equal(t, 20, person.Y())
	assert.Equal(t, 30, person.Z())
	assert.Equal(t, 100, person.Gold())
	assert.Equal(t, 50, person.Mana())
	assert.Equal(t, 75, person.Health())
	assert.Equal(t, 3, person.Respect())
	assert.Equal(t, 5, person.Strength())
	assert.Equal(t, 7, person.Experience())
	assert.Equal(t, 2, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasGun())
	assert.False(t, person.HasFamilty())
	assert.Equal(t, WarriorGamePersonType, person.Type())
}

func TestGamePersonOptionReapplication(t *testing.T) {
	person := NewGamePerson(
		WithName("Bob"),
		WithName("Charlie"),
		WithGold(100),
		WithGold(200),
		WithMana(50),
		WithMana(25),
	)

	assert.Equal(t, "Charlie", person.Name())
	assert.Equal(t, 200, person.Gold())
	assert.Equal(t, 25, person.Mana())
}

func TestGamePersonZeroValues(t *testing.T) {
	person := NewGamePerson(
		WithName(""),
		WithCoordinates(0, 0, 0),
		WithGold(0),
		WithMana(0),
		WithHealth(0),
		WithRespect(0),
		WithStrength(0),
		WithExperience(0),
		WithLevel(0),
		WithType(0),
	)

	assert.Equal(t, "", person.Name())
	assert.Equal(t, 0, person.X())
	assert.Equal(t, 0, person.Y())
	assert.Equal(t, 0, person.Z())
	assert.Equal(t, 0, person.Gold())
	assert.Equal(t, 0, person.Mana())
	assert.Equal(t, 0, person.Health())
	assert.Equal(t, 0, person.Respect())
	assert.Equal(t, 0, person.Strength())
	assert.Equal(t, 0, person.Experience())
	assert.Equal(t, 0, person.Level())
	assert.Equal(t, 0, person.Type())
}

func TestGamePersonMaxValues(t *testing.T) {
	person := NewGamePerson(
		WithName("MaxValueTest"),
		WithCoordinates(math.MaxInt32, math.MaxInt32, math.MaxInt32),
		WithGold(math.MaxInt32),
		WithMana(1000),
		WithHealth(1000),
		WithRespect(10),
		WithStrength(10),
		WithExperience(10),
		WithLevel(10),
		WithHouse(),
		WithGun(),
		WithFamily(),
		WithType(WarriorGamePersonType),
	)

	assert.Equal(t, "MaxValueTest", person.Name())
	assert.Equal(t, math.MaxInt32, person.X())
	assert.Equal(t, math.MaxInt32, person.Y())
	assert.Equal(t, math.MaxInt32, person.Z())
	assert.Equal(t, math.MaxInt32, person.Gold())
	assert.Equal(t, 1000, person.Mana())
	assert.Equal(t, 1000, person.Health())
	assert.Equal(t, 10, person.Respect())
	assert.Equal(t, 10, person.Strength())
	assert.Equal(t, 10, person.Experience())
	assert.Equal(t, 10, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasGun())
	assert.True(t, person.HasFamilty())
	assert.Equal(t, WarriorGamePersonType, person.Type())
}

func TestGamePersonNegativeCoordinates(t *testing.T) {
	person := NewGamePerson(
		WithCoordinates(-100, -200, -300),
	)

	assert.Equal(t, -100, person.X())
	assert.Equal(t, -200, person.Y())
	assert.Equal(t, -300, person.Z())
}

func TestGamePersonNameEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{"", ""},
		{"A", "A"},
		{"Hello", "Hello"},
		{"test\x00rest", "test"},
		{"test\x00", "test"},
		{"\x00test", ""},
		{"very_long_name_that_exceeds_the_maximum_length_of_forty_two_characters", "very_long_name_that_exceeds_the_maximum_le"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			person := NewGamePerson(WithName(tc.name))
			assert.Equal(t, tc.expected, person.Name())
		})
	}
}

func TestGamePersonTypeEdgeCases(t *testing.T) {
	testCases := []struct {
		personType int
		expected   int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 0},
		{5, 1},
		{6, 2},
		{7, 3},
		{8, 0},
		{-1, 3},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Type_%d", tc.personType), func(t *testing.T) {
			person := NewGamePerson(WithType(tc.personType))
			assert.Equal(t, tc.expected, person.Type())
		})
	}
}

func TestGamePersonFlagsCombinations(t *testing.T) {
	flagCombinations := []struct {
		house, gun, family bool
	}{
		{false, false, false},
		{true, false, false},
		{false, true, false},
		{false, false, true},
		{true, true, false},
		{true, false, true},
		{false, true, true},
		{true, true, true},
	}

	for i, combo := range flagCombinations {
		t.Run(fmt.Sprintf("Flags_%d", i), func(t *testing.T) {
			var options []Option
			if combo.house {
				options = append(options, WithHouse())
			}
			if combo.gun {
				options = append(options, WithGun())
			}
			if combo.family {
				options = append(options, WithFamily())
			}

			person := NewGamePerson(options...)
			assert.Equal(t, combo.house, person.HasHouse())
			assert.Equal(t, combo.gun, person.HasGun())
			assert.Equal(t, combo.family, person.HasFamilty())
		})
	}
}

func TestGamePersonManaHealthEdgeCases(t *testing.T) {
	testCases := []struct {
		mana, health                 int
		expectedMana, expectedHealth int
	}{
		{0, 0, 0, 0},
		{1000, 1000, 1000, 1000},
		{-1, 0, 0, 0},
		{0, -1, 0, 0},
		{1001, 1000, 1000, 1000},
		{1000, 1001, 1000, 1000},
		{500, 500, 500, 500},
		{999, 999, 999, 999},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("ManaHealth_%d", i), func(t *testing.T) {
			person := NewGamePerson(
				WithMana(tc.mana),
				WithHealth(tc.health),
			)
			assert.Equal(t, tc.expectedMana, person.Mana())
			assert.Equal(t, tc.expectedHealth, person.Health())
		})
	}
}

func TestGamePersonAttributesEdgeCases(t *testing.T) {
	testCases := []struct {
		respect, strength, experience, level                                 int
		expectedRespect, expectedStrength, expectedExperience, expectedLevel int
	}{
		{0, 0, 0, 0, 0, 0, 0, 0},
		{10, 10, 10, 10, 10, 10, 10, 10},
		{-1, 0, 0, 0, 0, 0, 0, 0},
		{0, -1, 0, 0, 0, 0, 0, 0},
		{0, 0, -1, 0, 0, 0, 0, 0},
		{0, 0, 0, -1, 0, 0, 0, 0},
		{11, 10, 10, 10, 10, 10, 10, 10},
		{10, 11, 10, 10, 10, 10, 10, 10},
		{10, 10, 11, 10, 10, 10, 10, 10},
		{10, 10, 10, 11, 10, 10, 10, 10},
		{5, 5, 5, 5, 5, 5, 5, 5},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Attributes_%d", i), func(t *testing.T) {
			person := NewGamePerson(
				WithRespect(tc.respect),
				WithStrength(tc.strength),
				WithExperience(tc.experience),
				WithLevel(tc.level),
			)
			assert.Equal(t, tc.expectedRespect, person.Respect())
			assert.Equal(t, tc.expectedStrength, person.Strength())
			assert.Equal(t, tc.expectedExperience, person.Experience())
			assert.Equal(t, tc.expectedLevel, person.Level())
		})
	}
}

func TestGamePersonMemoryLayout(t *testing.T) {
	var person GamePerson

	coordsSize := unsafe.Sizeof(person.Coordinates)
	goldSize := unsafe.Sizeof(person.gold)
	manaHealthSize := unsafe.Sizeof(person.manaHealth)
	paramsSize := unsafe.Sizeof(person.params)
	attrsSize := unsafe.Sizeof(person.attrs)
	nameSize := unsafe.Sizeof(person.name)

	totalExpectedSize := coordsSize + goldSize + manaHealthSize + paramsSize + attrsSize + nameSize
	actualSize := unsafe.Sizeof(person)

	t.Logf("Expected total size: %d, Actual size: %d", totalExpectedSize, actualSize)

	assert.GreaterOrEqual(t, actualSize, totalExpectedSize)
}

func TestGamePersonCopy(t *testing.T) {
	original := NewGamePerson(
		WithName("Original"),
		WithCoordinates(100, 200, 300),
		WithGold(500),
		WithMana(75),
		WithHealth(80),
		WithRespect(5),
		WithStrength(7),
		WithExperience(3),
		WithLevel(4),
		WithHouse(),
		WithGun(),
		WithType(BlacksmithGamePersonType),
	)

	copied := original

	assert.Equal(t, original.Name(), copied.Name())
	assert.Equal(t, original.X(), copied.X())
	assert.Equal(t, original.Y(), copied.Y())
	assert.Equal(t, original.Z(), copied.Z())
	assert.Equal(t, original.Gold(), copied.Gold())
	assert.Equal(t, original.Mana(), copied.Mana())
	assert.Equal(t, original.Health(), copied.Health())
	assert.Equal(t, original.Respect(), copied.Respect())
	assert.Equal(t, original.Strength(), copied.Strength())
	assert.Equal(t, original.Experience(), copied.Experience())
	assert.Equal(t, original.Level(), copied.Level())
	assert.Equal(t, original.HasHouse(), copied.HasHouse())
	assert.Equal(t, original.HasGun(), copied.HasGun())
	assert.Equal(t, original.HasFamilty(), copied.HasFamilty())
	assert.Equal(t, original.Type(), copied.Type())
}
