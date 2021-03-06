package size_test

import (
	"testing"

	. "github.com/openacid/low/size"
	"github.com/stretchr/testify/require"
)

const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

type intReader interface {
	Read() int
}
type intRead int32

func (m *intRead) Read() int {
	return 1
}

var myRead intRead = 0
var myReader intReader = &myRead
var myReaderNil intReader = nil

func TestSizeOf(t *testing.T) {

	ta := require.New(t)

	boolV := true

	cases := []struct {
		input interface{}
		want  int
	}{
		{true, 1},
		{boolV, 1},
		{&boolV, 8 + 1},
		{"abc", 16 + 3},

		{uint8(0), 1}, {int8(0), 1},
		{uint16(0), 2}, {int16(0), 2},
		{uint32(0), 4}, {int32(0), 4},
		{uint64(0), 8}, {int64(0), 8},

		{float32(0), 4},
		{float64(0), 8},

		{complex64(complex(1, 2)), 8},
		{complex128(complex(1, 2)), 16},
		{int(0), UintSize / 8},

		{[]int32{1, 2}, 24 + 8},
		{[3]int32{1, 2, 3}, 12},

		{map[int32]string{1: "a", 2: "b"}, 8 + (4 + (16 + 1)) + (4 + (16 + 1))},

		{struct{ a, b int64 }{1, 2}, 16},
		{&struct{ a, b int64 }{1, 2}, 8 + 16},

		{myReaderNil, 0},
		{myReader, 8 + 4},
		{struct{ a intReader }{nil}, 16},
		{struct{ a intReader }{&myRead}, 16 + 8 + 4},
	}

	for i, c := range cases {
		rst := Of(c.input)
		ta.Equal(c.want, rst, "%d-th: input: %+v", i+1, c.input)
	}
}

func TestSizeStat(t *testing.T) {

	ta := require.New(t)

	type my struct {
		a []int32
		b [3]int32
		c map[string]int8
		d *my
		e []*my
		f []string
		g intReader
		h intReader
	}

	v := my{
		a: []int32{1, 2, 3},
		b: [3]int32{4, 5, 6},
		c: map[string]int8{
			"abc": 3,
		},
		d: &my{
			a: []int32{1, 2},
		},
		e: []*my{
			{
				a: []int32{1, 2, 3},
			},
			{
				a: []int32{2, 3, 4},
			},
		},
		f: []string{
			"abc",
			"def",
		},
		g: nil,
		h: &myRead,
	}

	want10 := `
size_test.my: 658
    a: []int32: 36
        0: int32: 4
        1: int32: 4
        2: int32: 4
    b: [3]int32: 12
        0: int32: 4
        1: int32: 4
        2: int32: 4
    c: map[string]int8: 28
        abc: int8: 1
    d: *size_test.my: 148
        size_test.my: 140
            a: []int32: 32
                0: int32: 4
                1: int32: 4
            b: [3]int32: 12
                0: int32: 4
                1: int32: 4
                2: int32: 4
            c: map[string]int8: 8
            d: *size_test.my: 8
            e: []*size_test.my: 24
            f: []string: 24
            g: size_test.intReader: 16
                <nil>
            h: size_test.intReader: 16
                <nil>
    e: []*size_test.my: 328
        0: *size_test.my: 152
            size_test.my: 144
                a: []int32: 36
                    0: int32: 4
                    1: int32: 4
                    2: int32: 4
                b: [3]int32: 12
                    0: int32: 4
                    1: int32: 4
                    2: int32: 4
                c: map[string]int8: 8
                d: *size_test.my: 8
                e: []*size_test.my: 24
                f: []string: 24
                g: size_test.intReader: 16
                    <nil>
                h: size_test.intReader: 16
                    <nil>
        1: *size_test.my: 152
            size_test.my: 144
                a: []int32: 36
                    0: int32: 4
                    1: int32: 4
                    2: int32: 4
                b: [3]int32: 12
                    0: int32: 4
                    1: int32: 4
                    2: int32: 4
                c: map[string]int8: 8
                d: *size_test.my: 8
                e: []*size_test.my: 24
                f: []string: 24
                g: size_test.intReader: 16
                    <nil>
                h: size_test.intReader: 16
                    <nil>
    f: []string: 62
        0: string: 19
        1: string: 19
    g: size_test.intReader: 16
        <nil>
    h: size_test.intReader: 28
        *size_test.intRead: 12
            size_test.intRead: 4`[1:]

	got10 := Stat(v, 10, 100)
	ta.Equal(want10, got10)

	want3 := `
size_test.my: 658
    a: []int32: 36
        0: int32: 4
        1: int32: 4
        2: int32: 4
    b: [3]int32: 12
        0: int32: 4
        1: int32: 4
        2: int32: 4
    c: map[string]int8: 28
        abc: int8: 1
    d: *size_test.my: 148
        size_test.my: 140
            a: []int32: 32
            b: [3]int32: 12
            c: map[string]int8: 8
            d: *size_test.my: 8
            e: []*size_test.my: 24
            f: []string: 24
            g: size_test.intReader: 16
            h: size_test.intReader: 16
    e: []*size_test.my: 328
        0: *size_test.my: 152
            size_test.my: 144
        1: *size_test.my: 152
            size_test.my: 144
    f: []string: 62
        0: string: 19
        1: string: 19
    g: size_test.intReader: 16
        <nil>
    h: size_test.intReader: 28
        *size_test.intRead: 12
            size_test.intRead: 4`[1:]
	got3 := Stat(v, 3, 100)
	ta.Equal(want3, got3)

	want32 := `
size_test.my: 658
    a: []int32: 36
        0: int32: 4
        1: int32: 4
    b: [3]int32: 12
        0: int32: 4
        1: int32: 4
    c: map[string]int8: 28
        abc: int8: 1
    d: *size_test.my: 148
        size_test.my: 140
            a: []int32: 32
            b: [3]int32: 12
            c: map[string]int8: 8
            d: *size_test.my: 8
            e: []*size_test.my: 24
            f: []string: 24
            g: size_test.intReader: 16
            h: size_test.intReader: 16
    e: []*size_test.my: 328
        0: *size_test.my: 152
            size_test.my: 144
        1: *size_test.my: 152
            size_test.my: 144
    f: []string: 62
        0: string: 19
        1: string: 19
    g: size_test.intReader: 16
        <nil>
    h: size_test.intReader: 28
        *size_test.intRead: 12
            size_test.intRead: 4`[1:]
	got32 := Stat(v, 3, 2)
	ta.Equal(want32, got32)
}

func TestSizeStat_AvgOf(t *testing.T) {

	ta := require.New(t)

	type my struct {
		a []int32
		b [3]int32
		c map[string]int8
		d *my
		e []*my
		f []string
		g intReader
		h intReader
	}

	v := my{
		a: []int32{1, 2, 3},
		b: [3]int32{4, 5, 6},
		c: map[string]int8{
			"abc": 3,
		},
		d: &my{
			a: []int32{1, 2},
		},
		e: []*my{
			{
				a: []int32{1, 2, 3},
			},
			{
				a: []int32{2, 3, 4},
			},
		},
		f: []string{
			"abc",
			"def",
		},
		g: nil,
		h: &myRead,
	}

	want10 := `
size_test.my: 658 /n = 65.800
    a: []int32: 36 /n = 3.600
        0: int32: 4 /n = 0.400
        1: int32: 4 /n = 0.400
        2: int32: 4 /n = 0.400
    b: [3]int32: 12 /n = 1.200
        0: int32: 4 /n = 0.400
        1: int32: 4 /n = 0.400
        2: int32: 4 /n = 0.400
    c: map[string]int8: 28 /n = 2.800
        abc: int8: 1 /n = 0.100
    d: *size_test.my: 148 /n = 14.800
        size_test.my: 140 /n = 14.000
            a: []int32: 32 /n = 3.200
                0: int32: 4 /n = 0.400
                1: int32: 4 /n = 0.400
            b: [3]int32: 12 /n = 1.200
                0: int32: 4 /n = 0.400
                1: int32: 4 /n = 0.400
                2: int32: 4 /n = 0.400
            c: map[string]int8: 8 /n = 0.800
            d: *size_test.my: 8 /n = 0.800
            e: []*size_test.my: 24 /n = 2.400
            f: []string: 24 /n = 2.400
            g: size_test.intReader: 16 /n = 1.600
                <nil>
            h: size_test.intReader: 16 /n = 1.600
                <nil>
    e: []*size_test.my: 328 /n = 32.800
        0: *size_test.my: 152 /n = 15.200
            size_test.my: 144 /n = 14.400
                a: []int32: 36 /n = 3.600
                    0: int32: 4 /n = 0.400
                    1: int32: 4 /n = 0.400
                    2: int32: 4 /n = 0.400
                b: [3]int32: 12 /n = 1.200
                    0: int32: 4 /n = 0.400
                    1: int32: 4 /n = 0.400
                    2: int32: 4 /n = 0.400
                c: map[string]int8: 8 /n = 0.800
                d: *size_test.my: 8 /n = 0.800
                e: []*size_test.my: 24 /n = 2.400
                f: []string: 24 /n = 2.400
                g: size_test.intReader: 16 /n = 1.600
                    <nil>
                h: size_test.intReader: 16 /n = 1.600
                    <nil>
        1: *size_test.my: 152 /n = 15.200
            size_test.my: 144 /n = 14.400
                a: []int32: 36 /n = 3.600
                    0: int32: 4 /n = 0.400
                    1: int32: 4 /n = 0.400
                    2: int32: 4 /n = 0.400
                b: [3]int32: 12 /n = 1.200
                    0: int32: 4 /n = 0.400
                    1: int32: 4 /n = 0.400
                    2: int32: 4 /n = 0.400
                c: map[string]int8: 8 /n = 0.800
                d: *size_test.my: 8 /n = 0.800
                e: []*size_test.my: 24 /n = 2.400
                f: []string: 24 /n = 2.400
                g: size_test.intReader: 16 /n = 1.600
                    <nil>
                h: size_test.intReader: 16 /n = 1.600
                    <nil>
    f: []string: 62 /n = 6.200
        0: string: 19 /n = 1.900
        1: string: 19 /n = 1.900
    g: size_test.intReader: 16 /n = 1.600
        <nil>
    h: size_test.intReader: 28 /n = 2.800
        *size_test.intRead: 12 /n = 1.200
            size_test.intRead: 4 /n = 0.400`[1:]

	got10 := Stat(v, 10, 100, Opt{AvgOf: 10})
	ta.Equal(want10, got10)

	want3 := `
size_test.my: 658 /n = 65.800
    a: []int32: 36 /n = 3.600
        0: int32: 4 /n = 0.400
        1: int32: 4 /n = 0.400
        2: int32: 4 /n = 0.400
    b: [3]int32: 12 /n = 1.200
        0: int32: 4 /n = 0.400
        1: int32: 4 /n = 0.400
        2: int32: 4 /n = 0.400
    c: map[string]int8: 28 /n = 2.800
        abc: int8: 1 /n = 0.100
    d: *size_test.my: 148 /n = 14.800
        size_test.my: 140 /n = 14.000
            a: []int32: 32 /n = 3.200
            b: [3]int32: 12 /n = 1.200
            c: map[string]int8: 8 /n = 0.800
            d: *size_test.my: 8 /n = 0.800
            e: []*size_test.my: 24 /n = 2.400
            f: []string: 24 /n = 2.400
            g: size_test.intReader: 16 /n = 1.600
            h: size_test.intReader: 16 /n = 1.600
    e: []*size_test.my: 328 /n = 32.800
        0: *size_test.my: 152 /n = 15.200
            size_test.my: 144 /n = 14.400
        1: *size_test.my: 152 /n = 15.200
            size_test.my: 144 /n = 14.400
    f: []string: 62 /n = 6.200
        0: string: 19 /n = 1.900
        1: string: 19 /n = 1.900
    g: size_test.intReader: 16 /n = 1.600
        <nil>
    h: size_test.intReader: 28 /n = 2.800
        *size_test.intRead: 12 /n = 1.200
            size_test.intRead: 4 /n = 0.400`[1:]
	got3 := Stat(v, 3, 100, Opt{AvgOf: 10})
	ta.Equal(want3, got3)

	want32 := `
size_test.my: 658 /n = 65.800
    a: []int32: 36 /n = 3.600
        0: int32: 4 /n = 0.400
        1: int32: 4 /n = 0.400
    b: [3]int32: 12 /n = 1.200
        0: int32: 4 /n = 0.400
        1: int32: 4 /n = 0.400
    c: map[string]int8: 28 /n = 2.800
        abc: int8: 1 /n = 0.100
    d: *size_test.my: 148 /n = 14.800
        size_test.my: 140 /n = 14.000
            a: []int32: 32 /n = 3.200
            b: [3]int32: 12 /n = 1.200
            c: map[string]int8: 8 /n = 0.800
            d: *size_test.my: 8 /n = 0.800
            e: []*size_test.my: 24 /n = 2.400
            f: []string: 24 /n = 2.400
            g: size_test.intReader: 16 /n = 1.600
            h: size_test.intReader: 16 /n = 1.600
    e: []*size_test.my: 328 /n = 32.800
        0: *size_test.my: 152 /n = 15.200
            size_test.my: 144 /n = 14.400
        1: *size_test.my: 152 /n = 15.200
            size_test.my: 144 /n = 14.400
    f: []string: 62 /n = 6.200
        0: string: 19 /n = 1.900
        1: string: 19 /n = 1.900
    g: size_test.intReader: 16 /n = 1.600
        <nil>
    h: size_test.intReader: 28 /n = 2.800
        *size_test.intRead: 12 /n = 1.200
            size_test.intRead: 4 /n = 0.400`[1:]
	got32 := Stat(v, 3, 2, Opt{AvgOf: 10})
	ta.Equal(want32, got32)
}

func TestSizeStat_AvgUnit(t *testing.T) {

	ta := require.New(t)

	type my struct {
		a []int32
		b [3]int32
		c map[string]int8
		d *my
		e []*my
		f []string
		g intReader
		h intReader
	}

	v := my{
		a: []int32{1, 2, 3},
		b: [3]int32{4, 5, 6},
		c: map[string]int8{
			"abc": 3,
		},
		d: &my{
			a: []int32{1, 2},
		},
		e: []*my{
			{
				a: []int32{1, 2, 3},
			},
			{
				a: []int32{2, 3, 4},
			},
		},
		f: []string{
			"abc",
			"def",
		},
		g: nil,
		h: &myRead,
	}

	want32 := `
size_test.my: 658 /n = 526.400
    a: []int32: 36 /n = 28.800
        0: int32: 4 /n = 3.200
        1: int32: 4 /n = 3.200
    b: [3]int32: 12 /n = 9.600
        0: int32: 4 /n = 3.200
        1: int32: 4 /n = 3.200
    c: map[string]int8: 28 /n = 22.400
        abc: int8: 1 /n = 0.800
    d: *size_test.my: 148 /n = 118.400
        size_test.my: 140 /n = 112.000
            a: []int32: 32 /n = 25.600
            b: [3]int32: 12 /n = 9.600
            c: map[string]int8: 8 /n = 6.400
            d: *size_test.my: 8 /n = 6.400
            e: []*size_test.my: 24 /n = 19.200
            f: []string: 24 /n = 19.200
            g: size_test.intReader: 16 /n = 12.800
            h: size_test.intReader: 16 /n = 12.800
    e: []*size_test.my: 328 /n = 262.400
        0: *size_test.my: 152 /n = 121.600
            size_test.my: 144 /n = 115.200
        1: *size_test.my: 152 /n = 121.600
            size_test.my: 144 /n = 115.200
    f: []string: 62 /n = 49.600
        0: string: 19 /n = 15.200
        1: string: 19 /n = 15.200
    g: size_test.intReader: 16 /n = 12.800
        <nil>
    h: size_test.intReader: 28 /n = 22.400
        *size_test.intRead: 12 /n = 9.600
            size_test.intRead: 4 /n = 3.200`[1:]

	got32 := Stat(v, 3, 2, Opt{AvgOf: 10, AvgUnit: 1.0 / 8})
	ta.Equal(want32, got32)
}
