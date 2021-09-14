package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//MD5 MD5加密（32位）
func MD5(str string) string {
	sb := []byte(str)
	hash := md5.Sum(sb)
	return fmt.Sprintf("%x", hash) //将[]byte转成16进制
}

// 相差天数 （t2 - t1）
func DiffDays(t1 time.Time, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int(t2.Sub(t1).Hours() / 24)
}

//Int32ToBytes 整型转换成字节数组
func Int32ToBytes(n int32) []byte {
	buffer := bytes.NewBuffer([]byte{})
	// binary.Write(buffer, binary.BigEndian, n)
	binary.Write(buffer, binary.LittleEndian, n)
	return buffer.Bytes()
}

//BytesToInt32 字节数组转换成整型
func BytesToInt32(b []byte) int32 {
	buffer := bytes.NewBuffer(b)
	var n int32
	// binary.Read(buffer, binary.BigEndian, &n)
	binary.Read(buffer, binary.LittleEndian, &n)
	return n
}

//CheckNickName 检查昵称格式
func CheckNickName(nickname string) bool {
	return true
}

//CheckPhone 检查手机号格式
func CheckPhone(phone string) bool {
	if len(phone) == 11 {
		if match, _ := regexp.MatchString(`^1[3|4|5|8][0-9]\d{4,8}$`, phone); !match {
			return false
		} else if match, _ := regexp.MatchString(`^1[3|4|5|8][0-9]\d{4,8}$`, phone); !match {
			return false
		}
		return true
	}
	return false
}

//CheckPassword 检查密码格式
func CheckPassword(password string) bool {
	return len(password) == 32
}

//CheckCode 检查短信验证码
func CheckCode(code uint32) bool {
	return code >= 100000 && code <= 999999
}

//CheckMail 检查邮箱格式
func CheckMail(mail string) bool {
	return false
}

//CreateUID 创建6位数字的ID
func CreateUID() uint32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return uint32(r.Intn(889999) + 110000)
}

//ShowName 显示名字
func ShowName(uid uint32, username string, phone string, nickname string) string {
	if nickname != "" {
		return nickname
	}

	if username != "" {
		return username
	}

	if phone != "" {
		return phone
	}

	return strconv.Itoa(int(uid))
}

//Pow 次方算法
func Pow(x, n int) int {
	// 结果初始为0次方的值，整数0次方为1
	ret := 1
	for n != 0 {
		if n%2 != 0 {
			ret = ret * x
		}
		n /= 2
		x = x * x
	}
	return ret
}

//CheckFileIsExist 检查文件是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//PrintBuffer 打印字节流
func PrintBuffer(buffer []byte) string {
	var sa = make([]string, 0)
	for _, v := range buffer {
		sa = append(sa, fmt.Sprintf("%02X", v))
	}
	ss := strings.Join(sa, "-")
	return ss
}

// 转 int 输出
func Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}
