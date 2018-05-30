package lib

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SmsSender 发送短信接口
type SmsSender interface {
	Send(mobile string, content string, template string) error
}

// NewAliSms Ali短信接口实现
func NewAliSms() SmsSender {
	return &aliSms{
		//下面三个参数需要替换
		accessKeyID:  "demokey",
		accessSecret: "demosecret",
		signName:     "demo",
	}
}

type aliSms struct {
	accessKeyID  string
	accessSecret string
	signName     string
}

// Send https://blog.wolfogre.com/posts/send-sms-by-aliyun-in-golang/
func (s *aliSms) Send(mobile string, content string, template string) error {
	timezone, err := time.LoadLocation("GMT0") // 这里一定要设置GMT时区
	if err != nil {
		panic(err)
	}

	paras := make(map[string]string)

	// 1. 系统参数
	paras["SignatureMethod"] = "HMAC-SHA1"
	paras["SignatureNonce"] = fmt.Sprintf("%s", uuid.New().String()) // 原例子中是使用 UUID，但 golang 原生包里并没有支持，故用随机字符串代替
	paras["AccessKeyId"] = s.accessKeyID
	paras["SignatureVersion"] = "1.0"
	paras["Timestamp"] = time.Now().In(timezone).Format("2006-01-02T15:04:05Z")
	paras["Format"] = "XML"

	// 2. 业务参数
	paras["Action"] = "SendSms"
	paras["Version"] = "2017-05-25"
	paras["RegionId"] = "cn-hangzhou"
	paras["PhoneNumbers"] = mobile
	paras["SignName"] = s.signName   //签名名称， 例如：云数创
	paras["TemplateCode"] = template //签名的模板编码， 例如： "SMS_71390007"
	paras["TemplateParam"] = content //短信模板替换Json， 例如： `{"code":"123"}`
	paras["OutId"] = "123"

	// 3. 去除签名关键字Key
	delete(paras, "Signature")

	// 4. 参数KEY排序
	parasIndex := make([]string, 0)
	for k := range paras {
		parasIndex = append(parasIndex, k)
	}
	sort.Strings(parasIndex)

	// 5. 构造待签名的字符串
	sortedQueryString := ""
	for _, v := range parasIndex {
		sortedQueryString = sortedQueryString + "&" + specialUrlEncode(v) + "=" + specialUrlEncode(paras[v])
	}

	// 去除第一个多余的&符号
	sortedQueryString = sortedQueryString[1:]

	stringToSign := "GET" + "&" + specialUrlEncode("/") + "&" + specialUrlEncode(sortedQueryString)

	signStr := sign(s.accessSecret+"&", stringToSign)

	// 6. 签名最后也要做特殊URL编码
	signature := specialUrlEncode(signStr)

	// 最终打印出合法GET请求的URL
	urlStr := "http://dysmsapi.aliyuncs.com/?Signature=" + signature + "&" + sortedQueryString

	resp, err := http.Get(urlStr)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	var sr smsResponse
	err = xml.Unmarshal(body, &sr)

	if err != nil {
		var er smsError
		err1 := xml.Unmarshal(body, &er)
		if err1 == nil {
			return errors.New(er.Message)
		}
		return err
	}

	if sr.Code != "OK" {
		return errors.New(sr.Message)
	}
	/*
		<?xml version='1.0' encoding='UTF-8'?><SendSmsResponse><Message>OK</Message><RequestId>0EF9C684-209A-4452-9593-82BEC145764F</RequestId><BizId>736704824194315045^0</BizId><Code>OK</Code></SendSmsResponse>
	*/

	return nil
}

type smsResponse struct {
	XMLName   xml.Name `xml:"SendSmsResponse"`
	Message   string   `xml:"Message"`
	RequestID string   `xml:"RequestId"`
	BizID     string   `xml:"BizId"`
	Code      string   `xml:"Code"`
}

/*
<?xml version='1.0' encoding='UTF-8'?><Error><RequestId>162837D0-EA67-43B2-8347-C21780154880</RequestId><HostId>dysmsapi.aliyuncs.com</HostId><Code>InvalidAccessKeyId.NotFound</Code><Message>Specified access key is not found.</Message><Recommend><![CDATA[https://error-center.aliyun.com/status/search?Keyword=InvalidAccessKeyId.NotFound&source=PopGw]]></Recommend></Error>
*/
type smsError struct {
	XMLName   xml.Name `xml:"Error"`
	Message   string   `xml:"Message"`
	RequestID string   `xml:"RequestId"`
	Recommend string   `xml:"Recommend"`
}

func randomString(length int) string {
	base := "abcdefghijklmnopqrstuvwxyz1234567890"
	result := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(result) < length {
		index := r.Intn(len(base))
		result = result + base[index:index+1]
	}
	return result
}

func specialUrlEncode(value string) string {
	result := url.QueryEscape(value)
	result = strings.Replace(result, "+", "%20", -1)
	result = strings.Replace(result, "*", "%2A", -1)
	result = strings.Replace(result, "%7E", "~", -1)
	return result
}

func sign(accessSecret, strToSign string) string {
	mac := hmac.New(sha1.New, []byte(accessSecret))
	mac.Write([]byte(strToSign))
	signData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signData)
}
