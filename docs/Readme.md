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
    Running `npm ci`/ `npm i` resolves in the following error until hash is correctly updated again:
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
## Adding package that is not in package.json
When adding a dependency to the package-lock.json file that is not a subdependency or a dependency in the package.json it does not get installed.