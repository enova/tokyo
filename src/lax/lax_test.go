package lax

import (
	"github.com/mgutz/ansi"
	"github.com/stretchr/testify/assert"
	"math"
	"strconv"
	"testing"
	"time"
)

func eqF32(x float32, y float32) bool {
	return math.Abs(float64(y)-float64(x)) < 1e-10
}

func TestSeparator(t *testing.T) {
	dashes := Separator(5)
	assert.Equal(t, dashes, "-----", "Separator(5) didn't produce 5 dashes: "+dashes)
}

func TestHeading(t *testing.T) {
	heading := Heading("Welcome", 5)
	assert.Equal(t, heading, "----- Welcome -----", "Heading(5) didn't produce ----- Welcome -----: "+heading)
}

func TestReadFile(t *testing.T) {
	contents := ReadFile("test.txt", "test")
	assert.Equal(t, contents, "Hello!\n", "ReadFile did not read contents of file test.txt: "+contents)
}

func TestFileMode(t *testing.T) {
	assert := assert.New(t)

	fileMode, _ := FileMode("test.txt")
	assert.True(fileMode.IsRegular(), "FileMode failed to recognize a file")

	dirMode, _ := FileMode("../")
	assert.True(dirMode.IsDir(), "FileMode failed to recognize a directory")

	_, err := FileMode("i_dont_exist")
	assert.NotNil(err)
}

func TestIsFile(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsDir("./"))
	assert.False(IsFile("./"))
}

func TestIsDir(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsFile("lax_test.go"))
	assert.False(IsDir("lax_test.go"))
}

// Test Struct (for YAML Parsing)
type sale struct {
	Seller struct {
		Name string
		Age  int
	}

	Quantity int
	Price    float32
}

func TestReadYAML(t *testing.T) {
	assert := assert.New(t)

	s := sale{}
	ReadYAML("test.yml", &s)

	assert.Equal(s.Seller.Name, "Chip", "ReadYAML failed to read Seller-Name: "+s.Seller.Name)
	assert.Equal(s.Seller.Age, 25, "ReadYAML failed to read Seller-Age: "+strconv.Itoa(s.Seller.Age))
	assert.Equal(s.Quantity, 1800, "ReadYAML failed to read Quantity: "+strconv.Itoa(s.Quantity))
	assert.True(eqF32(s.Price, 25.43), "ReadYAML failed to read Price: "+F32toa(s.Price))
}

func TestStripColors(t *testing.T) {
	hello := ansi.Color("hello", "red")
	world := ansi.Color("world", "green")
	assert.Equal(t, StripColors(hello+" "+world), "hello world", "StripColros failed to remove color codes")
}

func TestMakeDate(t *testing.T) {
	d := MakeDate(20150601)
	assert.Equal(t, d.Format("20060102"), "20150601", "MakeDate failed to convert 20150601 into a date: "+d.Format("20060102"))
}

func TestMakeDateI(t *testing.T) {
	d := MakeDateI(20150601)
	assert.Equal(t, d.Format("20060102 15:04:05 -0700 MST"), "20150601 00:00:00 +0000 UTC", "MakeDate failed to convert 20150601 into a date: "+d.Format("20060102"))
}

func TestMakeDateS(t *testing.T) {
	d := MakeDateS("2015-06-01")
	assert.Equal(t, d.Format("20060102 15:04:05 -0700 MST"), "20150601 00:00:00 +0000 UTC", "MakeDate failed to convert 20150601 into a date: "+d.Format("20060102"))

	d = MakeDateS("20150601")
	assert.Equal(t, d.Format("20060102 15:04:05 -0700 MST"), "20150601 00:00:00 +0000 UTC", "MakeDate failed to convert 20150601 into a date: "+d.Format("20060102"))
}

func TestAddDays(t *testing.T) {
	d := MakeDate(20150601)
	assert.Equal(t, AddDays(d, 0), MakeDate(20150601), "AddDate failed to add zero days")
	assert.Equal(t, AddDays(d, 2), MakeDate(20150603), "AddDate failed to add two days")
	assert.Equal(t, AddDays(d, -2), MakeDate(20150530), "AddDate failed to subtract two days")
}

func TestSubDays(t *testing.T) {
	d := MakeDate(20150601)
	assert.Equal(t, SubDays(d, 0), MakeDate(20150601), "SubDate failed to subtract zero days")
	assert.Equal(t, SubDays(d, 2), MakeDate(20150530), "SubDate failed to subtract two days")
	assert.Equal(t, SubDays(d, -2), MakeDate(20150603), "SubDate failed to add two days")
}

func TestInTimeRange(t *testing.T) {
	assert := assert.New(t)

	a := MakeDate(20150601)
	b := MakeDate(20150602)
	c := MakeDate(20150603)

	assert.False(InTimeRange(a, b, c), "Time (a) should not be between (b) and (c)")
	assert.True(InTimeRange(b, a, c), "Time (b) should be between (a) and (c)")
	assert.False(InTimeRange(c, a, b), "Time (c) should not be between (a) and (b)")
	assert.True(InTimeRange(a, a, b), "Time (a) should be between (a) and (b)")
	assert.True(InTimeRange(b, a, b), "Time (b) should be between (a) and (b)")
}

func TestDaysBetween(t *testing.T) {
	a := MakeDate(20150601)
	c := MakeDate(20150603)

	assert.Equal(t, DaysBetween(a, a), 0, "DaysBetween should be 0")
	assert.Equal(t, DaysBetween(a, c), 2, "DaysBetween should be 3")
	assert.Equal(t, DaysBetween(c, a), -2, "DaysBetween should be -3")
}

func TestMakeDateList(t *testing.T) {
	assert := assert.New(t)

	bgn := MakeDate(20150601)
	end := MakeDate(20150603)
	list := MakeDateList(bgn, end)

	assert.Equal(len(list), 3, "List should have three dates")
	assert.Equal(list[0], bgn, "First date is wrong")
	assert.Equal(list[1], bgn.Add(24*time.Hour), "Second date is wrong")
	assert.Equal(list[2], bgn.Add(48*time.Hour), "Third date is wrong")
}

func TestSort(t *testing.T) {
	assert := assert.New(t)

	s := []string{"B", "C", "A"}
	SortS(s)

	assert.Equal(len(s), 3, "Sort failed")
	assert.Equal(s[0], "A", "Sort failed")
	assert.Equal(s[1], "B", "Sort failed")
	assert.Equal(s[2], "C", "Sort failed")

	i := []int{2, 3, 1}
	SortI(i)

	assert.Equal(len(i), 3, "Sort failed")
	assert.Equal(i[0], 1, "Sort failed")
	assert.Equal(i[1], 2, "Sort failed")
	assert.Equal(i[2], 3, "Sort failed")
}

func TestParseFloat64(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(ParseFloat64("23.45"), 23.45, "Failed to parse float64")
}

func TestSideBySide(t *testing.T) {
	assert := assert.New(t)

	left := "Bear\n"
	left += "Cat\n"
	left += "Elephant\n"

	right := "Ant\n"
	right += "Bee\n"

	result := "Bear      Ant\n"
	result += "Cat       Bee\n"
	result += "Elephant  \n"

	assert.Equal(SideBySide(left, right, "  "), result)
}

func TestEnv(t *testing.T) {
	assert := assert.New(t)

	assert.NotEqual(Username(), "unknown_user")
	assert.NotEqual(Hostname(), "unknown_host")
	assert.NotEqual(ExecPath(), "unknown_exec")
}

func TestHasS(t *testing.T) {
	assert := assert.New(t)

	list := []string{"a", "b", "c"}

	assert.True(HasS(list, "a"))
	assert.True(HasS(list, "b"))
	assert.True(HasS(list, "c"))
	assert.False(HasS(list, "x"))
}

func TestAllFilesUnderDir(t *testing.T) {
	assert := assert.New(t)

	files := AllFilesUnderDir("dir")
	SortS(files)

	assert.Equal(4, len(files))
	assert.Equal("dir/dirA/fileC", files[0])
	assert.Equal("dir/dirA/fileD", files[1])
	assert.Equal("dir/fileA", files[2])
	assert.Equal("dir/fileB", files[3])
}
