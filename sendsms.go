package main

import (
    "bytes"
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "sort"
    "strings"

    "github.com/bitly/go-simplejson"
)

var (
    SENDCLOUD_HOST        string = "http://sendcloud.sohu.com/smsapi/send"
    SENDCLOUD_SMS_USER    string = ""
    SENDCLOUD_SMS_KEY     string = ""
    SENDCLOUD_TEMPLAGE_ID string = ""
    SENDCLOUD_MSG_TYPE    string = "0"
)

func HTTPRequest(URI string, Data string) {
    client := new(http.Client)
    req, err := http.NewRequest("POST", URI, strings.NewReader(Data))
    if err != nil {
        fmt.Println("HTTP Request Error, ", err)
        return
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    response, err := client.Do(req)

    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)

    RequestResult, err := simplejson.NewJson(body)
    fmt.Println(RequestResult)
}

func BuildQueryString(Param map[string]string, IsQueryEscapeString bool) string {

    ParamKeys := make([]string, 0)
    for k, _ := range Param {
        ParamKeys = append(ParamKeys, k)
    }

    sort.Strings(ParamKeys)

    var Buffer bytes.Buffer
    for _, k := range ParamKeys {
        Buffer.WriteString(k)
        Buffer.WriteByte('=')
        var value string = Param[k]
        if IsQueryEscapeString {
            value = url.QueryEscape(Param[k])
        }
        Buffer.WriteString(value)
        Buffer.WriteByte('&')
    }
    return Buffer.String()[0 : len(Buffer.String())-1]

}

func SendcloudSMSSend(PhoneNumber string, Message map[string]string) {

    var Buffer bytes.Buffer
    for k, v := range Message {
        Buffer.WriteString("\"%")
        Buffer.WriteString(k)
        Buffer.WriteString("%\":\"")
        Buffer.WriteString(v)
        Buffer.WriteString("\",")
    }

    var TempString string = "{" + Buffer.String()[0:len(Buffer.String())-1] + "}"

    var ParamMap = map[string]string{
        "smsUser":    SENDCLOUD_SMS_USER,
        "templateId": SENDCLOUD_TEMPLAGE_ID,
        "msgType":    SENDCLOUD_MSG_TYPE,
        "phone":      PhoneNumber,
        "vars":       TempString,
    }

    var ParamString string = BuildQueryString(ParamMap, false)

    HashObject := md5.New()
    HashObject.Write([]byte(SENDCLOUD_SMS_KEY + "&" + ParamString + "&" + SENDCLOUD_SMS_KEY))
    HashString := hex.EncodeToString(HashObject.Sum(nil))

    ParamMap["signature"] = HashString

    ParamString = BuildQueryString(ParamMap, true)

    HTTPRequest(SENDCLOUD_HOST, ParamString)
}

func main() {
    msg := map[string]string{
        "verify_code": "123456",
        "exp_time":    "300",
    }
    SendcloudSMSSend("13910248888", msg)
}
