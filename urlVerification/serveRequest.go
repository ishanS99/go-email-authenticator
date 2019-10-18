package urlVerification

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/go-redis/redis"
)

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// key is the user identification string
	// used as key to store the token in redis
	key := q.Get("username")
	token := q.Get("token")

	savedToken := getSavedToken(key)

	if savedToken != token {
		_, _ = fmt.Fprint(w, "NOT ALLOWED, could not verify the Unique Token")
		return
	}

	time.Sleep(3 * time.Second)
	http.Redirect(w, r, "verified", http.StatusSeeOther)
}

// Verified is called after successful verification of user
// authentication token and http redirection to "/verified"
func Verified(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "You have been "+path.Base(r.URL.Path))
}

// getSavedToken uses a key to get its value
// from cache (in this case - Redis)
func getSavedToken(key string) string {

	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		PoolSize: 2,
	})

	savedToken, err := c.Get(key).Result()
	if err != nil {
		log.Println(err)
	}

	return savedToken
}
