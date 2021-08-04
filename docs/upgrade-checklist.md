# Upgrade Checklist

- `docker-compose.localdev.yml`
  - [`https://hub.docker.com/r/amazon/dynamodb-local`](https://hub.docker.com/r/amazon/dynamodb-local)
- `docker/go/Dockerfile`
  - [`https://hub.docker.com/_/golang`](https://hub.docker.com/_/golang)
  - [`https://github.com/golangci/golangci-lint/releases`](https://github.com/golangci/golangci-lint/releases)
- `docker/web-fn/Dockerfile`
  - [`https://hub.docker.com/_/golang`](https://hub.docker.com/_/golang)
- `templates/base.html`
  - Alpine.js: [`https://cdnjs.com/libraries/alpinejs`](https://cdnjs.com/libraries/alpinejs)
  - Font Awesome [`https://cdnjs.com/libraries/font-awesome`](https://cdnjs.com/libraries/font-awesome)
  - Tailwind CSS [`https://cdnjs.com/libraries/tailwindcss`](https://cdnjs.com/libraries/tailwindcss)
- `terraform/env/terragrunt.hcl`
  - [`https://registry.terraform.io/providers/kreuzwerker/docker/latest/docs`](https://registry.terraform.io/providers/kreuzwerker/docker/latest/docs)
  - [`https://registry.terraform.io/providers/hashicorp/aws/latest/docs`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
  - [`https://formulae.brew.sh/formula/terraform`](https://formulae.brew.sh/formula/terraform)
  - [`https://formulae.brew.sh/formula/terragrunt`](https://formulae.brew.sh/formula/terragrunt)
- `terraform/modules/apigw/cloudfront.tf`
  - [`https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/secure-connections-supported-viewer-protocols-ciphers.html`](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/secure-connections-supported-viewer-protocols-ciphers.html)
