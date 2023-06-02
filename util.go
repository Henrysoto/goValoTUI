package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var rank = []string{
	"Unranked",
	"Iron 1",
	"Iron 2",
	"Iron 3",
	"Bronze 1",
	"Bronze 2",
	"Bronze 3",
	"Silver 1",
	"Silver 2",
	"Silver 3",
	"Gold 1",
	"Gold 2",
	"Gold 3",
	"Platinum 1",
	"Platinum 2",
	"Platinum 3",
	"Diamond 1",
	"Diamond 2",
	"Diamond 3",
	"Ascendant 1",
	"Ascendant 2",
	"Ascendant 3",
	"Immortal 1",
	"Immortal 2",
	"Immortal 3",
	"Radiant",
}

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GetUsername(name string, tag string) string {
	return fmt.Sprintf("%s#%s", name, tag)
}

func RandomInt(min, max int64) int64 {
	return min + rnd.Int63n(max-min+1)
}

func RandomString(length int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomPlayerName() string {
	return RandomString(8)
}

func RandomPlayerRank() string {
	return rank[RandomInt(0, int64(len(rank)))]
}

func FormatStringTable(content string, length int) string {
	for len(content) < length {
		content = content + " "
	}
	return content
}
