package common

import (
	"context"
	"testing"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	operationalinsights "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

func TestLogAnalyticsWorkspace(t *testing.T, ctx types.TestContext) {

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID is not set in the environment variables ")
	}

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %e\n", err)
	}

	options := arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzurePublic,
		},
	}

	logAnalyticsWorkspaceClient, err := operationalinsights.NewWorkspacesClient(subscriptionID, credential, &options)
	if err != nil {
		t.Fatalf("Error creating log analytics workspace client: %v", err)
	}

	t.Run("doesLogAnalyticsWorkspaceExist", func(t *testing.T) {
		checkLogAnalyticsWorkspaceExistence(t, logAnalyticsWorkspaceClient, ctx.TerratestTerraformOptions(), ctx)
	})
}

func checkLogAnalyticsWorkspaceExistence(t *testing.T, logAnalyticsWorkspaceClient *operationalinsights.WorkspacesClient, terraformOptions *terraform.Options, ctx types.TestContext) {
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	workspaceName := terraform.Output(t, terraformOptions, "workspace_name")
	id := terraform.Output(t, terraformOptions, "id")

	workspace, err := logAnalyticsWorkspaceClient.Get(context.Background(), resourceGroupName, workspaceName, nil)
	if err != nil {
		t.Fatalf("Error getting log analytics workspace: %v", err)
	}

	assert.Equal(t, id, *workspace.ID, "Log Analytics ID does not match.")
}
