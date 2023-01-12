import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as iam from "aws-cdk-lib/aws-iam";
import { Fn } from "aws-cdk-lib";

import * as apigateway from "aws-cdk-lib/aws-apigateway"

// import * as sqs from 'aws-cdk-lib/aws-sqs';

export class CdkStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here

    const pokeAppId = "main"

    const testLambda = new lambda.Function(this, pokeAppId, {
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

    const api = new apigateway.RestApi(this, "pokeapi", {
      restApiName: "PokeAPI",
      description: "This service serves the PokeAPI",
      deployOptions: {
        stageName: "dev",
      }
    });

    const integration = new apigateway.LambdaIntegration(testLambda, {
      proxy: true,
      // integrationResponses: [{ statusCode: "200" }]
    });

    api.root.addProxy({
      defaultIntegration: integration,
      defaultMethodOptions: {
        // authorizationType: apigateway.AuthorizationType.NONE,
        // apiKeyRequired: false,
      },
      anyMethod: true,
    })


    // example resource
    // const queue = new sqs.Queue(this, 'CdkQueue', {
    //   visibilityTimeout: cdk.Duration.seconds(300)
    // });
  }
}
