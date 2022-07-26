package ds_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directoryservice"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfds "github.com/hashicorp/terraform-provider-aws/internal/service/ds"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccDSSharedDirectoryAccepter_basic(t *testing.T) {
	var providers []*schema.Provider
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_directory_service_shared_directory_accepter.test"

	domainName := acctest.RandomDomainName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, directoryservice.EndpointsID),
		ProviderFactories: acctest.FactoriesAlternate(&providers),
		CheckDestroy:      testAccCheckSharedDirectoryAccepterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSharedDirectoryAccepterConfig_basic(rName, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSharedDirectoryAccepterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "method", directoryservice.ShareMethodHandshake),
					resource.TestCheckResourceAttr(resourceName, "notes", "There were hints and allegations"),
					resource.TestCheckResourceAttrPair(resourceName, "owner_account_id", "data.aws_caller_identity.current", "account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "owner_directory_id"),
					resource.TestCheckResourceAttrSet(resourceName, "shared_directory_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notes",
				},
			},
		},
	})

}

func testAccCheckSharedDirectoryAccepterExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return names.Error(names.DS, names.ErrActionCheckingExistence, tfds.ResourceNameSharedDirectoryAccepter, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return names.Error(names.DS, names.ErrActionCheckingExistence, tfds.ResourceNameSharedDirectoryAccepter, name, errors.New("no ID is set"))
		}

		ownerId := rs.Primary.Attributes["owner_directory_id"]
		sharedId := rs.Primary.Attributes["shared_directory_id"]

		conn := acctest.Provider.Meta().(*conns.AWSClient).DSConn
		out, err := conn.DescribeSharedDirectories(&directoryservice.DescribeSharedDirectoriesInput{
			OwnerDirectoryId:   aws.String(ownerId),
			SharedDirectoryIds: aws.StringSlice([]string{sharedId}),
		})

		if err != nil {
			return names.Error(names.DS, names.ErrActionCheckingExistence, tfds.ResourceNameSharedDirectoryAccepter, name, err)
		}

		if len(out.SharedDirectories) < 1 {
			return names.Error(names.DS, names.ErrActionCheckingExistence, tfds.ResourceNameSharedDirectoryAccepter, name, errors.New("not found"))
		}

		if aws.StringValue(out.SharedDirectories[0].SharedDirectoryId) != sharedId {
			return names.Error(names.DS, names.ErrActionCheckingExistence, tfds.ResourceNameSharedDirectoryAccepter, rs.Primary.ID, fmt.Errorf("shared directory ID mismatch - existing: %q, state: %q", aws.StringValue(out.SharedDirectories[0].SharedDirectoryId), sharedId))
		}

		if aws.StringValue(out.SharedDirectories[0].OwnerDirectoryId) != ownerId {
			return names.Error(names.DS, names.ErrActionCheckingExistence, tfds.ResourceNameSharedDirectoryAccepter, rs.Primary.ID, fmt.Errorf("owner directory ID mismatch - existing: %q, state: %q", aws.StringValue(out.SharedDirectories[0].OwnerDirectoryId), ownerId))
		}

		return nil
	}

}

func testAccCheckSharedDirectoryAccepterDestroy(s *terraform.State) error {
	// cannot be destroyed from consumer account
	return nil
}

func testAccSharedDirectoryAccepterConfig_basic(rName, domain string) string {
	return acctest.ConfigCompose(
		acctest.ConfigAlternateAccountProvider(),
		testAccDirectoryConfig_microsoftStandard(rName, domain),
		`
data "aws_caller_identity" "current" {}

resource "aws_directory_service_shared_directory" "test" {
  directory_id = aws_directory_service_directory.test.id
  notes        = "There were hints and allegations"

  target {
    id = data.aws_caller_identity.consumer.account_id
  }
}

data "aws_caller_identity" "consumer" {
  provider = "awsalternate"
}

resource "aws_directory_service_shared_directory_accepter" "test" {
  provider = "awsalternate"

  shared_directory_id = aws_directory_service_shared_directory.test.shared_directory_id
}
`)
}
