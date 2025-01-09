package installer

type Tool struct {
	Name         string
	Versions     []string
	CheckCommand string
	InfoCommand  string
}

func GetDefaultTools() []Tool {
	return []Tool{
		{
			Name:         "aws-cli",
			CheckCommand: "aws --version",
			InfoCommand:  "aws --version",
		},
		{
			Name:         "azure-cli",
			CheckCommand: "az --version",
			InfoCommand:  "az --version",
		},
		{
			Name:         "terraform",
			Versions:     []string{"1.5.0", "1.4.6", "1.3.9"},
			CheckCommand: "terraform --version",
			InfoCommand:  "terraform version",
		},
		{
			Name:         "packer",
			Versions:     []string{"1.9.1", "1.8.7"},
			CheckCommand: "packer --version",
			InfoCommand:  "packer --version",
		},
		{
			Name:         "vault",
			Versions:     []string{"1.13.3", "1.12.7"},
			CheckCommand: "vault --version",
			InfoCommand:  "vault version",
		},
		{
			Name:         "consul",
			Versions:     []string{"1.15.2", "1.14.7"},
			CheckCommand: "consul --version",
			InfoCommand:  "consul version",
		},
		{
			Name:         "go",
			Versions:     []string{"1.20.5", "1.19.10"},
			CheckCommand: "go version",
			InfoCommand:  "go version",
		},
		{
			Name:         "python",
			Versions:     []string{"3.11.4", "3.10.11"},
			CheckCommand: "python3 --version",
			InfoCommand:  "python3 --version",
		},
	}
}