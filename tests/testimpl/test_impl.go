package common

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	monitor "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	operationalinsights "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Expected values mirror the alert configured in examples/complete/test.tfvars.
const (
	expectedAlertKey            = "ingestion"
	expectedQuery               = "Usage | project TimeGenerated, DataType, Quantity"
	expectedEvaluationFrequency = "PT5M"
	expectedWindowDuration      = "PT30M"
	expectedMuteDuration        = "PT1H"
	expectedMetricMeasureColumn = "Quantity"
	expectedDimensionName       = "DataType"
	expectedEmailSubject        = "Azure Alert: Log Analytics ingestion volume"
	expectedCustomPropertyKey   = "alert_source"
	expectedCustomPropertyValue = "reference_log_analytics_workspace"
	expectedThreshold           = float64(1000)
)

func TestComposableLogAnalyticsWorkspace(t *testing.T, ctx types.TestContext) {

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
	alerts := terraform.OutputMapOfObjects(t, terraformOptions, "scheduled_query_alerts")
	if len(alerts) == 0 {
		t.Skip("No scheduled query alerts configured, skipping")
	}

	entry, ok := alerts[expectedAlertKey].(map[string]interface{})
	require.True(t, ok, "scheduled_query_alerts should contain the %q alert", expectedAlertKey)
	alertName, ok := entry["name"].(string)
	require.True(t, ok, "the %q alert should expose a name", expectedAlertKey)

	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	workspaceID := terraform.Output(t, terraformOptions, "id")
	actionGroupID := terraform.Output(t, terraformOptions, "action_group_id")

	rule, err := scheduledQueryAlertClient.Get(context.Background(), resourceGroupName, alertName, nil)
	require.NoError(t, err, "Error getting scheduled query alert")
	assert.Equal(t, alertName, *rule.Name)

	properties := rule.Properties
	require.NotNil(t, properties, "Scheduled query alert properties should not be nil")

	// The alert is scoped to the workspace this module creates, not to a child resource.
	assert.Equal(t, []string{workspaceID}, derefAll(properties.Scopes))
	assert.Equal(t, monitor.AlertSeverityTwo, *properties.Severity)
	assert.Equal(t, expectedEvaluationFrequency, *properties.EvaluationFrequency)
	assert.Equal(t, expectedWindowDuration, *properties.WindowSize)
	assert.Equal(t, true, *properties.Enabled)

	// Re-notifies on a fixed cadence rather than resolving itself.
	assert.Equal(t, false, *properties.AutoMitigate)
	assert.Equal(t, expectedMuteDuration, *properties.MuteActionsDuration)

	require.NotNil(t, properties.Criteria)
	require.Len(t, properties.Criteria.AllOf, 1)
	condition := properties.Criteria.AllOf[0]

	assert.Equal(t, expectedQuery, *condition.Query)
	assert.Equal(t, monitor.ConditionOperatorGreaterThan, *condition.Operator)
	assert.Equal(t, monitor.TimeAggregationTotal, *condition.TimeAggregation)
	assert.Equal(t, expectedThreshold, *condition.Threshold)
	assert.Equal(t, expectedMetricMeasureColumn, *condition.MetricMeasureColumn)

	require.Len(t, condition.Dimensions, 1)
	assert.Equal(t, expectedDimensionName, *condition.Dimensions[0].Name)
	assert.Equal(t, monitor.DimensionOperatorInclude, *condition.Dimensions[0].Operator)
	assert.Equal(t, []string{"*"}, derefAll(condition.Dimensions[0].Values))

	require.NotNil(t, condition.FailingPeriods)
	assert.Equal(t, int64(1), *condition.FailingPeriods.MinFailingPeriodsToAlert)
	assert.Equal(t, int64(1), *condition.FailingPeriods.NumberOfEvaluationPeriods)

	require.NotNil(t, properties.Actions)
	assert.Equal(t, []string{actionGroupID}, derefAll(properties.Actions.ActionGroups))

	// email_subject reaches the API as an action property, not a first class field.
	require.Contains(t, properties.Actions.ActionProperties, "Email.Subject")
	assert.Equal(t, expectedEmailSubject, *properties.Actions.ActionProperties["Email.Subject"])

	require.Contains(t, properties.Actions.CustomProperties, expectedCustomPropertyKey)
	assert.Equal(t, expectedCustomPropertyValue, *properties.Actions.CustomProperties[expectedCustomPropertyKey])
}

func derefAll(in []*string) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		out = append(out, *v)
	}
	return out
}
