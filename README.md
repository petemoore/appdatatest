# appdatatest

I had some problems with `APPDATA` and `LOCALAPPDATA` missing from the environment block after calling
[`SHSetKnownFolderPath`](https://msdn.microsoft.com/en-us/library/windows/desktop/bb762249(v=vs.85).aspx) so I wrote this
small program to test with.

You execute it, passing three arguments: a username, a password, and a filepath for the user's `APPDATA`. It will create a
user account with the given username and password, and then call `SHSetKnownFolderPath` to set the `APPDATA` folder location.
Then it will output the environment variables of the newly generated user.

My experience when using it, was if the `APPDATA` folder location was underneath the newly generated user's home directory,
then `APPDATA` would be included in the user's environment variables. If I used a folder location somewhere else on the same
drive, `APPDATA` would be missing from the environment. However, I was able to choose a location on a different drive, and
the environment variable would still be created. Note, even when the environment variable did not get created, all syscalls
are successful, and it is possible to query the location via
[`SHGetKnownFolderPath`](https://msdn.microsoft.com/en-us/library/windows/desktop/bb762188(v=vs.85).aspx), without problems.

This appears to be a bug in the Windows kernel, which I intend to report (but have not done yet)d.
