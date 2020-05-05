package list

import (
	"fmt"
	"github.com/yuanmomo/go-cidrman/cidr"
	"strings"
)

var (
	typeRegistry = make(map[string] FetchType)
)


type FetchType interface {
	Name() string
	FetchAndMerge() []string
}

func RegisterCommand(fetchType FetchType) string {
	entry := strings.ToLower(fetchType.Name())
	if entry == "" {
		return "empty command name"
	}
	typeRegistry[entry] = fetchType;
	return ""
}

func GetType(name string) FetchType {
	cmd, found := typeRegistry[name]
	if !found {
		return nil
	}
	return cmd
}

func Merge(fetchType FetchType, originCidrArray []string) []string {
	var cidrArray []string;

	for _, cidr := range originCidrArray {
		// trim space
		cidr = strings.TrimSpace(cidr);

		// filter ipv6
		if strings.Index(cidr,":") >= 0{
			continue
		}

		cidrArray = append(cidrArray,cidr)
	}

	res, err := cidr.MergeCIDRs(cidrArray)
	if err != nil {
		fmt.Printf("[%s] merge error : %s\n", fetchType.Name(), err.Error())
		return []string{};
	}
	return res;
}


func GetAll() map[string][]string{
	var cidrMap = map[string][]string{}

	for _, fetchType := range typeRegistry {
		merged := fetchType.FetchAndMerge();
		cidrMap[fetchType.Name()] = merged;
	}
	return  cidrMap;
}