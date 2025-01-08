package mockbun

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/kyswtn/terraform-provider-porkbun/internal/client"
)

func (s *Server) dns_edit(w http.ResponseWriter, r *http.Request) {
	domain := r.PathValue("domain")
	id := r.PathValue("id")

	if records, ok := s.dnsRecords[domain]; ok {
		recordIndex := slices.IndexFunc(records, func(record client.DNSRecord) bool {
			return record.ID == id
		})

		if recordIndex != -1 {
			var data client.DNSRecord
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				goto dns_edit_failure
			}

			// Merge the existing record with the new data and replace the object in the dnsRecords
			oldRecord := s.dnsRecords[domain][recordIndex]
			s.dnsRecords[domain][recordIndex] = oldRecord.Merge(data)

			s.Write(w, client.Status{
				Value: "SUCCESS",
			})
			return
		}
	}

dns_edit_failure:
	s.Write(w, client.Status{
		Value:   "FAILURE",
		Message: "Record not found",
	})
}
