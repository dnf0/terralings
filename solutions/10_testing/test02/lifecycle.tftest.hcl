run "test_cluster_provisioning" {
  command = apply

  variables {
    cluster_tier = "premium"
  }

  assert {
    condition     = terraform_data.cluster.output.tier == "premium" && terraform_data.cluster.output.nodes == 5
    error_message = "Cluster tier and node count were not configured properly"
  }
}
