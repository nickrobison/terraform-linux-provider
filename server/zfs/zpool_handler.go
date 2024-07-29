package zfs

import (
	"net/http"

	"github.com/nickrobison/terraform-linux-provider/common"
)

func HandleZpoolList(client ZfsClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pools, _ := client.ListPools()
		common.Encode(w, r, 200, pools)
	})
}
