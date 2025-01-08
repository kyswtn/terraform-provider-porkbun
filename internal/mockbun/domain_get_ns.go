package mockbun

import (
	"net/http"

	"github.com/kyswtn/terraform-provider-porkbun/internal/client"
)

// domain_get_ns is a handler for the /domain/getNs/{domain}
func (s *Server) domain_get_ns(w http.ResponseWriter, r *http.Request) {
	domain := r.PathValue("domain")

	if nameservers, ok := s.nameservers[domain]; ok {
		s.Write(w, client.GetNameserversResponse{
			Status: client.Status{
				Value: "SUCCESS",
			},
			Ns: nameservers,
		})
	} else {
		s.Write(w, client.Status{
			Value:   "FAILURE",
			Message: "Domain not found",
		})
	}
}
