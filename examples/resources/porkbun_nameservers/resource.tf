resource "porkbun_nameservers" "example" {
  domain = "example.com"
  nameservers = [
    "jim.ns.cloudflare.com",
    "pam.ns.cloudflare.com"
  ]
}
