package transformer

import (
	"encoding/json"
	"git-pd.megvii-inc.com/liuwei02/Edward/taurusrpc"
	"log"
)

func Transform(slice []string) []interface{} {
	res := make([]interface{}, len(slice))
	for _, content := range slice {
		/**
		value, ok := content.(string)
		if !ok {
			continue
		}
		*/
		var info taurusrpc.SearchXIDInfo
		err := json.Unmarshal([]byte(content), &info)
		if err != nil {
			log.Fatal(err)
			continue
		}
		res = append(res, info)
	}
	return res
}
