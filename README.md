# Azure Verified Modules TFLint Ruleset

[![Build Status](https://github.com/Azure/tflint-ruleset-avm/workflows/build/badge.svg?branch=main)](https://github.com/Azure/tflint-ruleset-avm/actions)

This repository contains the TFLint ruleset for Azure Verified Modules.

## Requirements

- TFLint v0.42+
- Go v1.22

## Installation

TODO: This template repository does not contain release binaries, so this installation will not work. Please rewrite for your repository. See the "Building the plugin" section to get this template ruleset working.

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "avm" {
  enabled = true

  version = "0.2.0"
  source  = "github.com/Azure/tflint-ruleset-avm"

  signing_key = <<-KEY
----BEGIN PGP PUBLIC KEY BLOCK-----
Version: BSN Pgp v1.1.0.0

mQENBF9hII8BCADEOCDl3/1tAZQp/1BCVJN+tqIRCd3ywzhOXTC38XWC0zVbFtiA
vbBFL1e78aoDIyUFDZcphCyYDqBkweXeYyYVCojZFVniyKklc2xZ15LDwlMBhneU
yEPSzDCltFn67wMPQMKa4+TujZJ3TIs1OUnUTsCPrjavGgmrfAdxAF/EjCDrnVp9
XmRWJii/9elAnMqWLDkMDfPaWkv3lWuyYCBHc7avOJE9oWypmWoEPOujwmtika/i
FhmvZbojZN6huf7pykXGRl1wEpu0MMEFvm4UsfEOv8JHVBZEu2w6glQugT6a+IZ6
atH3zyy+i1mmgsJPlMF1soHNEufeK1CabMklABEBAAG0Q1RlcnJhZm9ybSBBRE8g
cHJvdmlkZXIgcmVsZWFzZSA8dGVycmFmb3JtYWRvcHJvdmlkZXJAbWljcm9zb2Z0
LmNvbT6JATgEEwEIACIFAl9hII8CGwMGCwkIBwMCBhUIAgkKCwQWAgMBAh4BAheA
AAoJEG8Lkb3phHjPT+YH/3aksw2yhoqVl+Dxkrpsq9LIsXBHmHfbk8/nwbZ7F6o6
fZetwozQzS/v5IriE42NFdk2omilDa/Iumk5soPrCamIIToYMbGvZJ9MJzCflXzp
H3crqEgoCwu/93FVot4hhNOGmS2ra538zDQ3JsSbsVSc2TyPeBCF08+qJrr9VSML
LceuEvCKUN8P8LH+PXN4kKM1xNlSVw4RfH6mNJKdUG1Klvh2nbq0kuw8jiHITn2F
ALGvKXPLwggdNA86RIQc9tc3z/uJrBGSA2n6UkJbV1gFZDETjHzVtgDqqEQwap7D
/i9e5KqIAEIf14OPm3h+e6kCdWXRG0RJWWVWeOHIEfQ=
=KwXd
-----END PGP PUBLIC KEY BLOCK-----
  KEY
}
```

## Rules

|Name|Description|Severity|Enabled|Link|
| --- | --- | --- | --- | --- |

> *TBC*

## Building the plugin

Clone the repository locally and run the following command:

```bash
make
```

You can easily install the built plugin with the following:

```bash
make install
```

You can run the built plugin like the following:

```bash
$ cat << EOS > .tflint.hcl
plugin "avm" {
  enabled = true
}
EOS
$ tflint
```
