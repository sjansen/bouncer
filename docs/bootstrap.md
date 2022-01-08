1. Build Cloudfront viewer-request Lambda
```
cd cloudfront/viewer-request/
yarn build
```

2. Configure Terraform

    a. `.env.local`
    ```
    AWS_PROFILE=staging
    ```

    b. `terraform/env/terragrunt-local.json`
    ```
    {
      "prefix": "example",
      "providers": {
        "aws": {
          "profile": "staging"
        },
        "route53": {
          "profile": "staging"
        }
      }
    }
    ```

    c. `terraform/env/$ENV/terraform.tfvars`
    ```
    dns-name = "docs.example.com"
    dns-zone = "example.com"
    ```

3. Create AWS resources
```
cd terraform/env/staging/
terragrunt apply
```

4. Generate or purchase cert & key
```
scripts/gen-test-cert
```

5. Add config data to AWS SSM Parameter Store

https://console.aws.amazon.com/systems-manager/parameters/


6. Build & upload Go Lambdas
```
make staging
```

7. Rotate JWK keys
```
TODO
```
