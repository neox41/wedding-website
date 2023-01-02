package security

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"net"
	"time"
)

var (
	login *cache.Cache
)
func init(){
	login = cache.New(10*time.Minute, 10*time.Minute)
}

type loginAttempts struct {
	Count int
	Passwords []string
}
func Login(remoteAddress, username, password string) (bool, string){
	ipaddress, _, err := net.SplitHostPort(remoteAddress)
	if err != nil{
		return false, ""
	}
	const(
		limit = 10
	)

	loginAttempt, found := attendance.Get(username)
	if found {
		l := loginAttempt.(loginAttempts)
		for _, p := range l.Passwords{
			if password == p{
				// Submitted the same password, so no brute force
				return true, ""
			}
		}

		// This is a new password and new attempt
		l.Passwords = append(l.Passwords, password)
		l.Count += 1
		if l.Count > limit{
			LogToFile(fmt.Sprintf("IP %s triggered too many login for username %s", ipaddress, username))
			return false, "Too many login failed."
		}
		if err := attendance.Replace(username, l, cache.DefaultExpiration); err != nil {
			return true, ""
		}
	}else{
		newCreds := loginAttempts{
			Count: 1,
			Passwords: []string{password},
		}
		attendance.Set(username, newCreds, cache.DefaultExpiration)
	}
	return true, ""
}
