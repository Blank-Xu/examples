package ftp

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

var testEqualSlice = []string{
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

var testEqualSliceByte [][]byte

func init() {
	testEqualSliceByte = make([][]byte, 0, len(testEqualSlice))
	for _, v := range testEqualSlice {
		b := []byte(v)
		testEqualSliceByte = append(testEqualSliceByte, b)
	}
}

func BenchmarkByte(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < b.N; i++ {
		idx := rnd.Int31n(int32(l))
		var s = []byte("PASV")
		_ = bytes.Equal(testEqualSliceByte[idx], s)
	}
}

func BenchmarkByteEqual(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < b.N; i++ {
		idx := rnd.Int31n(int32(l))
		var s = []byte("PASV")
		_ = bytes.EqualFold(testEqualSliceByte[idx], s)
	}
}

func BenchmarkString(b *testing.B) {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < b.N; i++ {
		idx := rnd.Int31n(int32(l))
		var s = "PASV"
		_ = string(testEqualSliceByte[idx]) == s
	}
}
