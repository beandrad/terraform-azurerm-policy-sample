ARG UBUNTU_VERSION=20.04

FROM ubuntu:${UBUNTU_VERSION}

ARG DEBIAN_FRONTEND=noninteractive
ARG TERRAFORM_VERSION=0.13.0
ARG GO_VERSION=1.15.1

RUN apt-get update \
    && apt-get -y install git unzip vim azure-cli curl zsh \
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /tmp/docker-downloads

# Install Terraform
RUN curl -sSL -o /tmp/docker-downloads/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip \
    && unzip /tmp/docker-downloads/terraform.zip \
    && mv terraform /usr/local/bin

# Install Go
RUN curl -fsS https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz -o /tmp/docker-downloads/golang.tar.gz \
    && tar -C /usr/local -xvf /tmp/docker-downloads/golang.tar.gz

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV PATH=$PATH:/usr/local/go/bin

RUN rm -rf /tmp/docker-downloads

ARG DEBIAN_FRONTEND=interactive
