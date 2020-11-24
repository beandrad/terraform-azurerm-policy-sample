ARG UBUNTU_VERSION=20.04

FROM ubuntu:${UBUNTU_VERSION}

ARG DEBIAN_FRONTEND=noninteractive
ARG TERRAFORM_VERSION=0.13.0

RUN apt-get update \
    && apt-get -y install git unzip vim azure-cli curl zsh \
    # Install Terraform
    && mkdir -p /tmp/docker-downloads \
    && curl -sSL -o /tmp/docker-downloads/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && unzip /tmp/docker-downloads/terraform.zip \
    && mv terraform /usr/local/bin \
    && rm -rf /tmp/docker-downloads \
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

ARG DEBIAN_FRONTEND=interactive
