import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';
import * as iam from '@aws-cdk/aws-iam';

export class CdkLambdaGoStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const myFuncIamRole = new iam.Role(this, 'MyFuncIAMRole', {
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com')
    })

    const adminStatement = new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      actions: ['*'],
      resources: ['*']
    })

    const myFuncIamPolicy = new iam.Policy(this, 'AdminForMyFunc');
    myFuncIamPolicy.addStatements(adminStatement)

    myFuncIamRole.attachInlinePolicy(myFuncIamPolicy)

    const myfunc = new lambda.Function(this, 'MyFuncHandler', {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset('lambda/bin'),
      handler: 'main',
      role: myFuncIamRole,
      timeout: cdk.Duration.minutes(1),
      tracing: lambda.Tracing.ACTIVE
    })
  }
}
