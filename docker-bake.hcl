group "default" {
    targets = [
        "server-latest",
        ]
}

variable "REGISTRY" {
    default = ""
}

variable "IMAGE_NAME" {
    default = "obalunenko/shipping-pack-optimizer-server"
}

variable "IMAGE_WITH_REGISTRY" {
    default = notequal("",REGISTRY) ? "${REGISTRY}/${IMAGE_NAME}": "${IMAGE_NAME}"
}

variable "IMAGE_TITLE" {
    default = "shipping-pack-optimizer"
}

variable "CI_COMMIT_TAG" {
    default = ""
}

variable "CI_COMMIT_REF_SLUG" {
    default = ""
}

variable "CI_COMMIT_SHORT_SHA" {
    default = ""
}

variable "IMAGE_TAG" {
    default = notequal("", CI_COMMIT_TAG) ? "${CI_COMMIT_TAG}" : "${CI_COMMIT_REF_SLUG}-${CI_COMMIT_SHORT_SHA}"
}

variable "BUILD_TAG" {
    default = notequal("-", IMAGE_TAG) ? "${IMAGE_WITH_REGISTRY}:${IMAGE_TAG}" : "${IMAGE_WITH_REGISTRY}:latest"
}

variable "IMAGE_DESCRIPTION" {
    default = ""
}

target "docker-metadata-action" {}

target "server-latest" {
    inherits = ["docker-metadata-action"]
    dockerfile = "Dockerfile"
    context    = "."
    platforms = [
        "linux/amd64",
        "linux/arm64"
    ]
    labels = {
        "org.opencontainers.image.title"       = "${IMAGE_TITLE}"
        "org.opencontainers.image.description" = "${IMAGE_DESCRIPTION}"
    }
    tags       = [
        "${BUILD_TAG}"
    ]
}