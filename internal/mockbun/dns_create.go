package mockbun

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/kyswtn/terraform-provider-porkbun/internal/client"
)

// newId is a helper to make predictable id's to check with within tests
func newRecordId(domain string, r client.DNSRecord) string {
	acc := fmt.Sprintf("%s-%s-%s-%s-%s-%s", domain, r.Name, r.Type, r.Priority, r.TTL, r.Content)
	sum := md5.New()
	sum.Write([]byte(acc))
	seed := binary.BigEndian.Uint64(sum.Sum(nil))
	random := rand.New(rand.NewSource(int64(seed)))
	return fmt.Sprintf("%d", random.Uint64())
}

func (s *Server) dns_create(w http.ResponseWriter, r *http.Request) {
	domain := r.PathValue("domain")

	var record client.DNSRecord
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		log.Printf("Failed to decode request: %v", err)
		return
	}

	record.ID = newRecordId(domain, record)

	if record.Name == "" {
		record.Name = fmt.Sprintf("%s.%s", record.Name, domain)
	}

	s.dnsRecords[domain] = append(s.dnsRecords[domain], record)
	s.Write(w, client.CreateDNSRecordResponse{
		Status: client.Status{
			Value: "SUCCESS",
		},
		ID: record.ID,
	})
}
