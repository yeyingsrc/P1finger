package detectbyFofa

import (
	"encoding/json"
)

// API响应文档
//{
//  "error": false, // 是否出现错误
//  "email": "fo****t@baimaohui.net", // 邮箱地址：
//  "username": "fofabot", // 用户名
//  "category": "user", // 用户种类
//  "fcoin": 0, // F币
//  "fofa_point": 49200, // F点
//  "remain_free_point": 0, // 剩余免费F点
//  "remain_api_query": 49992, // API月度剩余查询次数
//  "remain_api_data": 499398, // API月度剩余返回数量
//  "isvip": true, // 是否是会员
//  "vip_level": 12, // page.api.whether.level
//  "is_verified": false,
//  "avatar": "https://nosec.org/missing.jpg",
//  "message": "",
//  "fofacli_ver": "4.0.3",
//  "fofa_server": true
//}

type AccountInfo struct {
	Error           bool     `json:"error"`            // 是否出现错误
	ErrMsg          string   `json:"errmsg,omitempty"` // error string message
	Email           string   `json:"email"`
	Username        string   `json:"username"`
	Category        string   `json:"category"`          // 用户种类
	FCoin           int      `json:"fcoin"`             // F币
	FofaPoint       int64    `json:"fofa_point"`        // F点
	RemainFreePoint int      `json:"remain_free_point"` // 剩余免费F点
	RemainApiQuery  int      `json:"remain_api_query"`  // API月度剩余查询次数
	RemainApiData   int      `json:"remain_api_data"`   // API月度剩余返回数量
	IsVIP           bool     `json:"isvip"`             // 是否是会员
	VIPLevel        VipLevel `json:"vip_level"`
	IsVerified      bool     `json:"is_verified"`
	Avatar          string   `json:"avatar"` // page.api.whether.level
	Message         string   `json:"message"`
	FofacliVer      string   `json:"fofacli_ver"`
	FofaServer      bool     `json:"fofa_server"`
}

// DeductMode should deduct fcoin automatically or just use free limit
type DeductMode int

const (
	// DeductModeFree only use free limit size
	DeductModeFree DeductMode = 0
	// DeductModeFCoin deduct fcoin automatically if account has fcoin
	DeductModeFCoin DeductMode = 1
)

// ParseDeductMode parse string to DeductMode
func ParseDeductMode(v string) DeductMode {
	switch v {
	case "0", "free":
		return DeductModeFree
	case "1", "fcoin":
		return DeductModeFCoin
	default:
		panic("unknown deduct mode")
	}
}

type VipLevel int

const (
	VipLevelNone        VipLevel = 0 // 注册用户
	VipLevelNormal      VipLevel = 1 // 普通会员
	VipLevelAdvanced    VipLevel = 2 // 高级会员
	VipLevelEnterprise  VipLevel = 3 // 企业版
	VipLevelEnterprise2 VipLevel = 5 // 企业版
)

const (
	VipLevelSubPersonal VipLevel = 11 // 订阅个人
	VipLevelSubPro      VipLevel = 12 // 订阅专业
	VipLevelSubBuss     VipLevel = 13 // 订阅商业版
)

const (
	VipLevelRed     VipLevel = 20  // 红队版
	VipLevelStudent VipLevel = 22  // 教育账户
	VipLevelNever   VipLevel = 100 // 不可能的等级
)

func (ai AccountInfo) String() string {
	d, _ := json.MarshalIndent(ai, "", "  ")
	return string(d)
}

// AccountInfo fetch account info from fofa
func (c *FofaClient) AccountInfo() (ac AccountInfo, err error) {
	err = c.Fetch("info/my", nil, &ac)
	return
}

// freeSize 当前用户可以免费使用的数据量
func (c *FofaClient) freeSize() int {
	if !c.Account.IsVIP {
		return 0
	}

	switch c.Account.VIPLevel {
	case VipLevelNormal:
		return 100
	case VipLevelAdvanced:
		return 10000
	case VipLevelEnterprise, VipLevelEnterprise2:
		return 100000
	case VipLevelRed:
		return 10000
	case VipLevelStudent:
		return 10000
	// 订阅用户：通过 api 查询余额
	case VipLevelSubPersonal:
		fallthrough
	case VipLevelSubPro:
		fallthrough
	case VipLevelSubBuss:
		info, err := c.AccountInfo()
		if err != nil {
			info = c.Account
		}
		if info.RemainApiQuery > 0 {
			return info.RemainApiData
		}
	}
	// other level, ignore free limit check
	return -1
}
