package mockbun

import (
	"net/http"
	"slices"

	"github.com/kyswtn/terraform-provider-porkbun/internal/client"
)

func (s *Server) dns_delete(w http.ResponseWriter, r *http.Request) {
	domain := r.PathValue("domain")
	id := r.PathValue("id")

	if records, ok := s.dnsRecords[domain]; ok {
		recordIndex := slices.IndexFunc(records, func(record client.DNSRecord) bool {
			return record.ID == id
		})

		if recordIndex != -1 {
			records = append(records[:recordIndex], records[recordIndex+1:]...)
			s.dnsRecords[domain] = records

			s.Write(w, client.Status{
				Value: "SUCCESS",
			})
			return
		}
	}

	s.Write(w, client.Status{
		Value:   "FAILURE",
		Message: "Record not found",
	})
}
