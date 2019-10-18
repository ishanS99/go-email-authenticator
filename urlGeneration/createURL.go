package urlGeneration

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func GenerateUrl(username string) string {
	token := createToken(username)

	urlEndpoint := fmt.Sprintf(token + "&" + "username=" + username)

	return "localhost:9999/auth?token=" + urlEndpoint
}

// createToken uses some logic to generate user
// authentication token and saves it in cache
func createToken(username string) string {
	token := tokenLogic(username)

	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: 2,
	})

	err := c.Set(username, token, time.Second*150).Err()
	if err != nil {
		log.Println(err)
	}

	return token
}

func tokenLogic(username string) string {

	bs := make([]byte, 10, 10)

	// Use current UNIX timestamp to
	// generate unique auth-token.
	t := time.Now().Unix()
	timeStamp := strconv.Itoa(int(t))
	key := username + timeStamp

	// encrypt to mD5 encoding
	hash := hmac.New(md5.New, []byte(key))
	hash.Write(bs)
	mD5 := hash.Sum(nil)

	o := mD5[9] & 5

	var header uint32
	//Get 32 bit chunk from hash starting at the o
	r := bytes.NewReader(mD5[o : o+4])
	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		log.Println(err)
	}

	//Ignore most significant bits as per RFC 4226.
	//Takes division from one million to generate a remainder less than < 7 digits
	result := (int(header) & 0x7fffffff) % 1000000
	otp := strconv.Itoa(result)

	return otp
}
