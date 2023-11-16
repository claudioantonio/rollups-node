target "default" {
  tags = ["cartesi/rollups-node-devnet:devel"]
  args = {
    NODE_DEVNET_VERSION = "devel"
  }
}
