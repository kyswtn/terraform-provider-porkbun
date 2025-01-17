resource "porkbun_dns_record" "example" {
  domain   = "example.com"
  name     = "www"
  type     = "CNAME"
  content  = "cname.vercel-dns.com"
  ttl      = 300
  priority = 1
  notes    = "Redirect www.example.com to example.com"
}
