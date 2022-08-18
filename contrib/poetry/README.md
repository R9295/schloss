## Poetry
*Python packaging and dependency management made easy.*

Version: ``1.1.14``

### Allows ommited fields
No. All of the fields in a package entry must be present
```toml
[[package]]
name = "asgiref"
version = "3.5.2"
description = "ASGI specs, helper code, and adapters"
category = "main"
optional = false
python-versions = ">=3.7"
```
### Duplicate package field handling
Poetry errors on duplicate fields
```toml
[[package]]
name = "asgiref"
version = "3.5.2"
version = "3.5.4" # duplicate
description = "ASGI specs, helper code, and adapters"
category = "main"
optional = false
python-versions = ">=3.7"
```
Will throw:
```
RuntimeError

Unable to read the lock file (Invalid TOML file /somewhere/poetry.lock: Key "version" already exists.).
```

### Duplicate [metadata.files] handling
The ``metadata.files`` section is for storing file(tar/wheel/misc) file hashes.

The following
``` toml
[metadata.files]
asgiref = [
    {file = "asgiref-3.5.2-py3-none-any.whl", hash = "sha256:...."},
    {file = "asgiref-3.5.2.tar.gz", hash = "sha256:...."},
]
asgiref = [
    {file = "asgiref-3.5.2-py3-none-any.whl", hash = "sha256:....."},
    {file = "asgiref-3.5.2.tar.gz", hash = "sha256:...."},
] # duplicate entry
```
Would throw:
```
RuntimeError

Unable to read the lock file (Invalid TOML file /somewhere/poetry.lock: Key "asgiref" already exists.).
```

### Duplicate package handling
Poetry **does not error** on duplicate package entry!
It will **override** the previous package entry!
``` toml
[[package]]
name = "asgiref"
version = "3.5.2"
description = "ASGI specs, helper code, and adapters"
category = "main"
optional = false
python-versions = ">=3.7"

# duplicate
[[package]]
name = "asgiref"
version = "3.5.1"
description = "ASGI specs, helper code, and adapters"
category = "main"
optional = false
python-versions = ">=3.7"
```

### Hash handling
The values of entries in ``[metadata.files]`` appear to not be mandatory. For example, when installing Django,
the list is empty.

However, each package must have an entry in ``[metadata.files]``, even if it is an empty list.
Otherwise, it would throw:
```
NonExistentKey

'Key "<some_key>" does not exist.'
```

One particularly interesting thing:

``` toml
[metadata.files]
asgiref = [
    # the hash for the tar file is set as the hash for the wheel file and vice versa
    {file = "asgiref-3.5.2-py3-none-any.whl", hash = "sha256:HASH_FOR_TAR_FILE"},
    {file = "asgiref-3.5.2.tar.gz", hash = "sha256:HASH_FOR_WHEEL_FILE"},
]
```
The hashes have been swapped **BUT** it would still work as long as the right hash
is in the list, be it under any filename!

Since this is the case, duplicates do not affect anything
``` toml
[metadata.files]
asgiref = [
    # the hash for the tar file is set as the hash for the wheel file and vice versa
    {file = "asgiref-3.5.2-py3-none-any.whl", hash = "sha256:HASH_FOR_TAR_FILE"},
    {file = "asgiref-3.5.2.tar.gz", hash = "sha256:HASH_FOR_WHEEL_FILE"},
    {file = "asgiref-3.5.2.tar.gz", hash = "sha256:DUPLICATE_ENTRY"}, # duplicate
]
```
But this is of major concern!
