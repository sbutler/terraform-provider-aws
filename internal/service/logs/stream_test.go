package logs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tflogs "github.com/hashicorp/terraform-provider-aws/internal/service/logs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccLogsStream_basic(t *testing.T) {
	var ls cloudwatchlogs.LogStream
	resourceName := "aws_cloudwatch_log_stream.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamExists(resourceName, &ls),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccStreamImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLogsStream_disappears(t *testing.T) {
	var ls cloudwatchlogs.LogStream
	resourceName := "aws_cloudwatch_log_stream.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamExists(resourceName, &ls),
					acctest.CheckResourceDisappears(acctest.Provider, tflogs.ResourceStream(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccLogsStream_Disappears_logGroup(t *testing.T) {
	var ls cloudwatchlogs.LogStream
	var lg cloudwatchlogs.LogGroup
	resourceName := "aws_cloudwatch_log_stream.test"
	logGroupResourceName := "aws_cloudwatch_log_group.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, cloudwatchlogs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStreamDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStreamExists(resourceName, &ls),
					testAccCheckGroupExists(logGroupResourceName, &lg),
					acctest.CheckResourceDisappears(acctest.Provider, tflogs.ResourceGroup(), logGroupResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckStreamExists(n string, v *cloudwatchlogs.LogStream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CloudWatch Logs Log Stream ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LogsConn

		output, err := tflogs.FindLogStreamByTwoPartKey(context.Background(), conn, rs.Primary.Attributes["log_group_name"], rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckStreamDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).LogsConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cloudwatch_log_stream" {
			continue
		}

		_, err := tflogs.FindLogStreamByTwoPartKey(context.Background(), conn, rs.Primary.Attributes["log_group_name"], rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("CloudWatch Logs Log Stream still exists: %s", rs.Primary.ID)
	}

	return nil
}

func testAccStreamImportStateIdFunc(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["log_group_name"], rs.Primary.ID), nil
	}
}

func testAccStreamConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_log_group" "test" {
  name = %[1]q
}

resource "aws_cloudwatch_log_stream" "test" {
  name           = %[1]q
  log_group_name = aws_cloudwatch_log_group.test.id
}
`, rName)
}
