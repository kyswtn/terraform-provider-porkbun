package mockbun

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	porkbun "github.com/kyswtn/terraform-provider-porkbun/internal/client"
	"google.golang.org/appengine/log"
)

type Server struct {
	mux         *http.ServeMux
	server      *httptest.Server
	URL         string
	nameservers map[string][]string
	dnsRecords  map[string][]porkbun.DNSRecord
}

func New() *Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	m := &Server{
		mux:         mux,
		server:      server,
		URL:         server.URL,
		nameservers: make(map[string][]string),
		dnsRecords:  make(map[string][]porkbun.DNSRecord),
	}

	m.addPorkbunHandlers()
	return m
}

func (m *Server) Close() {
	m.server.Close()
}

func (m *Server) SetNameservers(domain string, nameservers []string) {
	m.nameservers[domain] = nameservers
}

func (m *Server) SetDNSRecords(domain string, records []porkbun.DNSRecord) {
	m.dnsRecords[domain] = records
}

func (s *Server) Write(writer http.ResponseWriter, body any) {
	writer.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(body)
	if err != nil {
		log.Warningf(context.TODO(), "Failed to marshal response: %v", err)
	}

	_, err = writer.Write(data)
	if err != nil {
		log.Warningf(context.TODO(), "Failed to write response: %v", err)
	}
}

func (s *Server) addPorkbunHandlers() {
	s.mux.HandleFunc("/domain/getNs/{domain}", s.domain_get_ns)
	s.mux.HandleFunc("/domain/updateNs/{domain}", s.domain_update_ns)
	s.mux.HandleFunc("/dns/create/{domain}", s.dns_create)
	s.mux.HandleFunc("/dns/retrieve/{domain}/{id}", s.dns_retrieve)
	s.mux.HandleFunc("/dns/edit/{domain}/{id}", s.dns_edit)
	s.mux.HandleFunc("/dns/delete/{domain}/{id}", s.dns_delete)
}
