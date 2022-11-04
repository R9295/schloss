# Risk analysis


| Risk | RoE (Result | Risk approximation | Solution |
|------|-------------|--------------------|----------|
| Omitted/Modified/Duplicated  version value | Break functionality | Low | Error on missing/duplicate version value/ notify on modified resolve value
| Omitted/Modified/Duplicate integrity value | Break functionality | Low | Error on missing/duplicate integrity value/ notify on modified resolve value
| Omitted/Modified/Duplicate resolve value   | Break functionality | Low | Error on missing/duplicate resolve value/ notify on modified resolve value
| Omitted/Modified/Duplicated integrity value + Omitted/Modified/Duplicated  resolve url| Install malicious/vulnerable/unintended package | high | Error if both integrity and resolve url of same package have been omitted/ modified/ duplicated
| Duplicate package precendence | Install malicious/vulnerable/unintended package | high| Error on duplicate packages
| Missing checks for package verification by name | Break functionality, install malicious/vulnerable/unintended package | middle | Verify the package name from the downloaded root file (eg. package json)	

To see which package manager is vulnerable to which attack check the Readme files in the individual folders of the package managers.
