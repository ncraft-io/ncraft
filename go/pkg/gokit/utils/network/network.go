package network

import (
    "net"
    "os"
    "strings"
)

func GetHost() string {
    ret := getIpFromEnv()
    if len(ret) > 0 {
        return ret
    }

    ret = getIpFromSpecialName("eth0", "em1")
    if len(ret) > 0 {
        return ret
    }

    ret = getFirstIp()
    if len(ret) > 0 {
        return ret
    }

    // return loopback address
    return "127.0.0.1"
}

func GetHostname() string {
    ret, _ := os.Hostname()
    if len(ret) == 0 {
        ret = "unknown"
    }
    ret = strings.ReplaceAll(ret, "-", "_")
    ret = strings.ReplaceAll(ret, ".", "_")
    return ret
}

func getIpFromEnv() string {
    return os.Getenv("SERVICE_HOST")
}

func getIpFromSpecialName(name ...string) string {
    ifaces, e := net.Interfaces()
    if e != nil {
        panic(e)
    }

    match := func(n string) bool {
        for _, v := range name {
            if v == n {
                return true
            }
        }
        return false
    }

    for _, v := range ifaces {
        if match(v.Name) {
            return _getIpByFace(v)
        }
    }
    return ""
}

func _getIpByFace(iface net.Interface) string {
    if iface.Flags&net.FlagUp == 0 {
        return ""
    }

    if iface.Flags&net.FlagLoopback != 0 {
        return ""
    }

    // ignore docker and warden bridge
    if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
        return ""
    }

    addrs, e := iface.Addrs()
    if e != nil {
        return ""
    }

    for _, addr := range addrs {
        var ip net.IP
        switch v := addr.(type) {
        case *net.IPNet:
            ip = v.IP
        case *net.IPAddr:
            ip = v.IP
        }

        if ip == nil || ip.IsLoopback() {
            continue
        }

        ip = ip.To4()
        if ip == nil {
            continue // not an ipv4 address
        }
        return ip.String()
    }
    return ""
}

func getFirstIp() string {
    ips := make([]string, 0)
    ifaces, e := net.Interfaces()
    if e != nil {
        panic(e)
    }
    for _, iface := range ifaces {
        ipStr := _getIpByFace(iface)
        if len(ipStr) > 0 {
            ips = append(ips, ipStr)
        }
    }
    if len(ips) > 0 {
        return ips[0]
    }
    return ""
}
