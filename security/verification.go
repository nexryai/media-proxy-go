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
	if ip != nil && (ip.IsLoopback() || ip.IsPrivate()) {
		return true
	}

	// プライベートアドレスの判定
	privateCIDRs := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"100.64.0.0/10",  // Tailscaleとかで使うやつ: ip.IsPrivate()では判定できない
		"169.254.0.0/16", // リンクローカルアドレス: ip.IsPrivate()では判定できない
	}

	for _, privateCIDR := range privateCIDRs {
		_, privateNet, err := net.ParseCIDR(privateCIDR)
		if err == nil && privateNet.Contains(ip) {
			return true
		}
	}

	// その他の条件が満たされない場合はパブリックアドレスとみなす
	return false
}

func IsSafeUrl(requestedUrl string) bool {
	// URLをパース
	parsedURL, err := url.Parse(requestedUrl)
	if err != nil {
		core.MsgErr("Failed to parse URL")
		return false
	}

	// https以外は拒否
	if parsedURL.Scheme != "https" {
		core.MsgDebug("Invalid protocol")
		return false
	}

	// UnixソケットとIPv6アドレス指定を拒否
	if strings.Contains(parsedURL.Hostname(), ":") {
		core.MsgDebug("Unix socket detected")
		return false
	}

	// おかしいね
	if !strings.Contains(parsedURL.Hostname(), ".") {
		core.MsgDebug("Hostname does not contain dot")
		return false
	}

	// 認証情報を含むのは拒否
	if parsedURL.User != nil {
		core.MsgDebug("User info detected")
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
