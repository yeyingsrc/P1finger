package detectbyFofa

import (
	"P1finger/modules/detectbyRule"
	"reflect"
	"testing"
)

func TestSplitDomainsAndIPs(t *testing.T) {
	inputs := []string{
		"https://ne.zzuli.edu.cn",
		"https://jgxy.zzuli.edu.cn",
		"www.zzuli.edu.cn",
		"https://zcgs.zzuli.edu.cn",
		"136.243.12.69:8083",
		"136.243.12.67",
		"https://jiaowu.zzuli.edu.cn",
		"jxpg.zzuli.edu.cn",
		"https://kys.zzuli.edu.cn",
		"https://tw.zzuli.edu.cn",
		"http://jfzx-xmsb.zzuli.edu.cn",
		"136.243.12.69:9000",
	}

	expectedDomains := []string{
		"https://ne.zzuli.edu.cn",
		"https://jgxy.zzuli.edu.cn",
		"www.zzuli.edu.cn",
		"https://zcgs.zzuli.edu.cn",
		"https://jiaowu.zzuli.edu.cn",
		"jxpg.zzuli.edu.cn",
		"https://kys.zzuli.edu.cn",
		"https://tw.zzuli.edu.cn",
		"http://jfzx-xmsb.zzuli.edu.cn",
	}

	expectedIPs := []string{
		"136.243.12.69:8083",
		"136.243.12.67",
		"136.243.12.69:9000",
	}

	domains, ips := detectbyRule.splitDomainsAndIPs(inputs)

	if !reflect.DeepEqual(domains, expectedDomains) {
		t.Errorf("Expected domains: %v, got: %v", expectedDomains, domains)
	}

	if !reflect.DeepEqual(ips, expectedIPs) {
		t.Errorf("Expected IPs: %v, got: %v", expectedIPs, ips)
	}
}
