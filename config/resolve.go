package config

func ResolveConfigFile(configFile string) (string, error) {
	if configFile != "" {
		return configFile, nil
	}

	patterns := []string{
		"backscribe.yaml", "backscribe.yml", "backscribe.json",
		".backscribe.yaml", ".backscribe.yml", ".backscribe.json",
	}
	return findFirstMatchingFile(patterns)
}
