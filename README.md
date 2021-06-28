# missingFileFinder

Given a directory structure SRC (source), find all files not matching in directory structure DEST (destination)

Generates a json file:

```json
[
{"key":"20120721182752EE.jpg", "mc":1, "sc":1, "source":"testData/source/2012-07-21 18.27.52_EE.jpg", "dest":"testData/dest/sub1/2012-07-21 18.27.52_EE.jpg", "si":3520946, "match":"name+size"},
{"key":"20120721183027FF.jpg", "mc":1, "sc":1, "source":"testData/source/2012-07-21 18.30.27_FF.jpg", "dest":"testData/dest/2012-07-21 18.30.27_FF.jpg", "si":2710282, "match":"name+size"},
{"key":"20120721183440AA.jpg", "mc":2, "sc":2, "source":"testData/source/2012-07-21 18.34.40_AA.jpg", "dest":"testData/dest/2012-07-21 18.34.40_AA.jpg | testData/dest/sub1/2012-07-21 18.34.40_AA.jpg", "si":3766236, "match":"name+size"},
{"key":"20120721190625NM.jpg", "mc":0, "sc":1, "source":"testData/source/sub1/2012-07-21 19.06.25_NM.jpg", "dest":"", "si":2551696, "match":"no-match"},
{"key":"20151010122136GG.jpg", "mc":1, "sc":1, "source":"testData/source/sub1/20151010_122136_GG.jpg", "dest":"testData/dest/sub1/20151010_122136_GG.jpg", "si":2880596, "match":"name"}
]
```

Where the 'key' is the match name. If -c is passed on the command line teh key will have all chars '-'  '.'  ' ' and  '_' removed. 

Note that the last '.' before the file extenstion is always retained.

Where 'mc' is the count of matching files. This is grater than 1 if there are multiple matching files in the DEST directory.

Where 'sc' id the count of duplicates in the SRC directory.

Where 'source' is the relative path of the first (f duplcates) file in the source directory

Where 'dest' is a list of paths to the matching files in the destination directory. 

Where 'size' is the size of the first (f duplcates) file in the source directory.

Where 'match' is the match type.

* no-match: Meaning the source file were NOT found in the destination directory.
* name: Meaning files of the same name (key) but not the same size were found in the destination directory.
* name+size: Meaning files of the same name (key) and of the same size were found in the destination directory.

The process matches files with the specific extensions. 

* The default is any. 
* Use the 'ext:' command line option as follows to define a set of file extensions:

``` bash
ext:JPG ext:jpg ext:png ext:bmp ext:mp4 ext:gif
```
