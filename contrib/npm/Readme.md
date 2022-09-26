# NPM
*Common package manager for node.js projects.*

Version: 8.5.0


## Allows ommited fields
When deleting the version of a package in package-lock.json npm does not complain, as long as it is valid json.
````
"node_modules/@angular/core": {
      "version": "14.2.2", # allows omitted
      "resolved": "https://registry.npmjs.org/@angular/core/-/core-14.2.2.tgz", # allows omitted
      "integrity": "sha512-kG30b4RqjgWvaH9y4g95JRCzoROV+9/xgFH4hSRejFa//Yvw==", # allows omitted
    },
````
-  ``version`` 
   -  When running `npm ci` it checks package.json for version, and installs the latest
   -  When running `npm i` version NOT added again => once missing always missing 
- ``integrity`` 
  - When running `npm ci` npm does NOT complain, new integrity hash also does not get added.
  - After running `npm i` integrity hash does NOT get added again => once missing always missing 
- ``resolved``
  - When running `npm ci` npm does NOT complain, new resolved value also does not get added.
  - After running `npm i` resolved value does NOT get added again => once missing always missing 

## Allows modified fields

- `version`
  - Changing version of package, so that it is not compatible with the version defined in package.json:
    ````
    package.json:
    "@angular/forms": "^14.2.0",

    package-lock.json
    "node_modules/@angular/forms": {
      "version": "13.2.2", # old version "14.2.2"
      "integrity": "sha512-9h8MwFLvIJ5kB5L03cd3Cyl4ySKVzL/E/YYugfLvcAzYZ8Rief63gJnkcKN+YxDOw==",
    },
    ````
  
    Running `npm ci` will throw the following error:
    ````
    npm ERR! `npm ci` can only install packages when your package.json and package-lock.json or npm-shrinkwrap.json are in sync. Please update your lock file with `npm install` before continuing.
    npm ERR! 
    npm ERR! Invalid: lock file's @angular/forms@13.2.2 does not satisfy @angular/forms@14.2.2
    npm ERR! 

    npm ERR! A complete log of this run can be found in:
    npm ERR!     /
    ````
  - Changing version of package that is still compatible with the version defined in package.json is possible, modified package version is gonna be installed:
    ````
    package.json:
    "@angular/forms": "^14.2.0",

    package-lock.json
    "node_modules/@angular/forms": {
      "version": "14.2.1", # old version "14.2.2"
      "integrity": "sha512-9h8MwFLvIJ5kB5L03cd3Cyl4ySKVzL/E/YYugfLvcAzYZ8Rief63gJnkcKN+YxDOw==",
    },
    ````
- `integrity`
  - modifying integrity hash
    ````
    package.json:
    "@angular/forms": "^14.2.0",

    package-lock.json
    "node_modules/@angular/forms": {
      "version": "14.2.2", 
      "integrity": "sha512-sdfsdf/E/YYugfLvcAzYZ8Rief63gJnkcKN+YxDOw==", # modified
    },
    ````
    Running `npm ci`/ `npm i` resolves in the following error until has is correctly updated again:
    ````
    npm ERR! code EINTEGRITY
    npm ERR! sha512-30b4RqjgWvaH9y4g95JRCzoROV+9/xgFH4hSRejFa/IcapMfvCmONJtJzwTjdsEUQAbiFohF/z9bx3QA/Yvw== integrity checksum failed when using sha512: wanted sha512-30b4RqjgWvaH9y4g95JRCzoROV+9/xgFH4hSRejFa/IcapMfvCmONJtJzwTjdsEUQAbiFohF/z9bx3QA/Yvw== but got sha512-kG30b4RqjgWvaH9y4g95JRCzoROV+9/xgFH4hSRejFa/IcapMfvCmONJtJzwTjdsEUQAbiFohF/z9bx3QA/Yvw==. (4407879 bytes)

    npm ERR! A complete log of this run can be found in:
    npm ERR!     /
    ````
- `resolved`
  - changing to a non-valid url may or may not cause an error:
    ````
     "node_modules/accepts": {
      "version": "1.3.8",
      "resolved": "https://registry.npmjs.dsdfsfdsdfsdforg/accepts/-/accepts-1.3.8.tgz", # modified
      "integrity": "sha512-PYAthTa2m2VKxuvSD3DPC/Gy+U+sOA1LAuT8mkmRuvw+NACSaeXEQ+NHcVF7rONl6qcaxV3Uuemwawk+7+SJLw==",
    },
    ```
    Even though the url is not valid npm does not complain and installs the right package from npm. What if someone would mess with the resolved url of a package in a private repository? Would npm per default search in npm for for a valid package?

    ````
      "node_modules/angular": {
        "version": "1.8.3",
        "resolved": "https://registry.npmjs.org/angulsafdasfdasfdaar/-/angular-1.8.3.tgz", # modified
        "integrity": "sha512-5qjkWIQQVsHj4Sb5TcEs4WZWpFedasdfasdfasfdasdfasdfVFHXwxEBHUhrny41D8UrBAd6T/6nPPAsLngJCReIOqi95W3mxdveveutpZw==",
      },
    ```
    Running `npm ci` will throw the following error:
    ````
      npm ERR! code E404
      npm ERR! 404 Not Found - GET https://registry.npmjs.org/sdfasfdasfdg/angular/-/angular-1.8.3.tgz
      npm ERR! 404 
      npm ERR! 404  'angular@https://registry.npmjs.org/sdfasfdasfdg/angular/-/angular-1.8.3.tgz' is not in this registry.
      npm ERR! 404 You should bug the author to publish it (or use the name yourself!)
      npm ERR! 404 
      npm ERR! 404 Note that you can also install from a
      npm ERR! 404 tarball, folder, http url, or git url.
    ````

  - Changing resolved value to an other, valid url points to different version
    points to different version:
    ````
      "node_modules/@angular/common": {
        "version": "14.2.2",
        "resolved": "https://registry.npmjs.org/@angular/common/-/common-14.2.1.tgz", # originally version 14.2.2
        "integrity": "sha512-9h8MwFLvIJ5kB5L03cd3Cyl4ySKVzL/E/YYugfLvcAzYZ8Rief63gJnkcKNjoS1A5DTxHhOBQL7pLZpj+YxDOw==",
      }
    ````
    Npm does not throw an error or corrects the modified field but still manages to install the correct package with the original version 14.2.2. How? What happens if executed on a cicd pipeline?

    points to different package:
     ````
      "node_modules/@angular/common": {
        "version": "14.2.2",
        "resolved": "https://registry.npmjs.org/@types/whatwg-url/-/whatwg-url-8.2.2.tgz",, # points to diffferent package
        "integrity": "sha512-9h8MwFLvIJ5kB5L03cd3Cyl4ySKVzL/E/YYugfLvcAzYZ8Rief63gJnkcKNjoS1A5DTxHhOBQL7pLZpj+YxDOw==",
      }
    ````
    Npm does not throw an error or corrects the modified field but still manages to install the original package with the correct version. How? What happens if executed on a cicd pipeline?

## Switching integrity and resolved value of two packages

## Duplicated field handling
- duplicated fields are ignored, looks like file does not actually get updated, even if field is duplicated and both versions are changed, the original correct version gets installed



## Duplicate package handling
- duplicated package are ignored, file does not get updated, even if package is duplicated and both versions are changed, the original correct version gets installed
  ````
    package.json
    "@angular/forms": "^14.2.1",

    "node_modules/@angular/forms": {
      "version": "14.2.2",
      "resolved": "https://registry.npmjs.org/@angular/forms/-/forms-14.2.2.tgz",
      "integrity": "sha512-lq6PpLMNfs0KcIhkQrQZKh6XHsbQRcuhWKHy6IITldz6sg8pWVrTw==",
    },
    "node_modules/@angular/forms": {
      "version": "14.2.0",
      "resolved": "https://registry.npmjs.org/@angular/forms/-/forms-14.2.2.tgz",
      "integrity": "sha512-lq6PpLMPhDiakvOWYoQZKh6XHsbQRcuhWKHy6IITldz6sg8pWVrTw==",
    },
  ````
  version 14.2.2 gets installed


  Both versions have been modified:
    ````
    package.json
    "@angular/forms": "^14.2.1", 

    "node_modules/@angular/forms": {
      "version": "14.2.0", # modified
      "resolved": "https://registry.npmjs.org/@angular/forms/-/forms-14.2.2.tgz",
      "integrity": "sha512-lq6PpLMNfs0KcIhkQrQZKh6XHsbQRcuhWKHy6IITldz6sg8pWVrTw==",
    },
    "node_modules/@angular/forms": {
      "version": "14.2.1",  # modified
      "resolved": "https://registry.npmjs.org/@angular/forms/-/forms-14.2.2.tgz",
      "integrity": "sha512-lq6PpLMPhDiakvOWYoQZKh6XHsbQRcuhWKHy6IITldz6sg8pWVrTw==",
    },
  ````
  version 14.2.2 gets installed

## missing: check other fields for vulnerabilities

## npm i vs. npm ci

``npm install``/ ``npm i``
- Installs all dependencies defined in package.json
- If you use [version ranges](https://github.com/npm/node-semver#ranges) there different versions may be installed depending at what time you run npm install.
- npm install can update/ create your package-lock.json when there are changes such as when you install or remove a dependency.

``npm clean install``/ ``npm ci``

- Deletes node_modules folder to ensure a clean state
- Installs all dependencies with the exact version defined in package-lock.json
- Unlike npm install, npm ci will never modify your package-lock.json. It does however expect a package-lock.json file in your project — if you do not have this file, npm ci will not work and you have to use npm install instead.
- If dependencies in the package lock do not match those in package.json, npm ci will exit with an error, instead of updating the package lock.
  
## Structure package.json
Simple example:
````
{
  "name": "schloss-npm",
  "version": "1.0.0",
  "description": "- find package manager vulnerabilities + document them with tests and doc - get familiar with go - write schloss just for npm and without paying attention to the subdependencies - understand current code base (fully)",
  "main": "index.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "author": "",
  "license": "ISC",
  "dependencies": {
    "@angular/forms": "^14.2.0",
    "datetime": "^0.0.3",
    "forms-angular": "^0.12.0-beta.191"
  },
  "peerDependencies": {},
  "devDependencies": {}
}

````
- dependencies
  - Dependencies are specified in a simple object that maps a package name to a version range. 
  - Dependencies can also be identified with a tarball or git URL.
- devDependencies
  - Only for development, not installed during build
- peerDependencies
  - In some cases, you want to express the compatibility of your package with a host tool or library, while not necessarily doing a require of this host. This is usually referred to as a plugin. Notably, your module may be exposing a specific interface, expected and specified by the host documentation.
  - From v7 also installed by default
  - Trying to install another plugin with a conflicting requirement may cause an error if the tree cannot be resolved correctly. 
- peerDependenciesMeta

More detailed [documentation](https://docs.npmjs.com/cli/v8/configuring-npm/package-json#dependencies) of the package.json.


## Structure package-lock.json
- Automatically generated for any operations where npm modifies either the node_modules tree, or package.json. 
- Describes the exact tree that was generated, to be able to recreate identical trees, regardless of intermediate dependency updates.
- package-lock.json has to be in the root folder, otherwise it's ignored, can be overwritten by adding a npm-shrinkwrap.json (exact same filestructure)
````
{
  "name": "schloss-npm",
  "version": "1.0.0",
  "lockfileVersion": 2,
  "requires": true,
  "packages": {
    "version": "14.2.0",
      "resolved": 
      "integrity": 
      "peer": true,
      "dependencies": {}, => subdependencies
      "engines": {},
      "peerDependencies": {}
  }
  "dependencies": {
    
  }
}
````
- packages
  - This is an object that maps package locations to an object containing the information about that package.
  - The root project is typically listed with a key of "", and all other packages are listed with their relative paths from the root project folder.
  
  Package descriptors have the following fields:

  - `version`: The version found in package.json
  - `resolved`: The place where the package was actually resolved from. In the case of packages fetched from the registry, this will be a url to a tarball. In the case of git dependencies, this will be the full git url with commit sha. In the case of link dependencies, this will be the location of the link target. registry.npmjs.org is a magic value meaning "the currently configured registry".
  - `integrity`: A sha512 or sha1 Standard Subresource Integrity string for the artifact that was unpacked in this location.
  - `link`: A flag to indicate that this is a symbolic link. If this is present, no other fields are specified, since the link target will also be included in the lockfile.
  - `dev, optional, devOptional`: If the package is strictly part of the devDependencies tree, then dev will be true. If it is strictly part of the optionalDependencies tree, then optional will be set. If it is both a dev dependency and an optional dependency of a non-dev dependency, then devOptional will be set. (An optional dependency of a dev dependency will have both dev and optional set.)
  - `inBundle`: A flag to indicate that the package is a bundled dependency.
  - `hasInstallScript`: A flag to indicate that the package has a preinstall, install, or postinstall script.
  - `hasShrinkwrap`: A flag to indicate that the package has an npm-shrinkwrap.json file.
  - `bin, license, engines, dependencies, optionalDependencies`: fields from package.json
- dependencies
  - Legacy data for supporting versions of npm that use lockfileVersion: 1. This is a mapping of package names to dependency objects. Because the object structure is strictly hierarchical, symbolic link dependencies are somewhat challenging to represent in some cases.
  - npm v7 ignores this section entirely if a packages section is present, but does keep it up to date in order to support switching between npm v6 and npm v7.

More detailed [documentation](https://docs.npmjs.com/cli/v8/configuring-npm/package-lock-json) of the package-lock.json.