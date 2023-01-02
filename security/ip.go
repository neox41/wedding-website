package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)
var cacheIP *cache.Cache
func init(){
	cacheIP = cache.New(60*time.Minute, 60*time.Minute)
}
// IPInfo describes a particular IP address.
type IPInfo struct {
	// IP holds the described IP address.
	IP string
	// Hostname holds a DNS name associated with the IP address.
	Hostname string
	// City holds the city of the ISP location.
	City string
	// Country holds the two-letter country code.
	Country string
	// Loc holds the latitude and longitude of the
	// ISP location as a comma-separated northing, easting
	// pair of floating point numbers.
	Loc string
	// Org describes the organization that is
	// responsible for the IP address.
	Org string
	// Postal holds the post code or zip code region of the ISP location.
	Postal string
}
// MyIP provides information about the public IP address of the client.
func MyIP() (*IPInfo, error) {
	return ForeignIP("")
}

// ForeignIP provides information about the given IP address,
// which should be in dotted-quad form.
func ForeignIP(ip string) (*IPInfo, error) {
	if ipParsed := net.ParseIP(ip); ipParsed == nil {
		return nil, errors.New("Invalid IP address")
	}
	response, err := http.Get(fmt.Sprintf("http://ipinfo.io/%s/json", ip))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var ipinfo IPInfo
	if err := json.Unmarshal(contents, &ipinfo); err != nil {
		return nil, err
	}
	return &ipinfo, nil
}
func IsForItalian(remoteAddr string)bool{
	var country string
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil{
		return true
	}
	countryCached, found := cacheIP.Get(ip)
	if !found {
		ipinfo, err := ForeignIP(ip)
		if err != nil{
			return true
		}
		country = ipinfo.Country
		cacheIP.Set(ip, country, cache.DefaultExpiration)
	}else{
		country = countryCached.(string)
	}

	if country == "IT" || country == "DE"{
		return true
	}

	return false
}
