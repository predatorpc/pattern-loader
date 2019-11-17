/****************************************************************************************************
*
*
* Utils module
* by Michael S. Merzlyakov AFKA predator_pc@12122018
*
* version v2.2.10
*
* created at 04122018
* last edit: 17062019
*
*****************************************************************************************************/

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"math/rand"
	"net/http"
	"os"
	_ "os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/labstack/echo"
	"errors"
)

const utilsModuleName = "utils.go"

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"
const WriteToLogOnly = false
const LogFileName = "out.log"
const LogRequestFileName = "in.log"

var CURRENT_TIMESTAMP string
var CURRENT_UNIXTIME time.Time
var CURRENT_TIMESTAMP_FS string

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func BToKb(b uint64) uint64 {
	return b / 1024
}

func WriteLog(FileName, Header, Module string, Message interface{}) {
	item := CURRENT_TIMESTAMP + " [ " + Header + " ] " + fmt.Sprintf("%s", Message) + ", " + Module + "\n"
	f, _ := os.OpenFile(FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	_, _ = f.WriteString(item)
	_ = f.Close()
}

func WriteCustomLog(FileName, Header string, Message interface{}) {
	item := Header + fmt.Sprintf("%s", Message)
	f, _ := os.OpenFile(FileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	_, _ = f.WriteString(item)
	_ = f.Close()
}

func PrintError(header string, message interface{}, module string) {
	sentry.CaptureException(errors.New(fmt.Sprintf("[ %s ] %s in module %s",header, message, module)))
	if WriteToLogOnly {
		WriteLog(LogFileName, header, module, message)
	} else {
		WriteLog(LogFileName, header, module, message)
		_, _ = fmt.Fprintf(color.Output, "[ %s ]", color.RedString(header))
		fmt.Println(" ", message, " - ", module)
	}
}

func PrintInfo(header string, message interface{}, module string) {
	if WriteToLogOnly {
		WriteLog(LogFileName, header, module, message)
	} else {
		WriteLog(LogFileName, header, module, message)
		_, _ = fmt.Fprintf(color.Output, "[ %s ]", color.CyanString(header))
		fmt.Println(" ", message, " - ", module)
	}
}

func PrintSuccess(header string, message interface{}, module string) {
	if WriteToLogOnly {
		WriteLog(LogFileName, header, module, message)
	} else {
		WriteLog(LogFileName, header, module, message)
		_, _ = fmt.Fprintf(color.Output, "[ %s ]", color.GreenString(header))
		fmt.Println(" ", message, " - ", module)
	}
}

func PrintDebug(header string, message interface{}, module string) {
	sentry.CaptureMessage(fmt.Sprintf("[ %s ] %s in module %s",header, message, module))
	if WriteToLogOnly {
		WriteLog(LogFileName, header, module, message)
	} else {
		WriteLog(LogFileName, header, module, message)
		_, _ = fmt.Fprintf(color.Output, "[ %s ]", color.YellowString(header))
		fmt.Println(" ", message, " - ", module)
	}
}

func LogRequest(header string, message string) {
	WriteLog(LogRequestFileName, header, "", message)
}

// URIByMap Заполняем наш мап параметрами из УРИ
func URIByMap(c echo.Context, keyMap []string) (map[string][]string, string) {
	var foreignQueryParams string
	resultMap := make(map[string][]string)
	for _, item := range keyMap {
		tmp := c.Param(strings.ToLower(item))
		if tmp == "" {
			tmp = c.QueryParam(strings.ToLower(item))
		}
		resultMap[item] = append(resultMap[item], tmp)
	}


	// fmt.Println("RESULT MAP = ", resultMap)
	// //------ Support old version of TDSs with incorrect ID and HASH representation  ------------------
	//
	// //resultMap["click_id"] = append(resultMap["click_id"], strings.Join(resultMap["click_hash"],""))
	// if strings.Join(resultMap["flow_id"], "") == "" {
	// 	resultMap["flow_id"] = append(resultMap["flow_id"], strings.Join(resultMap["flow_hash"], ""))
	// }
	// if strings.Join(resultMap["flow_hash"], "") == "" {
	// 	resultMap["flow_hash"] = append(resultMap["flow_hash"], strings.Join(resultMap["flow_hash"], ""))
	// }

	//------ Suppor incompatibility of naming webmaster / publisher -----------------------------------
	// it should not be in advertising link!!!!
	// if strings.Join(resultMap["publisher"], "") != "" {
	// 	resultMap["webmaster_id"] = nil
	// 	resultMap["webmaster_id"] = append(resultMap["webmaster_id"], strings.Join(resultMap["publisher"], ""))
	// 	resultMap["webmaster_id"] = append(resultMap["webmaster_id"], strings.Join(resultMap["publisher"], ""))
	// }
	// if strings.Join(resultMap["publisher_id"], "") != "" {
	// 	resultMap["webmaster_id"] = nil
	// 	resultMap["webmaster_id"] = append(resultMap["webmaster_id"], strings.Join(resultMap["publisher_id"], ""))
	// }
	//--------------------------------------------------------------------------------------------------

	// // Compatibility with other trackers
	// if strings.Join(resultMap["utm_source"], "") != "" {
	// 	resultMap["sub1"] = nil
	// 	resultMap["sub1"] = append(resultMap["sub1"], strings.Join(resultMap["utm_source"], ""))
	// }
	// if strings.Join(resultMap["utm_campaign"], "") != "" {
	// 	resultMap["sub2"] = nil
	// 	resultMap["sub2"] = append(resultMap["sub2"], strings.Join(resultMap["utm_campaign"], ""))
	// }
	// if strings.Join(resultMap["utm_medium"], "") != "" {
	// 	resultMap["sub3"] = nil
	// 	resultMap["sub3"] = append(resultMap["sub3"], strings.Join(resultMap["utm_medium"], ""))
	// }
	// if strings.Join(resultMap["utm_content"], "") != "" {
	// 	resultMap["sub4"] = nil
	// 	resultMap["sub4"] = append(resultMap["sub4"], strings.Join(resultMap["utm_content"], ""))
	// }
	// if strings.Join(resultMap["utm_term"], "") != "" {
	// 	resultMap["sub5"] = nil
	// 	resultMap["sub5"] = append(resultMap["sub5"], strings.Join(resultMap["utm_term"], ""))
	// }
	//
	// // Compatibility with other trackers
	// if strings.Join(resultMap["aff_sub1"], "") != "" {
	// 	resultMap["sub1"] = nil
	// 	resultMap["sub1"] = append(resultMap["sub1"], strings.Join(resultMap["aff_sub1"], ""))
	// }
	// if strings.Join(resultMap["aff_sub2"], "") != "" {
	// 	resultMap["sub2"] = nil
	// 	resultMap["sub2"] = append(resultMap["sub2"], strings.Join(resultMap["aff_sub2"], ""))
	// }
	// if strings.Join(resultMap["aff_sub3"], "") != "" {
	// 	resultMap["sub3"] = nil
	// 	resultMap["sub3"] = append(resultMap["sub3"], strings.Join(resultMap["aff_sub3"], ""))
	// }
	// if strings.Join(resultMap["aff_sub4"], "") != "" {
	// 	resultMap["sub4"] = nil
	// 	resultMap["sub4"] = append(resultMap["sub4"], strings.Join(resultMap["aff_sub4"], ""))
	// }
	// if strings.Join(resultMap["aff_sub5"], "") != "" {
	// 	resultMap["sub5"] = nil
	// 	resultMap["sub5"] = append(resultMap["sub5"], strings.Join(resultMap["aff_sub5"], ""))
	// }


	// forward other Major params
	keyMapMirror := keyMap
	allParams := c.QueryParams()

	// отсортируем все согласно нашей схеме
	for key, _ := range allParams {
		for _, keyMapItem := range keyMapMirror {
			if key == keyMapItem {
				delete(allParams, key)
			}
		}
	}

	// пройдемся по остаткам от сортировки
	for key, value := range allParams {
		foreignQueryParams = foreignQueryParams + "&" + key + "=" + strings.Join(value, "")
	}

	if len(foreignQueryParams) > 0 {
		return resultMap, foreignQueryParams
	} else {
		return resultMap, ""
	}

}

func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
	b, err := json.Marshal(v)

	if safeEncoding {
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func JSONPretty(Data interface{}) string {
	var out bytes.Buffer //буфер конвертации джейсона в красивый джейсон
	jsonData, _ := json.Marshal(Data)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)
	_ = json.Indent(&out, jsonData, "", "    ")
	return out.String()
}

// StringWithCharset this is very good function
func StringWithCharset(length int, charset string) string {
	var SeededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[SeededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

// returns real arguments
func Explode(str string, delimiter string) []string {
	result := strings.Split(str, delimiter)
	final := []string{}

	for i := 0; i < len(result); i++ {
		if result[i] != "" {
			final = append(final, result[i])
		}
	}
	return final
}

func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

// мы не хотим узнать сохранилась ли она или нет, пока-что
func SaveCookieToUser(value, path string) *http.Cookie {
	cookie := new(http.Cookie)
	// ставим куку на этот урл если у нас не прочиталось из запроса
	cookie.Name = "CID"
	cookie.Value = value
	cookie.Expires = time.Now().Add(365 * 24 * time.Hour) // for an year
	cookie.Path = path
	return cookie
}

func GetCertMap(path string) map[string]map[int]string {
	fds, _ := filepath.Glob(path + "/*")
	var mainCert, mainKey string
	var certMap map[string]map[int]string

	certMap = make(map[string]map[int]string, len(fds))
	for _, item := range fds {
		s, _ := filepath.Abs(item)
		certMap[s] = make(map[int]string, 2)
		certMap[s][0] = "false"
		certMap[s][1] = "false"

		certs, _ := filepath.Glob(s + "/fullchain1.pem")
		for _, certFiles := range certs {
			certFileName, _ := filepath.Abs(certFiles)
			if certFileName != "" {
				mainCert = certFileName

				if mainCert != "" {
					certMap[s][0] = mainCert
					break
				}
			}
		}

		keys, _ := filepath.Glob(s + "/privkey1.pem")
		for _, keyFiles := range keys {
			keyFileName, _ := filepath.Abs(keyFiles)
			if keyFileName != "" {
				mainKey = keyFileName

				if mainKey != "" {
					certMap[s][1] = mainKey
					break
				}
			}
		}
	}
	return certMap
}

func CheckGeo(geoSource []string, geoDestination string) bool{
	PrintDebug("GEO_INFO_SRC = ", geoSource, utilsModuleName)
	PrintDebug("GEO_INFO_DST = ", geoDestination, utilsModuleName)
	for _, item:= range geoSource {
		if item == geoDestination {
			return false
		}
	}
	return true
}

func ExtractLocale(AcceptLang string) string{
	if AcceptLang!="" {
		extract := Explode(AcceptLang, ";")
		for _, value := range extract {
			locale := Explode(value, ",")
			return locale[0];
		}
	}
	return "";
}
