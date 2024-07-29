package zfs

import (
	"net/http"

	"github.com/nickrobison/terraform-linux-provider/common"
	"github.com/rs/zerolog"
)

func HandleZpoolList(client ZfsClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)
		objects, err := client.ListPools(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Cannot list zpools")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pools := make([]common.ZPool, len(objects))
		for i, v := range objects {
			name, err := v.Name()
			if err != nil {
				log.Error().Err(err).Msgf("Cannot get name for pool %s", v.obj.Path())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pools[i] = common.ZPool{
				Name: name,
			}
		}
		common.Encode(w, r, 200, pools)
	})
}
