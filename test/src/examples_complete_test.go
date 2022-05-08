package test

import (
  "os"
  "strings"
  "testing"

  "github.com/gruntwork-io/terratest/modules/random"
  "github.com/gruntwork-io/terratest/modules/terraform"
  test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
  "github.com/stretchr/testify/assert"
)

func cleanup(t *testing.T, terraformOptions *terraform.Options, tempTestFolder string) {
  terraform.Destroy(t, terraformOptions)
  os.RemoveAll(tempTestFolder)
}

// Test the Terraform module in examples/complete using Terratest.
func TestExamplesComplete(t *testing.T) {
  t.Parallel()
  randID := strings.ToLower(random.UniqueId())
  attributes := []string{randID}

  rootFolder := "../../"
  terraformFolderRelativeToRoot := "examples/complete"
  varFiles := []string{"fixtures.us-east-2.tfvars"}

  tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)

  terraformOptions := &terraform.Options{
    // The path to where our Terraform code is located
    TerraformDir: tempTestFolder,
    Upgrade:      true,
    // Variables to pass to our Terraform code using -var-file options
    VarFiles: varFiles,
    Vars: map[string]interface{}{
      "attributes": attributes,
    },
  }

  // At the end of the test, run `terraform destroy` to clean up any resources that were created
  defer cleanup(t, terraformOptions, tempTestFolder)

  // This will run `terraform init` and `terraform apply` and fail the test if there are any errors
  terraform.InitAndApply(t, terraformOptions)

  zone_name := "cloudposse-terraform-aws-route53.cloudposse.com"
  // Run `terraform output` to get the value of an output variable
  zone_name_output := terraform.Output(t, terraformOptions, "zone_name")

  // Verify we're getting back the outputs we expect
  // Ensure zone successfully created
  assert.Equal(t, zone_name, zone_name_output)

}