target "default" {
  tags = ["cartesi/anvil_devnet:devel"]
  args = {
    NODE_ANVIL_VERSION = "devel"
  }
}
