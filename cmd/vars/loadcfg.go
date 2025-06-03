package vars

import (
	"fmt"
	"github.com/projectdiscovery/gologger"
	"gopkg.in/yaml.v3"
	"os"
)

type P1fingerConf struct {
	CustomizeFingerFiles  []string `yaml:"CustomizeFingerFiles"`
	UseDefaultFingerFiles bool     `yaml:"UseDefaultFingerFiles"`

	FofaCredentials struct {
		Email  string `yaml:"Email"`
		ApiKey string `yaml:"ApiKey"`
	} `yaml:"FofaCredentials"`
}

func LoadAppConf(filePath string, config *P1fingerConf) error {

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		defaultConfig := P1fingerConf{
			CustomizeFingerFiles:  []string{},
			UseDefaultFingerFiles: true,
			FofaCredentials: struct {
				Email  string `yaml:"Email"`
				ApiKey string `yaml:"ApiKey"`
			}{
				Email:  "P001water@163.com",
				ApiKey: "xxxx",
			},
		}
		data, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			return fmt.Errorf("生成默认配置时出错: %v", err)
		}

		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			return fmt.Errorf("无法创建文件并写入默认配置: %v", err)
		}

		gologger.Info().Msgf("配置文件不存在，已在当前目录创建文件并写入默认配置")
		gologger.Info().Msgf("文件路径: %s", filePath)
		os.Exit(0)
	} else if err != nil {
		return fmt.Errorf("检查文件状态时出错: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("无法读取配置文件: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("解析配置文件时出错: %v", err)
	}

	return nil
}
