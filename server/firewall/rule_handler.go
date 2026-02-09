package firewall

import (
	"net/http"

	"github.com/nickrobison/terraform-linux-provider/common"
	"github.com/rs/zerolog"
)

func HandleRuleAdd(client FirewallClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)

		var req common.FirewallRuleRequest
		err := common.DecodeRequest(r, &req)
		if err != nil {
			log.Error().Err(err).Msg("Cannot decode request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch req.RuleType {
		case "rich":
			err = client.AddRichRule(ctx, req.Zone, req.Rule)
		case "port":
			err = client.AddPort(ctx, req.Zone, req.Port, req.Protocol)
		case "service":
			err = client.AddService(ctx, req.Zone, req.Service)
		default:
			http.Error(w, "Invalid rule type. Must be 'rich', 'port', or 'service'", http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Error().Err(err).Msgf("Cannot add rule to zone %s", req.Zone)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := common.FirewallRuleResponse{
			Zone:     req.Zone,
			RuleType: req.RuleType,
			Rule:     req.Rule,
			Port:     req.Port,
			Protocol: req.Protocol,
			Service:  req.Service,
		}
		common.Encode(w, r, http.StatusCreated, response)
	})
}

func HandleRuleRemove(client FirewallClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)

		var req common.FirewallRuleRequest
		err := common.DecodeRequest(r, &req)
		if err != nil {
			log.Error().Err(err).Msg("Cannot decode request")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch req.RuleType {
		case "rich":
			err = client.RemoveRichRule(ctx, req.Zone, req.Rule)
		case "port":
			err = client.RemovePort(ctx, req.Zone, req.Port, req.Protocol)
		case "service":
			err = client.RemoveService(ctx, req.Zone, req.Service)
		default:
			http.Error(w, "Invalid rule type. Must be 'rich', 'port', or 'service'", http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Error().Err(err).Msgf("Cannot remove rule from zone %s", req.Zone)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
