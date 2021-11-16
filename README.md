# Envelope Rain Manager

本项目是[Group8-红包雨](https://github.com/ohroffen/Envelope)的配置服务，包含以下的功能：
- 初始化配置
- 初始化 Redis `envelope_list` 红包列表
- 于 `/reconfig` 接收热更新总金额请求，回收 `envelope_list` 中尚未放出的红包金额，并重新生成红包金额

请求参数：
```json
{
    "newBudget": 100000000
}
```

## HTTP Basic Authentication
HTTP Basic Auth 账号：`kxc3esdiu23d`

HTTP Basic Auth 密码：`xz8cvjs9q3m1`

## Curl Example:
```bash
curl -u kxc3esdiu23d:xz8cvjs9q3 -d 'newBudget=20000' -X POST 180.184.64.140:9090/reconfig
```