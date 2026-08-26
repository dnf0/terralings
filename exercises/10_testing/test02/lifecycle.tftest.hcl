run "test_cluster_provisioning" {
  # TODO: Set command = apply
  command = plan

  variables {
    cluster_tier = "premium"
  }

  assert {
    # TODO: Fix assertion
    condition     = terraform_data.cluster.output.tier == "standard"
    error_message = "Cluster tier was not configured properly"
  }
}
