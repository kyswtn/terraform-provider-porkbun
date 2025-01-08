package mockbun

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kyswtn/terraform-provider-porkbun/internal/client"
	"google.golang.org/appengine/log"
)

// domain_update_ns is a handler for the /domain/updateNs/{domain}
func (s *Server) domain_update_ns(w http.ResponseWriter, r *http.Request) {
	domain := r.PathValue("domain")

	if _, ok := s.nameservers[domain]; ok {
		var nameservers client.UpdateNameserversPayload
		err := json.NewDecoder(r.Body).Decode(&nameservers)
		if err != nil {
			log.Warningf(context.TODO(), "Failed to decode request: %v", err)
		}

		s.nameservers[domain] = nameservers.Ns

		s.Write(w, client.Status{
			Value: "SUCCESS",
		})
	} else {
		s.Write(w, client.Status{
			Value:   "FAILURE",
			Message: "Domain not found",
		})
	}
}
