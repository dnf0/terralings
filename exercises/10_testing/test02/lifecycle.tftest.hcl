run "test_cluster_provisioning" {
  # TODO (What): Set command = apply.
  # TODO (Why): Using command = apply provisions resources in memory/ephemeral state to assert against real computed attributes rather than unknown plan values.
  command = plan

  variables {
    cluster_tier = "premium"
  }

  assert {
    # TODO (What): Change condition to terraform_data.cluster.output.tier == "premium".
    # TODO (Why): The test must verify that the applied cluster resource output accurately reflects the input variable value.
    condition     = terraform_data.cluster.output.tier == "standard"
    error_message = "Cluster tier was not configured properly"
  }
}
