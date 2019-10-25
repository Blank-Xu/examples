package ftp

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"
)

var (
	commandMap = map[string]bool{
		"ALLO": true,
		"CDUP": true,
		"CWD":  true,
		"DELE": true,
		"EPRT": true,
		"EPSV": true,
		"FEAT": true,
		"LIST": true,
		"NLST": true,
		"MDTM": true,
		"MKD":  true,
		"MODE": true,
		"NOOP": true,
		"OPTS": true,
		"PASS": true,
		"PASV": true,
		"PORT": true,
		"PWD":  true,
		"QUIT": true,
		"RETR": true,
		"RNFR": true,
		"RNTO": true,
		"RMD":  true,
		"SIZE": true,
		"STOR": true,
		"STRU": true,
		"SYST": true,
		"TYPE": true,
		"USER": true,
		"XCUP": true,
		"XCWD": true,
		"XPWD": true,
		"XRMD": true,
	}

	commandSlice = []string{
		"ALLO",
		"CDUP",
		"CWD",
		"DELE",
		"EPRT",
		"EPSV",
		"FEAT",
		"LIST",
		"NLST",
		"MDTM",
		"MKD",
		"MODE",
		"NOOP",
		"OPTS",
		"PASS",
		"PASV",
		"PORT",
		"PWD",
		"QUIT",
		"RETR",
		"RNFR",
		"RNTO",
		"RMD",
		"SIZE",
		"STOR",
		"STRU",
		"SYST",
		"TYPE",
		"USER",
		"XCUP",
		"XCWD",
		"XPWD",
		"XRMD",
	}

	l = len(commandSlice)

	commandByte = make([]int, 0, l)

	commandMap2 = make(map[int]bool, l)
)

func init() {
	sort.Slice(commandSlice, func(i, j int) bool {
		return commandSlice[i] < commandSlice[j]
	})

	for _, value := range commandSlice {
		i := stringToInt(value)
		commandByte = append(commandByte, i)
		commandMap2[i] = true
	}
}

func stringToInt(s string) int {
	var str string
	for _, b := range s {
		str += fmt.Sprint(b)
	}
	i, _ := strconv.Atoi(str)
	return i
}

func TestSlice(t *testing.T) {
	var s []int
	if len(s) == 0 {
		t.Log("0")
	} else {
		t.Log("1")
	}
}

func BenchmarkStringToInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sl = make([]int, 0, l)
		// for _, value := range commandSlice {
		i := stringToInt("PASV")
		sl = append(sl, i)
		// }
	}
}

func TestRune(t *testing.T) {

	fmt.Println("ab" < "bb")

	// for _, v := range commandSlice {
	// 	var a []int32
	// 	var s string
	// 	for _, r := range v {
	// 		a = append(a, int32(r))
	// 		// s += string(r)
	// 		fmt.Print(r)
	// 		s += fmt.Sprint(r)
	// 	}
	// 	fmt.Println()
	// 	fmt.Println(string(a))
	// 	fmt.Println(s)
	// 	fmt.Println()

	// r, _ := utf8.DecodeRuneInString(v)
	// t.Logf("%v\n", r)
	// }
}

// func BenchmarkSearch(b *testing.B) {
// 	BenchmarkMapSearch(b)
// 	BenchmarkSliceSearch(b)
// }

func BenchmarkMapSearch(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < b.N; i++ {
		idx := rnd.Int31n(int32(l))
		_, ok := commandMap[commandSlice[idx]]
		if !ok {
			b.Fatal("not found")
		}
	}
}

func BenchmarkSliceSearch(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < b.N; i++ {
		idx := rnd.Int31n(int32(l))
		n := sort.SearchStrings(commandSlice, commandSlice[idx])
		if n == l {
			b.Fatal("not found")
		}
	}
}

func BenchmarkByteSearch(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	// b.Log(l)
	for i := 0; i < b.N; i++ {
		idx := rnd.Int31n(int32(l))
		// b.Log(idx)
		n := sort.SearchInts(commandByte, stringToInt(commandSlice[idx]))
		if n == l {
			b.Fatal("not found")
		}
	}
}

func BenchmarkBinarySearch(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	// b.Log(l)
	for i := 0; i < b.N; i++ {
		idx := rnd.Int31n(int32(l))
		// if commandSlice[idx] > "PASV" {
		// 	idx = -1
		// }
		// b.Log(idx)
		num := BinarySearch(commandSlice, commandSlice[idx])
		if num == -1 {
			b.Fatal("not found")
		}
	}
}

func BinarySearch(sortedData []string, s string) int {
	var low, high, mid = 0, l - 1, 0
	for low <= high {
		mid = low + (high-low)/2
		var value = sortedData[mid]
		if value == s {
			return mid
		} else if value > s {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}

func BinarySearch2(sortedData []int, s int) int {
	var low, high, mid = 0, l - 1, 0
	for low <= high {
		mid = low + (high-low)/2
		var value = sortedData[mid]
		if value == s {
			return mid
		} else if value > s {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}
