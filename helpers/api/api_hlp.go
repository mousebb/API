package api_helpers

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	EARTH               = 3963.1676 // radius of Earth in miles
	SOUTWEST_LATITUDE   = -90.00
	SOUTHWEST_LONGITUDE = -180.00
	NORTHEAST_LATITUDE  = 90.00
	NORTHEAST_LONGITUDE = 180.00
	CENTER_LATITUDE     = 44.79300
	CENTER_LONGITUDE    = -91.41048

	API_DOMAIN       = "https://API.curtmfg.com"
	AUTH_KEY_TYPE    = "AUTHENTICATION"
	PUBLIC_KEY_TYPE  = "PUBLIC"
	PRIVATE_KEY_TYPE = "PRIVATE"
)

func RandGenerator(max int) int {
	if max == 0 {
		return 0
	}
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return r.Intn(max)
}

func ValueOrFileContents(value string, filename string) string {
	if value != "" {
		return value
	}
	slurp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading %q: %v", filename, err)
	}
	return strings.TrimSpace(string(slurp))
}

func Md5Encrypt(str string) (string, error) {
	if str == "" {
		return "", errors.New("Invalid string parameter")
	}

	h := md5.New()
	io.WriteString(h, str)

	var buf bytes.Buffer
	_, err := buf.Write(h.Sum(nil))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func FromWindows1252(str string) string {
	var arr = []byte(str)
	var buf = bytes.NewBuffer(make([]byte, 512))
	var r rune

	for _, b := range arr {
		switch b {
		case 0x80:
			r = 0x20AC
		case 0x82:
			r = 0x201A
		case 0x83:
			r = 0x0192
		case 0x84:
			r = 0x201E
		case 0x85:
			r = 0x2026
		case 0x86:
			r = 0x2020
		case 0x87:
			r = 0x2021
		case 0x88:
			r = 0x02C6
		case 0x89:
			r = 0x2030
		case 0x8A:
			r = 0x0160
		case 0x8B:
			r = 0x2039
		case 0x8C:
			r = 0x0152
		case 0x8E:
			r = 0x017D
		case 0x91:
			r = 0x2018
		case 0x92:
			r = 0x2019
		case 0x93:
			r = 0x201C
		case 0x94:
			r = 0x201D
		case 0x95:
			r = 0x2022
		case 0x96:
			r = 0x2013
		case 0x97:
			r = 0x2014
		case 0x98:
			r = 0x02DC
		case 0x99:
			r = 0x2122
		case 0x9A:
			r = 0x0161
		case 0x9B:
			r = 0x203A
		case 0x9C:
			r = 0x0153
		case 0x9E:
			r = 0x017E
		case 0x9F:
			r = 0x0178
		default:
			r = rune(b)
		}

		buf.WriteRune(r)
	}

	return string(buf.Bytes())
}

func Escape(txt string) string {
	txt = escapeQuotes(txt)
	return escapeString(txt)
}

//https://github.com/ziutek/mymysql/blob/master/native/codecs.go#L462
func escapeString(txt string) string {
	var (
		esc string
		buf bytes.Buffer
	)
	last := 0
	for ii, bb := range txt {
		switch bb {
		case 0:
			esc = `\0`
		case '\n':
			esc = `\n`
		case '\r':
			esc = `\r`
		case '\\':
			esc = `\\`
		case '\'':
			esc = `\'`
		case '"':
			esc = `\"`
		case '\032':
			esc = `\Z`
		default:
			continue
		}
		io.WriteString(&buf, txt[last:ii])
		io.WriteString(&buf, esc)
		last = ii + 1
	}
	io.WriteString(&buf, txt[last:])
	return buf.String()
}

//https://github.com/ziutek/mymysql/blob/master/native/codecs.go#L495
func escapeQuotes(txt string) string {
	var buf bytes.Buffer
	last := 0
	for ii, bb := range txt {
		if bb == '\'' {
			io.WriteString(&buf, txt[last:ii])
			io.WriteString(&buf, `''`)
			last = ii + 1
		}
	}
	io.WriteString(&buf, txt[last:])
	return buf.String()
}
