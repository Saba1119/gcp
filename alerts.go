package main

import (
    "context"
    "fmt"
    "log"

    "cloud.google.com/go/logging"
    "google.golang.org/api/option"
    "google.golang.org/api/monitoring/v3"
)

const (
    projectID = "my-project-id"
)

func main() {
    ctx := context.Background()

    // Create a Stackdriver Logging client
    logClient, err := logging.NewClient(ctx, projectID)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer logClient.Close()

    // Create a Stackdriver Monitoring client
    monitoringService, err := monitoring.NewService(ctx, option.WithCredentialsFile("path/to/creds.json"))
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Define a new alert policy
    alertPolicy := &monitoring.AlertPolicy{
        DisplayName: "High CPU usage",
        Conditions: []*monitoring.AlertPolicy_Condition{
            {
                DisplayName: "High CPU usage",
                ConditionThreshold: &monitoring.AlertPolicy_Condition_MetricThreshold{
                    Comparison: monitoring.Comparison_GT,
                    ThresholdValue: &monitoring.TypedValue{
                        DoubleValue: 1.0,
                    },
                    TimeWindow: &monitoring.Duration{
                        Seconds: 60,
                    },
                    Metric: &monitoring.Metric{
                        Type: "compute.googleapis.com/instance/cpu/utilization",
                        Labels: map[string]string{
                            "instance_name": "my-instance-vm",
                        },
                    },
                    Aggregations: []*monitoring.Aggregation{
                        {
                            AlignmentPeriod: &monitoring.Duration{
                                Seconds: 60,
                            },
                            PerSeriesAligner: monitoring.Aggregation_Aligner_ALIGN_MAX,
                        },
                    },
                },
            },
        },
        CombinationOfPolicies: &monitoring.AlertPolicy_CombinationOfPolicies{
            Conditions: []*monitoring.AlertPolicy_CombinationCondition{
                {
                    Condition: &monitoring.AlertPolicy_Condition{
                        Name: "projects/" + projectID + "/alertPolicies/" + "policy-id-1" + "/conditions/0",
                    },
                    Trigger: &monitoring.AlertPolicy_CombinationCondition_Trigger{},
                },
                {
                    Condition: &monitoring.AlertPolicy_Condition{
                        Name: "projects/" + projectID + "/alertPolicies/" + "policy-id-2" + "/conditions/0",
                    },
                    Trigger: &monitoring.AlertPolicy_CombinationCondition_Trigger{},
                },
            },
            Combination: monitoring.AlertPolicy_CombinationOfPolicies_OR,
        },
        Enabled: true,
    }

    // Create the alert policy
    resp, err := monitoringService.Projects.AlertPolicies.Create("projects/"+projectID, alertPolicy).Do()
    if err != nil {
        log.Fatalf("Failed to create alert policy: %v", err)
    }

    // Print the alert policy ID
    fmt.Printf("Created alert policy %v\n", resp.Name)

    // Create a new log entry
    logEntry := logging.Entry{
        Payload: "High CPU usage on my-instance",
        Labels: map[string]string{
            "type": "error",
        },
    }
    // Write the log entry
    if _, err := logClient.Logger("my-log").Log(logEntry); err != nil {
        log.Fatalf("Failed to write log entry: %v", err)
    }
}
