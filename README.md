# go-xamarin

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

## Analyze builders

# git@github.com:bitrise-samples/xamarin-sample-app.git Release|iPhone

__go-xamarin:__

```
=> /Library/Frameworks/Mono.framework/Commands/xbuild "/Users/godrei/Develop/xamarin/xamarin-sample-app/Droid/XamarinSampleApp.Droid.csproj" "/target:PackageForAndroid" "/p:Configuration=Release" "/verbosity:minimal" "/nologo"
=> /Applications/Xamarin Studio.app/Contents/MacOS/mdtool "build" "/Users/godrei/Develop/xamarin/xamarin-sample-app/XamarinSampleApp.sln" "-c:Release|iPhone" "-p:XamarinSampleApp.iOS"
=> /Applications/Xamarin Studio.app/Contents/MacOS/mdtool "archive" "/Users/godrei/Develop/xamarin/xamarin-sample-app/XamarinSampleApp.sln" "-c:Release|iPhone" "-p:XamarinSampleApp.iOS"
```

---

__steps-xamarin-builder:__

```
["\"/Applications/Xamarin Studio.app/Contents/MacOS/mdtool\"", "build", "\"-c:Release|iPhone\"", "\"XamarinSampleApp.sln\"", "\"-p:XamarinSampleApp.iOS\""]
["\"/Applications/Xamarin Studio.app/Contents/MacOS/mdtool\"", "archive", "\"-c:Release|iPhone\"", "\"XamarinSampleApp.sln\"", "\"-p:XamarinSampleApp.iOS\""]
["\"/Library/Frameworks/Mono.framework/Commands/xbuild\"", "/t:PackageForAndroid", "/p:Configuration=\"Release\"", "\"./Droid/XamarinSampleApp.Droid.csproj\"", "/verbosity:minimal", "/nologo"]
```

# git@github.com:bitrise-samples/xamarin-sample-app.git Debug|Any CPU

__go-xamarin:__

```
=> /Applications/Xamarin Studio.app/Contents/MacOS/mdtool "build" "/Users/godrei/Develop/xamarin/xamarin-sample-app/XamarinSampleApp.sln" "-c:Debug|iPhoneSimulator" "-p:XamarinSampleApp.iOS
=> /Library/Frameworks/Mono.framework/Commands/xbuild "/Users/godrei/Develop/xamarin/xamarin-sample-app/Droid/XamarinSampleApp.Droid.csproj" "/target:PackageForAndroid" "/p:Configuration=Debug" "/verbosity:minimal" "/nologo"
```

---

__steps-xamarin-builder:__

```
["\"/Applications/Xamarin Studio.app/Contents/MacOS/mdtool\"", "build", "\"-c:Debug|iPhoneSimulator\"", "\"XamarinSampleApp.sln\"", "\"-p:XamarinSampleApp.iOS\""]
["\"/Library/Frameworks/Mono.framework/Commands/xbuild\"", "/t:PackageForAndroid", "/p:Configuration=\"Debug\"", "\"./Droid/XamarinSampleApp.Droid.csproj\"", "/verbosity:minimal", "/nologo"]
```

# xamarin-mac Debug|x86

__go-xamarin:__

```
=> /Applications/Xamarin Studio.app/Contents/MacOS/mdtool "build" "/Users/godrei/Develop/xamarin/xamarin-mac/XamarinMac/XamarinMac.sln" "-c:Debug|x86" "-p:XamarinMac"
=> /Applications/Xamarin Studio.app/Contents/MacOS/mdtool "archive" "/Users/godrei/Develop/xamarin/xamarin-mac/XamarinMac/XamarinMac.sln" "-c:Debug|x86" "-p:XamarinMac"
```

---

__steps-xamarin-builder:__

```
["\"/Applications/Xamarin Studio.app/Contents/MacOS/mdtool\"", "build", "\"-c:Debug|x86\"", "\"/Users/godrei/Develop/xamarin/xamarin-mac/XamarinMac/XamarinMac.sln\"", "\"-p:XamarinMac\""]
["\"/Applications/Xamarin Studio.app/Contents/MacOS/mdtool\"", "archive", "\"-c:Debug|x86\"", "\"/Users/godrei/Develop/xamarin/xamarin-mac/XamarinMac/XamarinMac.sln\"", "\"-p:XamarinMac\""]
```

