package response

/**
 * @Author: 南宫乘风
 * @Description:
 * @File:  responsebody.go
 * @Email: 1794748404@qq.com
 * @Date: 2024-05-17 16:24
 */

type ResponseBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
