package ruleClient

import (
	"embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

//go:embed P1fingersYaml/*
var YamlFingerFs embed.FS

func (r *RuleClient) LoadFingersFromFile(files []string) error {

	for _, file := range files {
		fileInf, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("File %s not existing:", fileInf.Name())
		}

		if filepath.Ext(file) == ".yaml" {
			fileBytes, err := os.ReadFile(file)
			if err != nil {
				fmt.Println("❌ 无法读取文件:", fileInf.Name(), err)
				continue
			}

			var newFingerprints []FingerprintsType
			err = yaml.Unmarshal(fileBytes, &newFingerprints)
			if err != nil {
				fmt.Println("❌ 解析 YAML 失败:", fileInf.Name(), err)
				continue
			}

			r.FingersTdSafe.FingerSlice = append(r.FingersTdSafe.FingerSlice, newFingerprints...)
		}
	}

	return nil
}

func (r *RuleClient) LoadFingersFromExEfs() error {

	fingerfolderPath := r.FingerFilePath

	files, err := YamlFingerFs.ReadDir(fingerfolderPath)
	if err != nil {
		fmt.Println("❌ 无法打开嵌入文件夹:", err)
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			fileBytes, err := YamlFingerFs.ReadFile(fingerfolderPath + "/" + file.Name())
			if err != nil {
				fmt.Println("❌ 无法读取文件:", file.Name(), err)
				continue
			}

			var newFingerprints []FingerprintsType
			err = yaml.Unmarshal(fileBytes, &newFingerprints)
			if err != nil {
				fmt.Println("❌ 解析 YAML 失败:", file.Name(), err)
				continue
			}

			r.FingersTdSafe.FingerSlice = append(r.FingersTdSafe.FingerSlice, newFingerprints...)
		}
	}
	return nil
}
