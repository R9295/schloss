## Yarn
  *The recommended package manager for Node.js projects.*

  Version: 1

## Allows ommited fields
   Yes, allowsfollowing omitted fields:
   
 - `integrity`
 
 - when running yarn,  integrity not added to yarn.lock file
 
 - `resolved`

 -  when running yarn,  resolved added to yarn.lock file
   ````
   corepack@^0.14.1:
   version "0.14.1"
   resolved  #allows omission
   integrity #allows omission
   ````

## Allows modified fields

 - `version`
 
 - Changing version of package, causes warning, but it's ignored:
   ````
   corepack@^0.14.1:
   version "0.13.1" #different version
   resolved "https://registry.yarnpkg.com/corepack/-/corepack-0.14.1.tgz#c0a4a3e2cf972a6320ad459bd3a1b55650f35b55"
   integrity sha512-vTxmM8ktrK0OR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7/D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg==
   ````
 - running yarn:
   ````
   warning Lockfile has incorrect entry for "corepack@^0.14.1". Ignoring it.
   ````
 - original version is installed and yarn.lock updateded correctly
  
 - `integrity`
 
 - modifying integrity hash resolves in the following error:
   ````
   error Incorrect integrity when fetching from the cache for "corepack". Cache has "sha512-vTxmM8ktrK0OR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7 
   /D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg== sha1-wKSj4s+XKmMgrUWb06G1VlDzW1U=" and remote has "sha512-vTxmM8ktrK0OsdadR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7
   /D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg==". Run `yarn cache clean` to fix the problem
   ````
 - after cleaning cache:
   ````
   error https://registry.yarnpkg.com/corepack/-/corepack-0.14.1.tgz: Integrity check failed for "corepack" (computed integrity doesn't match our records,got 
   "sha512-vTxmM8ktrK0OR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7/D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg== sha1-wKSj4s+XKmMgrUWb06G1VlDzW1U=")
   ````                                                                                  
 - package is not installed 
  
 - `resolved`
 
 - url change needs valid integrity or no integrity 
   ````
   corepack@^0.14.1:
   version "0.14.1"
   resolved "https://malicious.url
   ````

 - `resolved + integrity`
 
 - changing url and using valid integrity/no integrity is possible. 
 
 - Sample - url points to different package:
   ````
   corepack@^0.14.1:
   version "0.14.1"
   resolved "https://registry.yarnpkg.com/dummy-pkg/-/dummy-pkg-0.0.1.tgz#a30fed3bfbc23db526fe9cec888548bd5ee180a8" #dummy-pkg
   integrity sha512-5LX4qKqvkd3R8noMWqDDy4VQpeM2uJU1OdFDjQhK1S9QB3V2vtjn3FUr8XCiP3zY6zqfFF3uWExed4DPIg81Iw== #dummy-pkg hash
   ````
 - running yarn will turn no errors and dummy-pkg package is installed under corepack package!!
   
## Duplicated field handling  
  
 - `version` `resolved`
 
 - duplicated fields are ignored, the original correct version gets installed and yarn.lock file updated with correct data. 
   ````
   corepack@^0.14.1:
   version "0.14.1"
   version "0.13.1" #duplicate field
   resolved "https://registry.yarnpkg.com/corepack/-/corepack-0.14.1.tgz#c0a4a3e2cf972a6320ad459bd3a1b55650f35b55"
   integrity sha512-vTxmM8ktrK0OR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7/D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg==
   ````
 - `integrity`
 
 - duplicate key precedence issue
   ````
   corepack@^0.14.1:
   version "0.14.1"
   resolved "https://registry.yarnpkg.com/corepack/-/corepack-0.14.1.tgz#c0a4a3e2cf972a6320ad459bd3a1b55650f35b55"
   integrity sha512-vTxmM8ktrK0OR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7/D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg==
   integrity sha512-5LX4qKqvkd3R8noMWqDDy4VQpeM2uJU1OdFDjQhK1S9QB3V2vtjn3FUr8XCiP3zY6zqfFF3uWExed4DPIg81Iw== #duplicate 
   ````
  - results:
   ````
   Incorrect integrity when fetching from the cache for "corepack". Cache has "sha512-vTxmM8ktrK0OR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7/
   D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg== sha1-wKSj4s+XKmMgrUWb06G1VlDzW1U=" and remote has
   "sha512-5LX4qKqvkd3R8noMWqDDy4VQpeM2uJU1OdFDjQhK1S9QB3V2vtjn3FUr8XCiP3zY6zqfFF3uWExed4DPIg81Iw=="
   ````
   
## Duplicate package handling

 - duplicated package are ignored, the original correct version gets installed,yarn.lock file not updated
   ````
   corepack@^0.14.1:
   version "0.14.1"
   resolved "https://registry.yarnpkg.com/corepack/-/corepack-0.14.1.tgz#c0a4a3e2cf972a6320ad459bd3a1b55650f35b55"
   integrity sha512-vTxmM8ktrK0OR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7/D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg==

   corepack@^0.14.1: #duplicate
   version "0.14.1"
   resolved "https://registry.yarnpkg.com/corepack/-/corepack-0.14.1.tgz#c0a4a3e2cf972a6320ad459bd3a1b55650f35b55"
   integrity sha512-vTxmM8ktrK0OR4qZskfCrjiybCRoLBZrx/UVO3srXnkJK3Kf7/D+DU2ZH2kqrev5NAKEbBivTydzWYdF08KrQg==
   ````

