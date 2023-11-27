package main

import (
	"errors"
	"math"

	t "github.com/SpaceDiverr/util/ternary"
)

// -------- My implementation of StringBuilder --------

const (
	defaultLenAndCap = 1 << 8
)

type StringBuilder struct {
	str        []rune
	i          int
	isDefault  bool
	defaultCap int

	_ iStringBuilder
}

type iStringBuilder interface {
	iModifier
	iGetter
	iWriter
	iGrower
	iCapacityLackChecker
}

type iModifier interface {
	iGrower
	iResetter
}

type iGetter interface {
	String() string
	Cap() int
	Len() int
}

type iWriter interface {
	WriteRune(r rune)
	WriteString(s string)
}

type iGrower interface {
	GrowWithRate(rate float64)
	GrowBy(capacity int)
}

type iResetter interface {
	Reset()
}

type iCapacityLackChecker interface {
	isLackCapacity([]rune)
}

func NewDefault() *StringBuilder {

	return &StringBuilder{
		str:        make([]rune, defaultLenAndCap),
		i:          0,
		isDefault:  true,
		defaultCap: defaultLenAndCap,
	}
}

// Panics if сapacity < 0.
func NewWithCap(capacity uint) *StringBuilder {
	return &StringBuilder{
		str:        make([]rune, capacity),
		i:          0,
		isDefault:  false,
		defaultCap: int(capacity),
	}
}

func (sb *StringBuilder) WriteRune(r rune) *StringBuilder {
	sb.updateLenToCap()
	defer func() { sb.i++ }()

	if sb.isLackCapacity([]rune{r}) {
		sb.str = append(sb.str, r)
		return sb
	}
	sb.str[sb.i] = r
	return sb
}

func (sb *StringBuilder) WriteString(s string) *StringBuilder {

	if len(s) == 1 {
		sb.WriteRune(rune(s[0]))
		return sb
	}

	toWrite := []rune(s)
	if sb.isLackCapacity(toWrite) {
		sb.str = append(sb.str, toWrite...)
		sb.i += len(toWrite)
		return sb
	}

	sb.updateLenToCap()

	for i := range toWrite {
		sb.str[sb.i+i] = toWrite[i]
	}
	sb.i += len(toWrite)

	return sb
}

func (sb *StringBuilder) String() string {
	return string(sb.str[:sb.i])
}

func (sb *StringBuilder) Len() int {
	return sb.i
}

func (sb *StringBuilder) Cap() int {
	return cap(sb.str)
}

func (sb *StringBuilder) Reset() *StringBuilder {
	sb = t.Ternary[*StringBuilder](sb.isDefault, NewDefault(), NewWithCap(uint(sb.defaultCap)))
	return sb
}

func (sb *StringBuilder) GrowBy(capacity int) {
	if capacity < 1 {
		panic("stringbuilder: Capacity must be >= 1. StringBuilder failed to initialize with negative capacity.")
	}
	buf := make([]rune, len(sb.str)+capacity)
	copy(buf, sb.str)
	sb.str = buf
}

func (sb *StringBuilder) GrowWithRate(rate float64) error {
	if rate < 0 {
		return errors.New("stringbuilder: 'rate' must be >= 0. StringBuilder failed to grow capacity with negative rate (rate < 0) ")
	}
	capToAdd := int(math.Round(float64(cap(sb.str)) * rate))
	if capToAdd == cap(sb.str) {
		return errors.New("stringbuilder: 'rate' is too small; failed to grow")
	}

	sb.GrowBy(capToAdd)
	return nil
}

func (sb *StringBuilder) isLackCapacity(toAppend []rune) bool {
	return cap(sb.str) < sb.i+len(toAppend)
}

func (sb *StringBuilder) updateLenToCap() {
	sb.str = sb.str[:cap(sb.str)]
}