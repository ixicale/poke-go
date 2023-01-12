import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as iam from "aws-cdk-lib/aws-iam";
import { Fn } from "aws-cdk-lib";

// import * as sqs from 'aws-cdk-lib/aws-sqs';

export class CdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here

    const pokeAppId = "main"

    new lambda.Function(this, pokeAppId, {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset("app"),
      handler: "main",
      timeout: cdk.Duration.seconds(300),
      memorySize: 1024,
      role: iam.Role.fromRoleArn(
        this,
        "lamda-pokeapi-role", // LambdaRole
        Fn.importValue("core-data-common-CoreDataLambdaRole-dev-Arn")
      )
    })
    // example resource
    // const queue = new sqs.Queue(this, 'CdkQueue', {
    //   visibilityTimeout: cdk.Duration.seconds(300)
    // });
  }
}
