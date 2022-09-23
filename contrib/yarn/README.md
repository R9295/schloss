## Yarn
  *The recommended package manager for Node.js projects.*

  Version: 1

## Allows ommited fields

   Yes, `integrity` and `resolved`

   ````
   corepack@^0.14.1:
   version "0.14.1"
   resolved  # allows omission
   integrity # allows omission
   ````
  Note: In case of omission of ``resolved`` and ``integrity``, running yarn will NOT populate ``integrity`` but will populate``resolved``.

## Duplicated field handling  
  Yarn has a duplicate field issue for fields ``resolved`` and ``integrity``, but not for ``version``.
```
   corepack@^0.14.1:
   version "0.14.1"
   resolved  https://something
   resolved  https://something.duplicate # duplicate field
   integrity sha512-xyz
   integrity sha512-xyz-duplicate # duplicate field
```
## Duplicate package handling

Duplicated package entries are ignored, the first entry gets installed.
   ````
   corepack@^0.14.1:
   version "0.14.1"
   resolved "https://registry.yarnpkg.com/corepack/-/corepack-0.14.1.tgz#c0a4a3e2cf972a6320ad459bd3a1b55650f35b55"
   integrity sha512-vTxmM8ktrK0OR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7/D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg==

   corepack@^0.14.1: # duplicate
   version "0.14.1"
   resolved "https://registry.yarnpkg.com/corepack/-/corepack-0.14.1.tgz#Duplicate"
   integrity sha512-Duplicate
   ````
   Note: yarn does not remove the duplicate entry from the lockfile.

## Swapping Package Fields
Swapping a package's ``resolved`` and ``integrity`` fields with another's will install the package as the other. For example:
````
   pkgA@^0.5.1:
   version "0.5.1"
   resolved "https://registry.yarnpkg.com/pkgB"  # swapped with pkgB
   integrity sha512-pkgB # swapped with pkgB

   pkgB@^0.14.1:
   version "0.14.1"
   resolved "https://registry.yarnpkg.com/pkgA" # swapped with pkgA
   integrity sha512-pkgA  # swapped with pkgA
````
This will install ``pkgB`` as ``pkgA`` and vice versa!
