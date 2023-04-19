package main

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Replace these values with your project ID and zone
	projectID := "develop-375210"
	zone := "us-central1-a"

	// Create a new Compute Engine service client
	computeService, err := compute.NewService(ctx, option.WithScopes(compute.ComputeScope))
	if err != nil {
		fmt.Printf("Failed to create Compute Engine service client: %v", err)
		return
	}

	// Define the new instance to create
	instanceName := "my-instance-saba"
	machineType := fmt.Sprintf("zones/%s/machineTypes/f1-micro", zone)
	sourceImage := "projects/ubuntu-os-cloud/global/images/family/ubuntu-2204-lts"
	diskSizeGb := int64(11)

	instance := &compute.Instance{
		Name:        instanceName,
		MachineType: machineType,
		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Mode:       "READ_WRITE",
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: sourceImage,
					DiskSizeGb:  diskSizeGb,
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
				AccessConfigs: []*compute.AccessConfig{
					{
						Type:        "ONE_TO_ONE_NAT",
						Name:        "External NAT",
						NetworkTier: "PREMIUM",
					},
				},
				Subnetwork: fmt.Sprintf("projects/%s/regions/%s/subnetworks/default", projectID, zone[:len(zone)-2]),
			},
		},
	}

	// Create the new instance
	op, err := computeService.Instances.Insert(projectID, zone, instance).Do()
	if err != nil {
		fmt.Printf("Failed to create instance: %v", err)
		return
	}

	fmt.Printf("Instance %q is being created with operation %q\n", instanceName, op.Name)
}
