package crypto

import (
    "crypto/hmac"
    "crypto/md5"

)

func Hmac(key, data []byte) []byte {
    hmac := hmac.New(md5.New, key)
    hmac.Write(data)
    return hmac.Sum(nil)
}