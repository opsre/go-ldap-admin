# Contribution Guide

We welcome feedback, bug reports, and pull requests. Please submit them through [issues](https://github.com/eryajf/go-ldap-admin/issues).

If this is your first time collaborating on GitHub, see the [Collaborative Development Process](https://howtosos.eryajf.net/HowToStartOpenSource/01-basic-content/03-collaborative-development-process.html).

1. We use `golangci-lint` for checks. Before submitting a PR, run `make lint` locally and make sure it passes.

2. Unless necessary, add configuration files sparingly to avoid unexpected issues during upgrades.

3. When a feature change involves frontend updates, link the two PRs whenever possible; otherwise, merging can be difficult.
