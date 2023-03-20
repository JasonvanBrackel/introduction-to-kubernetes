package main

import (
	"github.com/jasonvanbrackel/introduction-to-kubernetes/pulumi/iac"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := iac.DeployInfrastructure(ctx)
		return err
	})
}
