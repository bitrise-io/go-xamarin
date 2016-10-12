## Changelog (Current version: 1.0.0)

-----------------

## 1.0.0 (2016 Oct 12)

### Release Notes

* verbose mode, outputmap fix
* build tools package
* dsym search fix
* warning handling
* fixed errors
* edit build command
* mdtool timeout handling
* duplicated commands fix
* analyze target definition fix
* builder tests
* validate solution config before build
* copy command in diagnostic mode, instead of using directly

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-tools/go-xamarin/releases/download/1.0.0/go-xamarin-$(uname -s)-$(uname -m) > /usr/local/bin/go-xamarin
```

Then:

```
chmod +x /usr/local/bin/go-xamarin
```

That's all, you're ready to go!

### Release Commits - 0.9.0 -> 1.0.0

* [9939077] Krisztian Godrei - bitrise.yml update (2016 Oct 12)
* [8062602] Krisztian Godrei - prepare for 1.0.0 (2016 Oct 12)
* [395f6a1] Krisztian Godrei - copy command in diagnostic mode, instead of using directly (2016 Oct 12)
* [05242f3] Krisztian Godrei - assembly name regexp fix and tests (2016 Oct 12)
* [07642b4] Krisztian Godrei - fix and test for AM PM comparsion (2016 Oct 12)
* [4ff6181] Krisztian Godrei - only print builder.buildableProjects warnings if no project found (2016 Oct 09)
* [ebdf593] Krisztian Godrei - validate solution config before build (2016 Oct 09)
* [f2bd494] Krisztián Gödrei - buildableProjects warnings (#16) (2016 Oct 09)
* [fb7831a] Krisztián Gödrei - builder tests (#15) (2016 Oct 09)
* [8bca748] Krisztian Godrei - analyze target definition fix (2016 Oct 06)
* [c3d47e8] Krisztián Gödrei - Duplicated commands (#14) (2016 Oct 06)
* [5f95d43] Krisztián Gödrei - timeout handling (#13) (2016 Oct 05)
* [1a27993] Krisztián Gödrei - Edit command (#12) (2016 Oct 05)
* [8974f9d] Krisztian Godrei - fixed errors (2016 Oct 03)
* [868e774] Krisztian Godrei - warning handling (2016 Oct 03)
* [457452d] Krisztian Godrei - dsym search fix (2016 Oct 03)
* [63d9004] Krisztián Gödrei - cleanup (#11) (2016 Oct 03)
* [b7ab76d] Krisztián Gödrei - Develop (#10) (2016 Sep 30)
* [24f8f5c] Krisztián Gödrei - buildtool package, logging updates (#8) (2016 Sep 30)
* [73f5ddb] Krisztián Gödrei - Build tools (#9) (2016 Sep 30)
* [dce37fc] Krisztián Gödrei - verbose mode, outputmap fix (#7) (2016 Sep 30)


-----------------

Updated: 2016 Oct 12