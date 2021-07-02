package common

import (
	"crypto/md5"
	"fmt"
	"github.com/wuzehv/passport/util/static"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func GenPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("generate password error: %v\n", err)
	}

	return string(hash)
}

func VerifyPassword(hash string, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

func GenToken() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d", rand.Int()))))
}

type key []string

func (k key) Len() int {
	return len(k)
}

func (k key) Less(i, j int) bool {
	return k[j] > k[i]
}

func (k key) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func GenSign(m map[string]string, secret string) string {
	keyArr := make(key, len(m))
	i := 0
	for k, _ := range m {
		if k == static.Sign {
			continue
		}

		keyArr[i] = k
		i++
	}
	sort.Sort(keyArr)

	var b strings.Builder
	for _, v := range keyArr {
		b.WriteString(m[v])
	}
	b.WriteString(secret)

	return fmt.Sprintf("%X", md5.Sum([]byte(b.String())))
}
