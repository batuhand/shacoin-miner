package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"unsafe"
)

// Struct for making an object before sending the data
type MineHash struct {
	GeneratedStr  string `json:"generated_string"`
	GeneratedHash string `json:"generated_hash"`
	SenderWallet  string `json:"sender_wallet_id"`
}

// Struct for storing user wallet
type Wallet struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

func main() {
	URI := "http://localhost:8080/api/mine/upload/"
	args := os.Args[1:] // For getting the users wallet file
	fmt.Println(args)
	wallet_file := args[0]
	fmt.Println(wallet_file)
	my_wallet := ReturnWallet(wallet_file).PublicKey

	i := 0

	// TO-DO
	// Miner ilk çalıştığında rastgele 100 tane random string ve hashlerini hesaplayıp atar, sonrasında api dönüşünde
	// sıradaki 100 random string api tarafından daha önceden hesaplanmamış olanlar gönderilir.
	for {
		if i < 10 {

			// Creates a random string
			randStr := RandStringBytesMaskImprSrcUnsafe(15)
			// Finds the hash of the string
			hash := sha256.Sum256([]byte(randStr))
			fmt.Printf("'%x'\n", hash[:])
			fmt.Println("URL:>", URI)
			minedHash := MineHash{
				GeneratedStr:  randStr,
				GeneratedHash: fmt.Sprintf("%s", hash[:]),
				SenderWallet:  my_wallet,
			}

			jsonData, _ := json.Marshal(minedHash)

			var jsonStr = []byte(jsonData)
			req, err := http.NewRequest("POST", URI, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Accept-Charset", "utf-8")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))
			i = i + 1
		} else {
			break
		}

	}

}

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&\\'()*+,-./:;<=>?@[\\]^_`{|}~ \t\n\r\x0b\x0c"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func ReturnWallet(wallet_file string) Wallet {
	data, err := ioutil.ReadFile(wallet_file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
	wallet := Wallet{}
	json.Unmarshal(data, &wallet)
	return wallet
}
