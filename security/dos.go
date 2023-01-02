package security

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"net"
	"time"
)
var (
	attendance *cache.Cache
	confirm *cache.Cache
)
func init(){
	attendance = cache.New(5*time.Minute, 10*time.Minute)
	confirm = cache.New(1*time.Hour, 1*time.Hour)
}

func Attendance(remoteAddr string) (bool, string){
	ipaddress, _, err := net.SplitHostPort(remoteAddr)
	if err != nil{
		return false, "Invalid IP Address."
	}

	if _, found := attendance.Get(ipaddress); found {
		LogToFile(fmt.Sprintf("IP %s already submitted the attendanc", ipaddress))
		return false, "You have already submitted your attendance, try again in the next few minutes."
	}
	attendance.Set(ipaddress, true, cache.DefaultExpiration)

	return true, ""
}

func Link(remoteAddr string) (bool, string){
	const(
		limit = 60
	)
	var(
		c = 1
	)
	ipaddress, _, err := net.SplitHostPort(remoteAddr)
	if err != nil{
		return false, "Invalid IP Address."
	}

	counter, found := attendance.Get(ipaddress)
	if found {
		c += counter.(int)
		if err = attendance.Replace(ipaddress, c, cache.DefaultExpiration); err != nil {
			return true, ""
		}
		if c > limit{
			LogToFile(fmt.Sprintf("IP %s triggered too many requests", ipaddress))
			return false, "Too many requests, try again in the next few minutes."
		}
	}else{
		attendance.Set(ipaddress, c, cache.DefaultExpiration)
	}
	return true, ""
}