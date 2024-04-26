resource "porkbun_dns_record" "example" {
  domain   = "example.com"
  name     = "www"
  type     = "CNAME"
  content  = "cname.vercel-dns.com"
  TTL      = 300
  Priority = 1
  Notes    = "Redirect www.example.com to example.com"
}
