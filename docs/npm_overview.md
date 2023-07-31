# NPM 

NPM is the world's largest software registry. Open source developers from every continent use the registry and the CLI to share and borrow packages. 

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

npm ci is the command of choice when being used in a CI/CD pipeline. The problem is, that the package manager does not verify the integrity of a lockfile and when not checking carefully every character in a pull request it is easy to tempter the lockfiles. Read the docs to find out in detail how the lockfile can be tempered.
  
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