package security

import (
	"git.sda1.net/media-proxy-go/core"
	"net"
	"net/url"
	"strings"
)

func isPrivateAddress(address string) bool {
	// 0.0.0.0は127.0.0.1と扱うOSが多いのでプライベートアドレスとして扱う
	if address == "0.0.0.0" {
		return true
	}

	ip := net.ParseIP(address)
	return ip != nil && (ip.IsLoopback() || ip.IsPrivate())
}

func IsSafeUrl(requestedUrl string) bool {
	// URLをパース
	parsedURL, err := url.Parse(requestedUrl)
	if err != nil {
		core.MsgErr("Failed to parse URL")
		return false
	}

	// Unixソケットを拒否
	if strings.HasPrefix(parsedURL.Scheme, "unix") {
		return false
	}

	// ポート番号を取得
	port := parsedURL.Port()

	if port != "" && port != "80" && port != "443" {
		// 宛先が80と443以外ならブロック
		return false
	}

	// hostname検証
	hostname := parsedURL.Hostname()

	ip := net.ParseIP(hostname)
	if ip != nil {
		//IPアドレスが指定されている場合、それがプライベートアドレスならブロック
		if isPrivateAddress(ip.String()) {
			return false
		}
	} else {
		// ドメイン名を名前解決してIPアドレスを取得
		ips, err := net.LookupIP(hostname)
		if err != nil {
			// 失敗したらとりあえず拒否
			return false
		}

		// 取得したIPアドレスにプライベートアドレスが含まれないか確認
		for _, ip := range ips {
			core.MsgDebug("IP: " + ip.String())
			if isPrivateAddress(ip.String()) {
				core.MsgDebug("Private address detected")
				return false
			}
		}
	}

	return true
}
