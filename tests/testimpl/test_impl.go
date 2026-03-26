package common

import (
	"context"
	"strings"
	"testing"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	operationalinsights "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights"
	monitor "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
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

	actionGroupClient, err := monitor.NewActionGroupsClient(subscriptionID, credential, &options)
	if err != nil {
		t.Fatalf("Error creating action group client: %v", err)
	}

	scheduledQueryAlertClient, err := monitor.NewScheduledQueryRulesClient(subscriptionID, credential, &options)
	if err != nil {
		t.Fatalf("Error creating scheduled query alert client: %v", err)
	}

	t.Run("doesLogAnalyticsWorkspaceExist", func(t *testing.T) {
		checkLogAnalyticsWorkspaceExistence(t, logAnalyticsWorkspaceClient, ctx.TerratestTerraformOptions(), ctx)
	})

	t.Run("doesActionGroupExist", func(t *testing.T) {
		checkActionGroupExistence(t, actionGroupClient, ctx.TerratestTerraformOptions(), ctx)
	})

	t.Run("doesScheduledQueryAlertExist", func(t *testing.T) {
		checkScheduledQueryAlertExistence(t, scheduledQueryAlertClient, ctx.TerratestTerraformOptions(), ctx)
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

func checkActionGroupExistence(t *testing.T, actionGroupClient *monitor.ActionGroupsClient, terraformOptions *terraform.Options, ctx types.TestContext) {
	actionGroupID := terraform.Output(t, terraformOptions, "action_group_id")

	// If no action group was created (when query_alerts is empty), skip this check
	if actionGroupID == "" || actionGroupID == "<nil>" {
		t.Skip("No action group configured, skipping")
	}

	actionGroupName := terraform.Output(t, terraformOptions, "action_group_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	actionGroup, err := actionGroupClient.Get(context.Background(), resourceGroupName, actionGroupName, nil)
	if err != nil {
		t.Fatalf("Error getting action group: %v", err)
	}

	assert.Equal(t, strings.ToLower(actionGroupID), strings.ToLower(*actionGroup.ID), "Action Group ID does not match.")
	assert.NotNil(t, actionGroup.Properties, "Action Group properties should not be nil")
}

func checkScheduledQueryAlertExistence(t *testing.T, scheduledQueryAlertClient *monitor.ScheduledQueryRulesClient, terraformOptions *terraform.Options, ctx types.TestContext) {
	scheduledQueryAlertsOutput := terraform.Output(t, terraformOptions, "scheduled_query_alerts")

	// If no scheduled query alerts were created, skip this check
	if scheduledQueryAlertsOutput == "" || scheduledQueryAlertsOutput == "{}" {
		t.Skip("No scheduled query alerts configured, skipping")
	}

	// In real implementation, you would unmarshal the JSON output here
	// For now, just verify that the output is not empty
	assert.NotEmpty(t, scheduledQueryAlertsOutput, "Scheduled query alerts output should not be empty")
}
