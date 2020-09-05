#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { CdkLambdaGoStack } from '../lib/cdk-lambda-go-stack';

const app = new cdk.App();
new CdkLambdaGoStack(app, 'CdkLambdaGoStack');
