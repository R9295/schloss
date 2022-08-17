## Cargo
*The recommended package manager for Rust projects.*

Version: 3

### Allows ommited fields
Yes, ``source`` ``dependencies`` and ``checksum``
```toml
[[package]]
name = "autocfg"
version = "1.1.0"
source = "xyz" # allows ommission
checksum = "xyz" # allows ommmission
dependencies = [] # optional
```
We only care about ``source`` and ``checksum`` being ommited

### Duplicate field handling
Cargo errors on duplicate fields
```toml
[[package]]
name = "autocfg"
version = "1.1.0"
source = "xyz"
checksum = "xyz"
checksum = "asdf"
```
Will throw:
```
error: failed to parse lock file at: /somewhere/Cargo.lock

Caused by:
  could not parse input as TOML

Caused by:
  TOML parse error at line 29, column 1
     |
  29 | checksum = "asdf"
     | ^
  Duplicate key `checksum` in table `package`
```

### Duplicate package handling
Cargo errors on duplicate packages
``` toml
[[package]]
name = "autocfg"
version = "1.1.0"
source = "registry+https://github.com/rust-lang/crates.io-index"
checksum = "d468802bab17cbc0cc575e9b053f41e72aa36bfa6b7f55e3529ffa43161b97fa"

[[package]]
name = "autocfg"
version = "1.1.0"
source = "registry+https://malicious"
checksum = "malicious"
```
Will throw:
```
error: failed to parse lock file at: /somewhere/Cargo.lock

Caused by:
  package `autocfg` is specified twice in the lockfile
```
