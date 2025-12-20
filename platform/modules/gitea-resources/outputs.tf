output "private_key" {
  value     = tls_private_key.platform_ops_key.private_key_pem
  sensitive = true
}