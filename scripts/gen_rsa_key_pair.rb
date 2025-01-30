#!/usr/bin/ruby

require 'openssl'
require 'sshkey'
require 'base64'

private_key = OpenSSL::PKey::RSA.new(2048)

puts "RAW PRIVATE"
puts private_key.to_pem
puts ""

puts "BASE64 ENCODED PRIVATE"
puts Base64.urlsafe_encode64(private_key.to_pem)
puts ""

puts "BASE64 ENCODED PUBLIC"
puts Base64.urlsafe_encode64(private_key.public_key.to_pem)
puts ""