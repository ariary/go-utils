package host

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

//GetExternalIP: find your extarnl ip using dig and some dns resolbver
func GetExternalIP() (ip string, err error) {
	cmd := exec.Command("dig", "@resolver4.opendns.com", "myip.opendns.com", "+short")
	ipB, err := cmd.Output()
	if err != nil {
		return "", err
	}
	ip = string(ipB)
	ip = strings.ReplaceAll(ip, "\n", "")
	return ip, err
}

//GetHostIP: return the ip of the hostname -I (if a list is returned by hostname, it only take sthe first one)
func GetHostIP() (ip string, err error) {
	cmd := exec.Command("hostname", "-I")
	ipB, err := cmd.Output()
	if err != nil {
		//retry with -i
		cmd := exec.Command("hostname", "-i")
		ipB, err := cmd.Output()
		if err != nil {
			return "", err
		}
		ip = string(ipB)
		ip = strings.Split(ip, "\n")[0]
		ip = strings.Split(ip, " ")[0]
		return ip, nil
	}

	//Only take first result
	r := bytes.NewReader(ipB)
	reader := bufio.NewReader(r)
	line, _, err := reader.ReadLine()
	ip = string(line)
	ip = strings.Split(ip, "\n")[0]
	ip = strings.Split(ip, " ")[0]
	return ip, err
}
