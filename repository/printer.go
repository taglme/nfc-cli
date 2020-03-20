package repository

import (
	"fmt"
	apiModels "github.com/taglme/nfc-goclient/pkg/models"
)

func (s *RepositoryService) printAppInfo(info apiModels.AppInfo) {
	fmt.Println("Server:")
	if len(info.Version) > 0 {
		fmt.Printf("   Version: %s\n", info.Version)
	}
	if len(info.Commit) > 0 {
		fmt.Printf("   Commit: %s\n", info.Commit)
	}
	if len(info.SDKInfo) > 0 {
		fmt.Printf("   SDK: %s\n", info.SDKInfo)
	}
	if len(info.Platform) > 0 {
		fmt.Printf("   Platform: %s\n", info.Platform)
	}
	if len(info.BuildTime) > 0 {
		fmt.Printf("   Build time: %s\n", info.BuildTime)
	}
}

func (s *RepositoryService) printAdapters(adapters []apiModels.Adapter) {
	if len(adapters) == 0 {
		fmt.Println("Adapters not found")
		return
	}

	fmt.Println("Adapters:")

	for i, a := range adapters {
		fmt.Printf("[%d] %s\n", i+1, a.Name)
	}

	fmt.Println()
}
