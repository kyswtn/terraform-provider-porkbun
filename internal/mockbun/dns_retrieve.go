package mockbun

import (
	"net/http"
	"slices"

	"github.com/kyswtn/terraform-provider-porkbun/internal/client"
)

func (s *Server) dns_retrieve(w http.ResponseWriter, r *http.Request) {
	domain := r.PathValue("domain")
	id := r.PathValue("id")

	if records, ok := s.dnsRecords[domain]; ok {
		recordIndex := slices.IndexFunc(records, func(record client.DNSRecord) bool {
			return record.ID == id
		})

		if recordIndex != -1 {
			s.Write(w, client.RetrieveDNSRecordResponse{
				Status: client.Status{
					Value: "SUCCESS",
				},
				Records: []client.DNSRecord{records[recordIndex]},
			})
			return
		}
	}

	s.Write(w, client.Status{
		Value:   "FAILURE",
		Message: "Record not found",
	})
}
