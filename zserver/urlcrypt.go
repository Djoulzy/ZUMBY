package zserver

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/Djoulzy/Tools/clog"
)

type cypher struct {
	HashSize int // Should be 8
	HexKey   []byte
	HexIV    []byte
}

func (uc *cypher) genIVBin() []byte {
	ivBin := make([]byte, 16)

	newRand := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	for i := 0; i < 16; i++ {
		ivBin[i] = byte(newRand.Intn(256))
	}

	return ivBin
}

func (uc *cypher) genIVHex() []byte {
	bin := uc.genIVBin()

	dst := make([]byte, hex.EncodedLen(len(bin)))
	hex.Encode(dst, bin)
	return dst
}

func (*cypher) pkcs7pad(data []byte, blockSize int) []byte {

	var paddingCount int

	if paddingCount = blockSize - (len(data) % blockSize); paddingCount == 0 {
		paddingCount = blockSize
	}

	return append(data, bytes.Repeat([]byte{byte(paddingCount)}, paddingCount)...)
}

// RemovePkcs7 removes pkcs7 padding from previously padded byte array
func (uc *cypher) pkcs7unpad(padded []byte, blockSize int) []byte {

	dataLen := len(padded)
	paddingCount := int(padded[dataLen-1])

	if paddingCount > blockSize || paddingCount <= 0 {
		return padded //data is not padded (or not padded correctly), return as is
	}

	padding := padded[dataLen-paddingCount : dataLen-1]

	for _, b := range padding {
		if int(b) != paddingCount {
			return padded //data is not padded (or not padded correcly), return as is
		}
	}

	return padded[:len(padded)-paddingCount] //return data - padding
}

func (uc *cypher) encodeBase64(b []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(b))
}

func (uc *cypher) decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		clog.Fatal("Crypt", "decodeBase64", err)
	}
	return data
}

func (uc *cypher) getMD5Hash(text string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(text))
	md5 := hasher.Sum(nil)

	return md5
}

func (uc *cypher) encryptB64(text string) ([]byte, error) {
	var ivBin []byte

	keyBin := make([]byte, hex.DecodedLen(len(uc.HexKey)))
	hex.Decode(keyBin, uc.HexKey)

	if len(uc.HexIV) == 0 {
		ivBin = uc.genIVBin()
	} else {
		ivBin = make([]byte, hex.DecodedLen(len(uc.HexIV)))
		hex.Decode(ivBin, uc.HexIV)
	}

	block, err := aes.NewCipher(keyBin)
	if err != nil {
		return nil, err
	}

	textHash := (uc.getMD5Hash(text))[0:8]

	signedText := append(textHash, []byte(text)...)
	textPadded := uc.pkcs7pad(signedText, aes.BlockSize)

	cbc := cipher.NewCBCEncrypter(block, ivBin)
	cbc.CryptBlocks(textPadded, textPadded)

	ivB64 := bytes.TrimRight(bytes.Replace(bytes.Replace(uc.encodeBase64(ivBin), []byte{'/'}, []byte{'_'}, -1), []byte{'+'}, []byte{'-'}, -1), "=")
	textB64 := bytes.TrimRight(bytes.Replace(bytes.Replace(uc.encodeBase64(textPadded), []byte{'/'}, []byte{'_'}, -1), []byte{'+'}, []byte{'-'}, -1), "=")

	return append(append(ivB64, '/'), textB64...), nil
}

func (uc *cypher) decryptB64(encText string) ([]byte, error) {
	encodedStr := strings.Split(encText, "/")
	if len(encodedStr) < 2 {
		return nil, errors.New("Bad string scheme")
	}

	padder := "===="
	ivPad := ""
	if len(encodedStr[0])%4 > 0 {
		ivPad = padder[len(encodedStr[0])%4:]
	}
	textPad := ""
	if len(encodedStr[1])%4 > 0 {
		textPad = padder[len(encodedStr[1])%4:]
	}

	ivB64 := strings.Replace(strings.Replace(encodedStr[0], "_", "/", -1), "-", "+", -1) + ivPad
	textB64 := strings.Replace(strings.Replace(encodedStr[1], "_", "/", -1), "-", "+", -1) + textPad

	clog.Info("Crypt", "Decrypt_b64", "IV_B64: %s -- TXT_B64: %s", ivB64, textB64)

	keyBin := make([]byte, hex.DecodedLen(len(uc.HexKey)))
	hex.Decode(keyBin, uc.HexKey)
	ivBin := uc.decodeBase64(ivB64)
	textBin := uc.decodeBase64(textB64)
	block, err := aes.NewCipher(keyBin)
	if err != nil {
		return nil, err
	}

	// iv_bin = iv_bin[:aes.BlockSize]
	// textBin = textBin[aes.BlockSize:]
	if len(textBin)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	cbc := cipher.NewCBCDecrypter(block, ivBin)
	cbc.CryptBlocks(textBin, textBin)
	textUnpadded := uc.pkcs7unpad(textBin, aes.BlockSize)
	textUnpadded = textUnpadded[uc.HashSize:]

	return textUnpadded, nil
}
