package token

import (
    "encoding/json"
    "testing"
    "time"
    "reflect"
)

func TestTokenBase(t *testing.T) {
    key := "bajiuwenqingtian" // mod 16 == 0

    rawTokenContent := &TokenContent{
        CreateTime : uint64(time.Now().Unix()),
        ID : "sun-moon-star-star",
        Message : "天若不爱酒，酒星不在天。地若不爱酒，地应无酒泉。",
    }

    token, err := GenerateToken(rawTokenContent, []byte(key))

    if err != nil {
        t.Fatal(err.Error())
    }

    var tokenContent *TokenContent
    tokenContent, err = CheckToken(token, []byte(key))

    if err != nil {
        t.Fatal(err.Error())
    }

    if !reflect.DeepEqual(tokenContent, rawTokenContent)  {
        t.Fatal("not useful")
    }
}

// token仅做身份校验的时候，不该包含无关身份校验的信息，应当尽量精简
func TestJson(t *testing.T) {
    type Article struct {
        Title string
        Meta string
        Text string
    }

    key := "qingshanyijiuzai"

    text := "自暴者，不可与有言也；自弃者，不可与有为也。言非礼义，谓之自暴也；"
    text += "吾身不能居仁由义，谓之自弃也。仁，人之安宅也；义，人之正路也。旷安宅而弗居，舍正路而不由，哀哉！"

    article := &Article{
        Title: "《孟子·离娄章句上·第十节》",
        Meta: "孟子（约公元前372年—公元前289年），名轲，字子舆，邹国（今山东邹城东南）人。",
        Text: text,
    }

    articleStr, err := json.Marshal(article)

    if err != nil {
        t.Fatal(err.Error())
    }

    rawTokenContent := &TokenContent{
        CreateTime : uint64(time.Now().Unix()),
        ID : "sun-moon-star-star",
        Message : string(articleStr),
    }

    token, err := GenerateToken(rawTokenContent, []byte(key))

    if err != nil {
        t.Fatal(err.Error())
    }

    var tokenContent *TokenContent
    tokenContent, err = CheckToken(token, []byte(key))

    if err != nil {
        t.Fatal(err.Error())
    }
    
    if !reflect.DeepEqual(tokenContent, rawTokenContent)  {
        t.Fatal("not useful")
    }
}