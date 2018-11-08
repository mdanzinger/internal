
# composerToGit
composerToGit is a small script that removes the "composer" elements from our Wordpress projects.


### Usage
```
cd <PATH TO composerToGit>
./composerToGit -d <PATH TO PROJECT ROOT>
```



#### What does it do?
1. Removes composer.json and composer.lock files from project root
2. Removes .git files from all plugins and themes (to avoid having submodules)
3. Copies over any symlinked plugins / themes
4. Remove /themes/ and /plugins/ from `.gitignore`
5. *OPTIONAL* clears git cache, adds and commits the changes.
