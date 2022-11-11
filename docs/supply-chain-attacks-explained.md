# Supply Chain Attacks

Supply chain attacks are attacks that seek to damage a pice of software by targeting the less secure elements in the supply chain. When talking about a software supply chain that includes:
- CI/CD Pipelines
- Containers
- Package registries
- Container registries
- Container orchestrators
- Dev machines
- Vendors
- Hosting environments
- “The Cloud"
  
Basically anything that goes into or affects your code from developmet, through your CI/CD pipeline until it gets deployed into production.
=> Very difficult to secure due to breadth and complexity

Round about 78 % of a pice of software is not written by the original creator but comes from open source dependencies or language libraries.
On average a repository on github has [dependencies to 203](https://github.blog/2020-09-02-secure-your-software-supply-chain-and-protect-against-supply-chain-threats-github-blog/) repositories.

### Problems
- Almost impossible to audit all the code/ dependencies
- Managing vulnerabilities in dependencies is cumbersome
- Updating dependencies means continuously bringing in code that hasn't been evaluated yet.

## Big attacks in the past
- Solar Winds
- event-stream => malicious update (1.9 Million weekly downlaods)
- Poisoned PyPi packages
- Kesaya

- Typosquatting => colour-pickr
- Dependency Confusion => paypal, alex birsan

## How to protect oneself?
- Purchasing ALL potential typo'd domains
  - For less than 10$/year per deomain typo attacks can be prevented
  - Redirect: Make sure all domains point to the same main domain
- Register namespaces on public dependency managers
  - Register the names and similar spellings of internal libraries and packages being used by the organization on public repositories. This would prevent attacker from tricking a developer into installing an identically or similarly named package on a public repository.
- Activate namespaced modules
  - Many package managers support namedspaced modules, which prevents the same name from being used for two different resources
- Treat lists of internal packages as sensitive information
  - The success of these attacks depend on knowing the name of internal resources, organizations can no longer be casual about where this information is stored and shared
- Use lockfiles
  - ensures that you get the same package version on every install, so if you are secure today then you are secure tomorrow
  - **Make sure to review all lockfile changes very carefully and use an aditional linting tool**
- Review every new dependency
  - Have a policy on what to look out for when considering adding a new dependency to the project, e.g:
    - Do you need to install this dependency?
    - Is the package well maintained?
    - Does the package the source match the package manager?
    - Can you write it yourself?
    - Should you clone and host yourself? (still watch out for security updates)
- Use Subresource Integrity ([SRI](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity)) with CDNs
  - Direct the browser to only pull a specific resource if the hash of the file matches the cryptographic hash you have specified
- Continuously monitor dependencies
  - Have a update dependency policy 
  - Use tools to be alerted when dependencies have a security update
- Prefer signed packages
- Protect your deployment accounts
  - 2-factor-authentication
  - role based authentication on who can puss artifacts, etc