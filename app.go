package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "strings"

    "github.com/Tnze/CoolQ-Golang-SDK/cqp"
)

//go:generate cqcfg -c .
// cqp: 名称: demo
// cqp: 版本: 1.0.0:0
// cqp: 作者: mango
// cqp: 简介: 自动回复~
func main() { /*此处应当留空*/ }

var (
    loginQQ = ""
)

func init() {
    cqp.AppID = "cshare.site.demo"
    cqp.PrivateMsg = onPrivateMsg
    cqp.GroupMsg = onGroupMsg
}

func onPrivateMsg(subType, msgID int32, fromQQ int64, msg string, font int32) int32 {
    cqp.SendPrivateMsg(fromQQ, msg) //复读机
    return 0
}

func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
    if loginQQ == "" {
        loginQQ = strconv.Itoa(int(cqp.GetLoginQQ()))
        cqp.AddLog(cqp.Info, "qq", loginQQ)
    }
    if isAtMe(msg, loginQQ) {
        return robotAtMsg(fromGroup, fromQQ, msg)
    }

    return 0
}

func isCQMsg(msg string) bool {
    return strings.Contains(msg, "CQ:")
}

func isAtMe(msg, myqq string) bool {
    return isCQMsg(msg) && strings.Contains(msg, "CQ:at,qq="+ myqq)
}

func robotAtMsg(fromGroup, fromQQ int64, msg string) int32 {
    robotResp := Robot(msg)
    newMsg := fmt.Sprintf("[CQ:at,qq=%d] %s", fromQQ, robotResp)
    cqp.SendGroupMsg(fromGroup, newMsg)
    return 0
}


var (
    robot_url = "http://api.qingyunke.com/api.php?key=free&appid=0&msg="
)

func Robot(msg string) interface{} {
    values := url.PathEscape(msg)
    resp, err := http.Get(robot_url + values)
    if err != nil {
        fmt.Println("err", err)
        return ""
    }
    defer resp.Body.Close()
    bytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Read error", err)
    }
    fmt.Println(string(bytes))
    respMap := make(map[string]interface{})
    _ = json.Unmarshal(bytes, &respMap)
    return respMap["content"]
}
