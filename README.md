<p align="center">
    <img src="logo.png">
</p>

# Jack Hash Finder
Quick lookup for the original value of an hash

# Purpose
I was tired of looking up for common hashes values by hand. During CTFs you will eventually encounter some hashes. Instead of cracking them on your local machine or fire up a browser and look it up, the script does it for you. It tries some services to see if it's a common and known hash.

# Install
Install `golang` and then download and build the source code:
`go install github.com/jackrendor/jhf@latest`

On a linux machine, the binary should be located in `~/go/bin/`

On a windows machine, you should be able to run it from the terminal without specifying the path

# Example

```bash
jhf 21232f297a57a5a743894a0e4a801fc3
```
You can specify more than one hash
```bash
jhf b3ddbc502e307665f346cbd6e52cc10d 0bc11f2f3279555c317be9cf9e52645a
```
Or you can read from file by using `-f` or `--file`
```bash
jhf -file report/hashes.txt
```

