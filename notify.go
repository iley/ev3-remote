package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

func NotifyAgent(port int, agentUrl, secret string) error {
	ip, err := MyIp()
	if err != nil {
		return err
	}
	log.Printf("My IP is %s", ip.String())
	var myUrl string
	if port == 80 {
		myUrl = "http://" + ip.String()
	} else {
		myUrl = fmt.Sprintf("http://%s:%d", ip.String(), port)
	}
	log.Printf("My url is %s", myUrl)
	timestamp, err := GetAgentTimestamp(agentUrl)
	if err != nil {
		return err
	}
	signature := Sign(fmt.Sprintf("%s:%d", myUrl, timestamp), secret)
	err = SendAgentUpdate(agentUrl, myUrl, timestamp, signature)
	if err != nil {
		return err
	}
	return nil
}

func MyIp() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			if Verbose() {
				log.Printf("Error parsing addr %s, skipping. %s",
					addr, err.Error())
			}
			continue
		}
		if ip.IsLoopback() {
			if Verbose() {
				log.Printf("Skipping loopback address %s", addr)
			}
			continue
		}
		if ip.To4() == nil {
			if Verbose() {
				log.Printf("Skipping non-IPv4 address %s", addr)
			}
			continue
		}
		return ip, nil
	}
	return nil, errors.New("Couldn't find IP")
}

func GetAgentTimestamp(agentUrl string) (int64, error) {
	resp, err := http.Get(agentUrl + "/time")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	timestamp, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return 0, err
	}
	return timestamp, nil
}

func Sign(message string, secret string) string {
	h := sha1.New()
	str := fmt.Sprintf("%s:%s", message, secret)
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SendAgentUpdate(agentUrl, addr string, timestamp int64, signature string) error {
	time := fmt.Sprintf("%d", timestamp)
	resp, err := http.PostForm(agentUrl+"/update",
		url.Values{"url": {addr}, "time": {time}, "signature": {signature}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("Bad response from agent: %s, %s", resp.Status, body))
	}
	return nil
}
