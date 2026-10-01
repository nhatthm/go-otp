# OTP

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/go-otp)](https://github.com/nhatthm/go-otp/releases/latest)
[![Build Status](https://github.com/nhatthm/go-otp/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/go-otp/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/go-otp/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/go-otp)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/otp)
[![Donate](https://img.shields.io/badge/%20-Donate-%20?style=flat&logo=githubsponsors&color=E5E4E2)](http://donate.nhat.me)

Manage and generate one-time passwords.

## Prerequisites

- `Go >= 1.25`

### Keyring (Optional)

Support **OS X**, **Linux/BSD (dbus)** and **Windows**.

#### OS X

The OS X implementation depends on the `/usr/bin/security` binary for
interfacing with the OS X keychain. It should be available by default.

#### Linux and *BSD

The Linux and *BSD implementation depends on the [Secret Service][SecretService] dbus
interface, which is provided by [GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring).

It's expected that the default collection `login` exists in the keyring, because
it's the default in most distros. If it doesn't exist, you can create it through the
keyring frontend program [Seahorse](https://wiki.gnome.org/Apps/Seahorse):

* Open `seahorse`
* Go to **File > New > Password Keyring**
* Click **Continue**
* When asked for a name, use: **login**

## Install

```bash
go get go.nhat.io/otp
```

## Usage

Example 1: Generate a new TOTP using a provided secret.

```go
package main

import (
	"context"

	"go.nhat.io/otp"
)

func do(ctx context.Context) {
	result, err := otp.GenerateTOTP(ctx, otp.TOTPSecret("NBSWY3DP"))

	if err != nil {
		// Handle error.
	}

	// Use the result.
}
```

Example 2: Generate a new TOTP using a secret that persisted in keychain.

```go
package main

import (
    "context"

    "go.nhat.io/otp"
    "go.nhat.io/otp/keyring"
)

func do(ctx context.Context) {
    result, err := otp.GenerateTOTP(ctx, keyring.TOTPSecretFromKeyring("john.doe@example.com"))

    if err != nil {
        // Handle error.
    }

    // Use the result.
}
```

## Donation

If this project saved you some development time, buy me a cup of coffee :)

[![donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](http://donate.nhat.me)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://github.com/nhatthm/donate.nhat.me/blob/master/images/qr_sponsor.png" width="147px" />

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)
