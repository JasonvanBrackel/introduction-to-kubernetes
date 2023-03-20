package iac

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type mocks int

func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	id := args.Name + "_id"
	createdResources[args.TypeToken] = args.Inputs
	return id, args.Inputs, nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

var createdResources = make(map[string]resource.PropertyMap)

var _ = Describe("Deploy Infrastructure", Label("IAC"), func() {
	var err error
	BeforeEach(func() {
		createdResources = make(map[string]resource.PropertyMap)
		err = pulumi.RunErr(func(ctx *pulumi.Context) error {
			return DeployInfrastructure(ctx)
		}, pulumi.WithMocks("https://github.com/jasonvanbrackel/introduction-to-kuberneters", "stack", mocks(0)))
	})

	It("should create a VPC.", func() {
		Expect(createdResources["aws:ec2/vpc:Vpc"]).NotTo(BeNil())
		Expect(pulumi.String(createdResources["aws:ec2/vpc:Vpc"]["cidrBlock"].StringValue())).To(Equal(pulumi.String("10.0.0.0/16")))
	})

	It("should create a Subnet.", func() {
		Expect(createdResources["aws:ec2/subnet:Subnet"]).NotTo(BeNil())
		Expect(pulumi.String(createdResources["aws:ec2/subnet:Subnet"]["cidrBlock"].StringValue())).To(Equal(pulumi.String("10.0.1.0/24")))
		Expect(pulumi.Bool(createdResources["aws:ec2/subnet:Subnet"]["mapPublicIpOnLaunch"].BoolValue())).To(Equal(pulumi.Bool(true)))
		Expect(pulumi.String(createdResources["aws:ec2/subnet:Subnet"]["vpcId"].StringValue())).To(Equal(pulumi.String("vpc-intro_to_k8s-hf_id")))
	})

	It("should create a Security Group.", func() {
		Expect(createdResources["aws:ec2/securityGroup:SecurityGroup"]).NotTo(BeNil())
		Expect(pulumi.String(createdResources["aws:ec2/securityGroup:SecurityGroup"]["vpcId"].StringValue())).To(Equal(pulumi.String("vpc-intro_to_k8s-hf_id")))
	})

	It("should create an SSH ingress rule.", func() {
		Expect(createdResources["aws:ec2/securityGroupRule:SecurityGroupRule"]).NotTo(BeNil())
		Expect(pulumi.String(createdResources["aws:ec2/securityGroupRule:SecurityGroupRule"]["type"].StringValue())).To(Equal(pulumi.String("ingress")))
		Expect(pulumi.String(createdResources["aws:ec2/securityGroupRule:SecurityGroupRule"]["protocol"].StringValue())).To(Equal(pulumi.String("tcp")))
		Expect(createdResources["aws:ec2/securityGroupRule:SecurityGroupRule"]["cidrBlocks"].ArrayValue()).To(ContainElement(pulumi.String("0.0.0.0/0")))
	})

	It("should create a Tag.", func() {
		Expect(createdResources["aws:ec2/tag:Tag"]).NotTo(BeNil())
		Expect(pulumi.String(createdResources["aws:ec2/tag:Tag"]["key"].StringValue())).To(Equal(pulumi.String(OwnerTagKey)))
		Expect(pulumi.String(createdResources["aws:ec2/tag:Tag"]["value"].StringValue())).To(Equal(pulumi.String("jvb"))) // Replace "jvb" with the value you expect for the owner
		Expect(pulumi.String(createdResources["aws:ec2/tag:Tag"]["resourceId"].StringValue())).To(Equal(pulumi.String("vpc-intro_to_k8s-hf_id")))
	})

	It("should not error.", func() {
		Expect(err).NotTo(HaveOccurred())
	})
})
