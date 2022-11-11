### Todo
- https://www.youtube.com/watch?v=fCaglk1bpmI
- prepare a simple demo
- finish first minimum version for npm
- check also for version change in subdependencies
- writing tests
- further documentation + risk analysis
- prepare for assessment


integrity => sha hash (function that takes an input and returns an output of a specifc defined length, alsways same length)
An other security feature to check if url is valid, and website not attacked
Browser gets url and checks source of script and gets a hash back. Browser validates if returned hash is same as integrity hash. If yes then script can be downloaded from source and executed, if not, Script / website is malicious
??? If we are able to change url then we would also be able to change hash so not adding that much value ...? => more obvious to change url, human readable
??? When browser is evaluating if hashes match, browser already has a connection to the website, possible to run attacks based on that connection?
What happens if hash is removed?
How do we know that first hash has not been modified?

Now we only check for packages and dependencies of the packages but what about dependencies of the dependencies...?



transitive dependencies => dependencies of dependencies 

## Open Questions
- Why are npm peerDependencies installed by default?
- What happens now if different transitive dependencies have same peerDependency but different versions?
- What is a bundled dependency?
- Is there a way to check if the resolved and integrity value correspond to the official one? With schloss we can easily detect changes to fields of existing packages but if someone adds a new package and replaces the resolved and integrity values with some of a malicious package is there a way for us to detect such a change?
E.g. create from lockfile package a package.json and use precise versions, create a lockfile from the created package.json and compare hash and integrity value
with version check if integrity and resolved value correspond with the official sources

## Answered Questions
- How does integrity sha tag exactely work? How to make sure, that very first integrity value has not been modified? => very first time installing you never can be sure 100 % that's why being able to remove the integrity hash is a security vulnerabilit

## Todo
- messing with subdependencies
- messing with other fields
- checking how to downgrade minor version in package-lock.josn => I thought running npm ci changes would get applied but they don't


risk analysis
=> what is the risk, why is it bad, what can you do to omittigate it...
semantic diff algorithm
direct integration in git? => we could only show diffs that change functionality
=> are we talking now about 2 different tools
linter to throw errors for invalid structure
risk evaluation
*report security vulnerabilities to npm*

git checkout HEAD testdata/

### Todo
- first min version pull request review
- subdependencies of subdependencies in lockfile ignored
- write test
- demo
- read up, on next level supply chain attacks + rewatch talk
- risk analysis + work on documentation
- document code?
- issue if I simplified something


### Lockfile tampering as a sort of supply chain attacks
https://snyk.io/blog/why-npm-lockfiles-can-be-a-security-blindspot-for-injecting-malicious-modules/
https://www.youtube.com/watch?v=fCaglk1bpmI

lockfiles can have a huge impact in the supply chain, a changes in a lockfile may influnce the CI/CD pipeline, contaners, dev machines, vendors, ...

Editing lockfiles manaually is not prevented by package managers

### Attacker Perspective: Compromising Supply Chains using Lockfiles
Goal is to inject controlled dependency into supply chain without breaking functionality, requires a lot of knowledge about the software + techstack + build systems

Exploit Behavior
- Compromise dev machines
- Compromise CI, hosting environment (kubernetes, vms, databases, etc.)
Upstream injection
- A dependency upstream of a direct dependency is exploited
Direct injection
- A direct dependency of a project is tampered
Indirect injection
- A dependency of a dependency is tampered
Other flavors of this attack
- Anything in the software supply chain where lockfiles are present can be a target

E.g. somebody adds a dependency, lockfile diffs get very large very fast, quite easy to e.g. change resolved + hash in between of an other dependency so that it quite hard for the reviewer to detect them, with schloss it would print the changes in a more readable way "dependency xyz has been added, integrity hash of dependency z has been modified,..."

Accepted best practice: commit your lockfiles to source control
- Adds a layer of difficulty to lockfile tampering attack
- Lockfiles are huge!
- When was the last time your devs read through lockfile changes?


### Lockfile tampering exampe: 
updating version of package and changing div, resolved value in lock file, schloss could be added in the ci/cd pipeline and raise an error if fields of existing lockfiles have been modified

Wha happens mostly is that people check the changes in package.json thouroughly
    The packages I introduced are known in the ecosystem and are vulnerabilities free – check.
    There are no typosquatting attempts in these package names – check.
    These are valid versions of those packages and aren’t malicious in and of themselves – check.
But the package.lock file is often not checked thouroughly. Since changes can be pretty huge very fast just by adding one dependency it is also easy to overlook a change in an other package


### similar projects
- https://github.com/lirantal/lockfile-lint
- https://gitlab.com/gitlab-org/frontend/untamper-my-lockfile


https://www.atlanticcouncil.org/in-depth-research-reports/report/breaking-trust-shades-of-crisis-across-an-insecure-software-supply-chain/

https://github.blog/2020-09-02-secure-your-software-supply-chain-and-protect-against-supply-chain-threats-github-blog/



 ~/go/bin/golines -w *