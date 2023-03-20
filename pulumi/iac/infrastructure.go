package iac

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const OwnerTagKey = "Owner"
const SecurityGroupName = "intro_to_k8s-hf"
const VpcName = "vpc-intro_to_k8s-hf"
const SubnetName = "subnet–intro_to_k8s-hf"

func DeployInfrastructure(ctx *pulumi.Context) error {
	owner, ok := ctx.GetConfig("owner")
	if !ok {
		owner = "jvb"
	}

	// Create a new VPC
	vpc, err := ec2.NewVpc(ctx, VpcName, &ec2.VpcArgs{
		CidrBlock: pulumi.String("10.0.0.0/16"),
	})
	if err != nil {
		return err
	}

	_, err = ec2.NewSubnet(ctx, SubnetName, &ec2.SubnetArgs{
		CidrBlock:           pulumi.String("10.0.1.0/24"),
		MapPublicIpOnLaunch: pulumi.Bool(true),
		VpcId:               vpc.ID(),
	})

	// Create a new Security Group
	sg, err := ec2.NewSecurityGroup(ctx, SecurityGroupName, &ec2.SecurityGroupArgs{
		VpcId: vpc.ID(),
	})
	if err != nil {
		return err
	}

	_, err = ec2.NewSecurityGroupRule(ctx, "ssh", &ec2.SecurityGroupRuleArgs{
		SecurityGroupId: sg.ID(),
		Type:            pulumi.String("ingress"),
		Protocol:        pulumi.String("tcp"),
		FromPort:        pulumi.Int(22),
		ToPort:          pulumi.Int(22),
		CidrBlocks:      pulumi.StringArray{pulumi.String("0.0.0.0/0")},
	})
	if err != nil {
		return err
	}

	// Add a tag to the VPC
	_, err = ec2.NewTag(ctx, OwnerTagKey, &ec2.TagArgs{
		ResourceId: vpc.ID(),
		Key:        pulumi.String(OwnerTagKey),
		Value:      pulumi.String(owner),
	})
	if err != nil {
		return err
	}

	// Export the VPC and Security Group IDs
	ctx.Export("vpcID", vpc.ID())
	ctx.Export("securityGroupID", sg.ID())

	return nil
}
