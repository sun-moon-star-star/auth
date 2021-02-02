package token

import (
    "fmt"
    "math/rand"
    "time"
)

func RandSeedTime() {
    rand.Seed(time.Now().Unix())
}

// 注意: 如果分布式服务同时重启有一定的概率导致id冲突, 需要注意RandSeedTime的调用时机
// 格式: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func GetUUID() (string, error) {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        return "", err
    }
    uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
        b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    return uuid, nil
}