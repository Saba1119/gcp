package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// Replace these values with your project ID and zone
	projectID := "develop-375210"
	zone := "us-central1-c"
	//subnetRegion := "zone[:len(zone)-2]"

	// Create a new Compute Engine service client
	computeService, err := compute.NewService(ctx, option.WithScopes(compute.ComputeScope))
	if err != nil {
		fmt.Printf("Failed to create Compute Engine service client: %v", err)
		return
	}

	// // Define the new VPC to create
	// vpcName := "my-vpc-vm"
	// vpcCIDR := "10.0.0.0/16"

	// // Define the new subnet to create
	// subnetName := "my-subnet-vm"
	// subnetCIDR := "10.0.1.0/24"

	// Define the new VPC network to create
	networkName := "my-network-vm"
	network := &compute.Network{
		Name:                  networkName,
		AutoCreateSubnetworks: true,
	}

	// Create the new VPC network
	op, err := computeService.Networks.Insert(projectID, network).Do()
	if err != nil {
		fmt.Printf("Failed to create VPC network: %v", err)
		return
	}
	// Wait for the VPC network creation operation to complete
	for op.Status != "DONE" {
		op, err = computeService.GlobalOperations.Get(projectID, op.Name).Do()
		if err != nil {
			log.Fatalf("Failed to get VPC network creation operation status: %v", err)
		}
	}

	fmt.Printf("VPC network %q is created.\n", networkName)

	// // Define the new subnets to create
	// subnetName1 := "my-subnet-1-vm"
	// subnet1 := &compute.Subnetwork{
	// 	Name:        subnetName1,
	// 	IpCidrRange: "10.0.1.0/24",
	// 	Network:     fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName),
	// 	Region:      zone[:len(zone)-2],
	// }

	// subnetName2 := "my-subnet-2-vm"
	// subnet2 := &compute.Subnetwork{
	// 	Name:        subnetName2,
	// 	IpCidrRange: "10.0.2.0/24",
	// 	Network:     fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName),
	// 	Region:      zone[:len(zone)-2],
	// }

	// // Create the new subnets
	// _, err = computeService.Subnetworks.Insert(projectID, zone, subnet1).Do()
	// if err != nil {
	// 	fmt.Printf("Failed to create subnet %q: %v", subnetName1, err)
	// 	return
	// }

	// fmt.Printf("Subnets %q  is created.\n", subnetName1)

	// _, err = computeService.Subnetworks.Insert(projectID, zone, subnet2).Do()
	// if err != nil {
	// 	fmt.Printf("Failed to create subnet %q: %v", subnetName2, err)
	// 	return
	// }

	//fmt.Printf("Subnets %q  is created.\n", subnetName1)

	subnetName1 := "my-subnet-1-vm"
	region := "us-central1"

	// Define the new instance to create
	instanceName := "my-instance-vm"
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
				Network:    fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName),
				Subnetwork: fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", projectID, region, subnetName1),
				AccessConfigs: []*compute.AccessConfig{
					{
						Type:        "ONE_TO_ONE_NAT",
						Name:        "External NAT",
						NetworkTier: "PREMIUM",
					},
				},
			},
		},
		Tags: &compute.Tags{
			Items: []string{"http-server", "https-server"},
		},
	}

	// // Define the new VPC to create
	// vpcName := "my-vpc-vm"
	// vpcCIDR := "10.0.0.0/16"

	// // Create the new VPC
	// vpc, err := computeService.Networks.Insert(projectID, &compute.Network{
	// 	Name:      vpcName,
	// 	IPv4 : vpcCIDR,
	// }).Do()
	// if err != nil {
	// 	fmt.Printf("Failed to create VPC: %v", err)
	// 	return
	// }

	// fmt.Printf("VPC %q is created with selfLink %q.\n", vpc.Name, vpc.SelfLink)

	// // // Define the new subnet to create
	// // subnetName := "my-subnet-vm"
	// // subnetCIDR := "10.0.1.0/24"

	// // Create the new subnet
	// subnet, err := computeService.Subnetworks.Insert(projectID, zone, &compute.Subnetwork{
	// 	Name:        subnetName,
	// 	Network:     vpc.SelfLink,
	// 	IpCidrRange: subnetCIDR,
	// }).Do()
	// if err != nil {
	// 	fmt.Printf("Failed to create subnet: %v", err)
	// 	return
	// }

	// fmt.Printf("Subnet %q is created with selfLink %q.\n", subnet.Name, subnet.SelfLink)

	// Define the new subnets to create
	// subnetName1 := "my-subnet-1-vm"
	subnet1 := &compute.Subnetwork{
		Name:        subnetName1,
		IpCidrRange: "10.0.1.0/24",
		Network:     fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName),
		Region:      "us-central1",
	}
	// Create the new subnets
	_, err = computeService.Subnetworks.Insert(projectID, region, subnet1).Do()
	if err != nil {
		fmt.Printf("Failed to create subnet %q: %v", subnetName1, err)
		return
	}
	fmt.Printf("Subnets %q  is created.\n", subnetName1)

	// Wait for the operation to complete.
	for {
		time.Sleep(60 * time.Second)
		op, err = computeService.GlobalOperations.Get(projectID, op.Name).Do()
		if err != nil {
			panic(err)
		}

		if op.Status == "DONE" {
			if op.Error != nil {
				panic(fmt.Sprintf("Operation failed: %v", op.Error))
			}
		}

		// Create the new instance
		op, err := computeService.Instances.Insert(projectID, zone, instance).Do()
		if err != nil {
			fmt.Printf("Failed to create instance: %v", err)
			return
		}

		fmt.Printf("Instance %q is being created with operation %q\n", instanceName, op.Name)

		// Wait for the instance to be created
		fmt.Printf("Waiting for the instance to be created...\n")
		for {
			op, err := computeService.ZoneOperations.Get(projectID, zone, op.Name).Do()
			if err != nil {
				fmt.Printf("Failed to get zone operation %q: %v", op.Name, err)
				return
			}

			if op.Status == "DONE" {
				if op.Error != nil {
					fmt.Printf("Failed to create instance: %v", op.Error.Errors)
					return
				}

				fmt.Printf("Instance %q is created.\n", instanceName)
				break
			}

			time.Sleep(1 * time.Second)
		}

		// Define the new firewall rule to allow ingress and egress traffic on port 80 and 443
		firewallName := "my-firewal-rule-vm"
		firewall := &compute.Firewall{
			Name: firewallName,
			Allowed: []*compute.FirewallAllowed{
				{
					IPProtocol: "tcp",
					Ports:      []string{"8000", "8080"},
				},
			},
			Network: fmt.Sprintf("projects/%s/global/networks/default", projectID),
			SourceRanges: []string{
				"0.0.0.0/0",
			},
			DestinationRanges: []string{
				"0.0.0.0/0",
			},
			TargetTags: []string{"http-server", "https-server"},
		}

		// Create the new firewall rule
		_, err = computeService.Firewalls.Insert(projectID, firewall).Do()
		if err != nil {
			fmt.Printf("Failed to create firewall rule: %v", err)
			return
		}

		fmt.Printf("Firewall rule %q is created.\n", firewallName)

		// Get the new instance to update its network tags
		newInstance, err := computeService.Instances.Get(projectID, zone, instanceName).Do()
		if err != nil {
			fmt.Printf("Failed to get instance: %v", err)
			return
		}

		// Add the network tags to the new instance
		newInstance.Tags = &compute.Tags{
			Items: []string{
				"web-server",
			},
		}

		// Update the new instance to add the network tags
		_, err = computeService.Instances.Update(projectID, zone, instanceName, newInstance).Do()
		if err != nil {
			fmt.Printf("Failed to update instance: %v", err)
			return
		}

		fmt.Printf("Instance %q is updated with network tags.\n", instanceName)
	}
}
