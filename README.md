# xamarin-builder

## Building android Applications:  

command:  
`xbuild /t:SignAndroidPackage /p:Configuration=SOLUTION_CONFIGURATION SOLUTION_PATH`

example:  
`xbuild /t:SignAndroidPackage /p:Configuration=Ad-Hoc Path/To/MyProject.Android.csproj`
  
[source](https://forums.xamarin.com/discussion/68835/mdtool-archive-command-in-xamarin-studio-6)

## Building iOS Applications:  

__create ipa:__

command:  
`xbuild /p:Configuration=SOLUTION_CONFIGURATION /p:Platform=SOLUTION_PLATFORM  /p:BuildIpa=true /target:Build SOLUTION_PATH`

example:  
`xbuild /p:Configuration=AppStore /p:Platform=iPhone /p:BuildIpa=true /target:Build /path/to/solution.sln`

[source](https://developer.xamarin.com/guides/cross-platform/ci/teamcity/)

---

__create xcarchive:__

command:  
`xbuild /p:Configuration=PROJECT_CONFIGURATION /p:Platform=PROJECT_PLATFORM  /p:ArchiveOnBuild=true /target:Build PROJECT_PATH`

example:  
`xbuild /p:Configuration=Release /p:Platform=iPhone /p:ArchiveOnBuild=true /t:"Build" MyProject.csproj`

[source](https://forums.xamarin.com/discussion/42649/creating-archive-via-xbuild)
[source](http://chrisriesgo.com/icystt-command-line-tool-trouble-after-xamarin-cycle-7/)
