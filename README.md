# missingFileFinder

Given a directory structure SRC (source), find all files not matching in directory structure DEST (destination)

## Output to a file
Generates a json file. The default name is 'result.json'. The command line parameter 'out:' will override this.

For example 'out:results-run.json' creates a file of the name esults-run.json.

```json
[
{"key":"20120721182752EE.jpg", "mc":1, "sc":1, "source":"testData/source/2012-07-21 18.27.52_EE.jpg", "dest":"testData/dest/sub1/2012-07-21 18.27.52_EE.jpg", "si":3520946, "match":"name+size"},
{"key":"20120721183027FF.jpg", "mc":1, "sc":1, "source":"testData/source/2012-07-21 18.30.27_FF.jpg", "dest":"testData/dest/2012-07-21 18.30.27_FF.jpg", "si":2710282, "match":"name+size"},
{"key":"20120721183440AA.jpg", "mc":2, "sc":2, "source":"testData/source/2012-07-21 18.34.40_AA.jpg", "dest":"testData/dest/2012-07-21 18.34.40_AA.jpg | testData/dest/sub1/2012-07-21 18.34.40_AA.jpg", "si":3766236, "match":"name+size"},
{"key":"20120721190625NM.jpg", "mc":0, "sc":1, "source":"testData/source/sub1/2012-07-21 19.06.25_NM.jpg", "dest":"", "si":2551696, "match":"no-match"},
{"key":"20151010122136GG.jpg", "mc":1, "sc":1, "source":"testData/source/sub1/20151010_122136_GG.jpg", "dest":"testData/dest/sub1/20151010_122136_GG.jpg", "si":2880596, "match":"name"}
]
```

## Output fields
Where the 'key' is the match name. If -c is passed on the command line teh key will have all chars '-'  '.'  ' ' and  '_' removed. 

Note that the last '.' before the file extenstion is always retained.

Where 'mc' is the count of matching files. This is grater than 1 if there are multiple matching files in the DEST directory.

Where 'sc' is the count of duplicates in the SRC directory.

Where 'source' is the relative path of the first (f duplcates) file in the source directory

Where 'dest' is a list of paths to the matching files in the destination directory. 

Where 'size' is the size of the first (f duplcates) file in the source directory.

Where 'match' is the match type.

* no-match: Meaning the source file was NOT found in the destination directory.
* name: Meaning files of the same name (key) were found but not of the same size in the destination directory.
* name+size: Meaning files of the same name (key) and of the same size were found in the destination directory.

The process matches files with the specific extensions. 

* The default is any. 
* Use the 'ext:' command line option as follows to define a set of file extensions. Dor not include a '.' before the extension. The matches are case sensitive. For example:

``` bash
ext:JPG ext:jpg ext:png ext:bmp ext:mp4 ext:gif
```
## Generate a bash script
A templated file can be generated for each line in the JSON result file.

The requires options are as follows:
* bashfile:\<file-name\>
  * \<file-name\> is the name of the generated file. For example 'bash-copy.sh'
* bash:\<match-type\> \<template\>
  * \<match-type\> Can be one of 'no-match', 'name', name+size' or 'all'
  * \<template\> is a string that contains substitution names. For example 'cp "$source" ./temp'

### Template values

The template can contain names from the JSON above and will be replaced with the values for each generated row.
* $source - The file name from the source directory
* 



## example 1
``` bash
./missingFileFinder source:testData/source dest:testData/dest -v -c out:results-run.json ext:JPG ext:jpg ext:png ext:bmp ext:mp4 ext:gif 'bash:no-match cp "$source" ./temp' bashfile:bash-run.sh 
```
Given the results json above:

* The source directory is: testData/source
* The destination directory is: testData/dest
* -v gives verbose output.
* -c Compresses file names, removing '-', ' ', '\_' and all '.' characters except the last.
* out:results-run.json will rename the output results file to results-run.json.
* ext:JPG will include only files with the '.JPG' extension.
* ext:JPG ext:jpg will include only files with the '.JPG' and '.jpg' extension.
  * There is no practical limit to the number of ext: parameters that can be defined
* The file 'bash-run.sh' will contain one line for each 'no-match' file. 
*   Each line will contain 
``` bash
cp testData/source/sub1/2012-07-21 19.06.25_NM.jpg ./temp
```

## example 2

``` bash
./missingFileFinder source:testData/source dest:testData/dest -v -c out:results-run.json ext:JPG ext:jpg ext:png ext:bmp ext:mp4 ext:gif 'bash:no-match cp "$source" ./temp' 
```
Given the results json above:

As no bashfile: is defined there will be no template file generated. Every thing else will be the same.

