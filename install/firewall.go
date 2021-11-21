package install

import (
	"fmt"
	"github.com/kuaifan/sdos/pkg/logger"
	"strings"
)

//BuildFirewall is
func BuildFirewall() {
	if Exists("/usr/sbin/ufw") {
		if FirewallConfig.Mode == "add" {
			ufwFirewallAdd()
		} else {
			ufwFirewallDel()
		}
	} else if Exists("/usr/sbin/firewalld") {
		if FirewallConfig.Mode == "add" {
			cmdFirewallAdd()
		} else {
			cmdFirewallDel()
		}
	} else if Exists("/etc/init.d/iptables") {
		if FirewallConfig.Mode == "add" {
			iptablesFirewallAdd()
		} else {
			iptablesFirewallDel()
		}
	}
}

func ufwFirewallTemplate(mode string) string {
	FirewallConfig.Ports = strings.Replace("-", ":", FirewallConfig.Ports, -1)
	if FirewallConfig.Type == "accept" {
		FirewallConfig.Type = "allow"
	} else {
		FirewallConfig.Type = "deny"
	}
	value := ""
	if FirewallConfig.Address == "" {
		if strings.Contains(FirewallConfig.Protocol, "/") {
			tcp := fmt.Sprintf("ufw {MODE} %s %s/tcp", FirewallConfig.Type, FirewallConfig.Ports)
			udp := fmt.Sprintf("ufw {MODE} %s %s/udp", FirewallConfig.Type, FirewallConfig.Ports)
			value = fmt.Sprintf("%s && %s", tcp, udp)
		} else {
			value = fmt.Sprintf("ufw {MODE} %s %s/%s", FirewallConfig.Type, FirewallConfig.Ports, FirewallConfig.Protocol)
		}
	} else {
		if strings.Contains(FirewallConfig.Protocol, "/") {
			tcp := fmt.Sprintf("ufw {MODE} %s proto tcp from %s to any port %s", FirewallConfig.Type, FirewallConfig.Address, FirewallConfig.Ports)
			udp := fmt.Sprintf("ufw {MODE} %s proto udp from %s to any port %s", FirewallConfig.Type, FirewallConfig.Address, FirewallConfig.Ports)
			value = fmt.Sprintf("%s && %s", tcp, udp)
		} else {
			value = fmt.Sprintf("ufw {MODE} %s proto %s from %s to any port %s", FirewallConfig.Type, FirewallConfig.Protocol, FirewallConfig.Address, FirewallConfig.Ports)
		}
	}
	if mode == "del" {
		value = strings.ReplaceAll(value, "{MODE}", "delete")
	} else {
		value = strings.ReplaceAll(value, " {MODE}", "")
	}
	return value
}

func ufwFirewallAdd() {
	cmd := ufwFirewallTemplate("add")
	_, s, err := RunCommand("-c", cmd)
	if err != nil {
		logger.Error(err, s)
	}
}

func ufwFirewallDel() {
	cmd := ufwFirewallTemplate("del")
	_, s, err := RunCommand("-c", cmd)
	if err != nil {
		logger.Error(err, s)
	}
}

func cmdFirewallTemplate(mode string) string {
	value := ""
	if FirewallConfig.Address == "" {
		if strings.Contains(FirewallConfig.Protocol, "/") {
			if FirewallConfig.Type == "accept" {
				tcp := fmt.Sprintf("firewall-cmd --permanent --zone=public --{MODE}-port=%s/tcp", FirewallConfig.Ports)
				udp := fmt.Sprintf("firewall-cmd --permanent --zone=public --{MODE}-port=%s/udp", FirewallConfig.Ports)
				value = fmt.Sprintf("%s && %s", tcp, udp)
			} else {
				tcp := fmt.Sprintf("firewall-cmd --permanent --{MODE}-rich-rule=\"rule family=\"ipv4\" port protocol=\"tcp\" port=\"%s\" drop\"", FirewallConfig.Ports)
				udp := fmt.Sprintf("firewall-cmd --permanent --{MODE}-rich-rule=\"rule family=\"ipv4\" port protocol=\"udp\" port=\"%s\" drop\"", FirewallConfig.Ports)
				value = fmt.Sprintf("%s && %s", tcp, udp)
			}
		} else {
			if FirewallConfig.Type == "accept" {
				value = fmt.Sprintf("firewall-cmd --permanent --zone=public --{MODE}-port=%s/%s", FirewallConfig.Ports, FirewallConfig.Protocol)
			} else {
				value = fmt.Sprintf("firewall-cmd --permanent --{MODE}-rich-rule=\"rule family=\"ipv4\" port protocol=\"%s\" port=\"%s\" drop\"", FirewallConfig.Protocol, FirewallConfig.Ports)
			}
		}
	} else {
		if strings.Contains(FirewallConfig.Protocol, "/") {
			tcp := fmt.Sprintf("firewall-cmd --permanent --{MODE}-rich-rule=\"rule family=\"ipv4\" source address=\"%s\" port protocol=\"tcp\" port=\"%s\" %s\"", FirewallConfig.Address, FirewallConfig.Ports, FirewallConfig.Type)
			udp := fmt.Sprintf("firewall-cmd --permanent --{MODE}-rich-rule=\"rule family=\"ipv4\" source address=\"%s\" port protocol=\"udp\" port=\"%s\" %s\"", FirewallConfig.Address, FirewallConfig.Ports, FirewallConfig.Type)
			value = fmt.Sprintf("%s && %s", tcp, udp)
		} else {
			value = fmt.Sprintf("firewall-cmd --permanent --{MODE}-rich-rule=\"rule family=\"ipv4\" source address=\"%s\" port protocol=\"%s\" port=\"%s\" %s\"", FirewallConfig.Address, FirewallConfig.Protocol, FirewallConfig.Ports, FirewallConfig.Type)
		}
	}
	if mode == "del" {
		value = strings.ReplaceAll(value, "{MODE}", "remove")
	} else {
		value = strings.ReplaceAll(value, "{MODE}", "add")
	}
	return value
}

func cmdFirewallAdd() {
	cmd := cmdFirewallTemplate("add")
	_, s, err := RunCommand("-c", cmd)
	if err != nil {
		logger.Error(err, s)
	}
}

func cmdFirewallDel() {
	cmd := cmdFirewallTemplate("del")
	_, s, err := RunCommand("-c", cmd)
	if err != nil {
		logger.Error(err, s)
	}
}

func iptablesFirewallTemplate(mode string) string {
	FirewallConfig.Ports = strings.Replace("-", ":", FirewallConfig.Ports, -1)
	value := ""
	if FirewallConfig.Address == "" {
		if strings.Contains(FirewallConfig.Protocol, "/") {
			tcp := fmt.Sprintf("iptables {MODE} INPUT -p tcp -m state --state NEW -m tcp --dport %s -j %s", FirewallConfig.Ports, FirewallConfig.Type)
			udp := fmt.Sprintf("iptables {MODE} INPUT -p udp -m state --state NEW -m udp --dport %s -j %s", FirewallConfig.Ports, FirewallConfig.Type)
			value = fmt.Sprintf("%s && %s", tcp, udp)
		} else {
			value = fmt.Sprintf("iptables {MODE} INPUT -p tcp -m state --state NEW -m %s --dport %s -j %s", FirewallConfig.Protocol, FirewallConfig.Ports, FirewallConfig.Type)
		}
	} else {
		if strings.Contains(FirewallConfig.Protocol, "/") {
			tcp := fmt.Sprintf("iptables {MODE} INPUT -s %s -p tcp --dport %s -j %s", FirewallConfig.Address, FirewallConfig.Ports, FirewallConfig.Type)
			udp := fmt.Sprintf("iptables {MODE} INPUT -s %s -p udp --dport %s -j %s", FirewallConfig.Address, FirewallConfig.Ports, FirewallConfig.Type)
			value = fmt.Sprintf("%s && %s", tcp, udp)
		} else {
			value = fmt.Sprintf("iptables {MODE} INPUT -s %s -p %s --dport %s -j %s", FirewallConfig.Address, FirewallConfig.Protocol, FirewallConfig.Ports, FirewallConfig.Type)
		}
	}
	if mode == "del" {
		value = strings.ReplaceAll(value, "{MODE}", "-D")
	} else {
		value = strings.ReplaceAll(value, "{MODE}", "-I")
	}
	return value
}

func iptablesFirewallAdd() {
	cmd := iptablesFirewallTemplate("add")
	_, s, err := RunCommand("-c", cmd)
	if err != nil {
		logger.Error(err, s)
	}
}

func iptablesFirewallDel() {
	cmd := iptablesFirewallTemplate("del")
	_, s, err := RunCommand("-c", cmd)
	if err != nil {
		logger.Error(err, s)
	}
}
