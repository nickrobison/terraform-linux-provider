package firewall

import (
	"net/http"

	"github.com/nickrobison/terraform-linux-provider/common"
	"github.com/rs/zerolog"
)

func HandleZoneList(client FirewallClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)
		objects, err := client.ListZones(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Cannot list firewall zones")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		zones := make([]common.FirewallZoneResponse, len(objects))
		for i, v := range objects {
			name, err := v.Name()
			if err != nil {
				log.Error().Err(err).Msgf("Cannot get name for zone %s", v.obj.Path())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			description, err := v.Description()
			if err != nil {
				description = ""
			}
			target, err := v.Target()
			if err != nil {
				target = ""
			}
			services, err := v.Services()
			if err != nil {
				services = []string{}
			}
			ports, err := v.Ports()
			if err != nil {
				ports = []string{}
			}
			richRules, err := v.RichRules()
			if err != nil {
				richRules = []string{}
			}

			zones[i] = common.FirewallZoneResponse{
				Name:        name,
				Description: description,
				Target:      target,
				Services:    services,
				Ports:       ports,
				RichRules:   richRules,
			}
		}
		common.Encode(w, r, 200, zones)
	})
}

func HandleZoneGet(client FirewallClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)
		zoneName := r.PathValue("name")
		if zoneName == "" {
			http.Error(w, "Zone name is required", http.StatusBadRequest)
			return
		}

		zoneObj, err := client.GetZone(ctx, zoneName)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot get zone %s", zoneName)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		name, err := zoneObj.Name()
		if err != nil {
			log.Error().Err(err).Msgf("Cannot get name for zone %s", zoneObj.obj.Path())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		description, err := zoneObj.Description()
		if err != nil {
			description = ""
		}
		target, err := zoneObj.Target()
		if err != nil {
			target = ""
		}
		services, err := zoneObj.Services()
		if err != nil {
			services = []string{}
		}
		ports, err := zoneObj.Ports()
		if err != nil {
			ports = []string{}
		}
		richRules, err := zoneObj.RichRules()
		if err != nil {
			richRules = []string{}
		}

		zone := common.FirewallZoneResponse{
			Name:        name,
			Description: description,
			Target:      target,
			Services:    services,
			Ports:       ports,
			RichRules:   richRules,
		}
		common.Encode(w, r, 200, zone)
	})
}

func HandleZoneCreate(client FirewallClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)

		var req common.FirewallZoneCreateRequest
		err := common.DecodeRequest(r, &req)
		if err != nil {
			log.Error().Err(err).Msg("Cannot decode request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		settings := ZoneSettings{
			Description: req.Description,
			Target:      req.Target,
		}

		err = client.AddZone(ctx, req.Name, settings)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot create zone %s", req.Name)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		zone := common.FirewallZoneResponse{
			Name:        req.Name,
			Description: req.Description,
			Target:      req.Target,
			Services:    []string{},
			Ports:       []string{},
			RichRules:   []string{},
		}
		common.Encode(w, r, http.StatusCreated, zone)
	})
}

func HandleZoneDelete(client FirewallClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)
		zoneName := r.PathValue("name")
		if zoneName == "" {
			http.Error(w, "Zone name is required", http.StatusBadRequest)
			return
		}

		err := client.RemoveZone(ctx, zoneName)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot delete zone %s", zoneName)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
