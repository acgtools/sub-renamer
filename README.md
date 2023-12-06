# Sub-Renamer

English | [简体中文](./README_ZH_CN.md)

A CLI tool designed to automatically rename subtitle files corresponding to matched video files. 

<br/>

If this repo is helpful to you, please consider giving it a star (o゜▽゜)o☆ . Thank you OwO. 

> Random Wink OvO

<img src="https://waifu-getter.vercel.app/sfw?eps=wink" />

<br />

<!-- 
  If you prefer to use your own Moe-Counter
  please refer to the tutorial 
  in its original repo: https://github.com/journey-ad/Moe-Counter
  and deploy it to the Replit or Glitch
-->
![](https://political-capable-roll.glitch.me/get/@acgtoolssubrenamer?theme=rule34)

## Installation

### Using `go`

```sh
$ go install -ldflags "-w -s" github.com/acgtools/sub-renamer@latest
```

### Download from releases

[Release page](https://github.com/acgtools/sub-renamer/releases)

## Quick Start

```sh
$ sub-reanmer -h
sub-renamer <video dir> <sub dir>

Usage:
  sub-renamer [flags]

Flags:
  -c, --copy               copy subtitles after renaming
  -h, --help               help for sub-renamer
      --log-level string   log level, options: debug, info, warn, error (default "info")
  -v, --version            version for sub-renamer

```

### How to use on windows

1. Download the zip file and extract it to the folder you’d like.

2. Open the `cmd`: you can search for “cmd” in the Windows search bar, then click "Command Prompt"

3. Copy the path of the folder you just extracted. Input the following command and press `enter`:

   ```cmd
   cd /d "<path you copied>"
   ```

4. Copy the path of the video and subtitle folder, input the following command, and press `enter`:

   > **Caution:  Please ensure that the video file path is provided as the first argument.**
   >
   > Since it will rename files in the second parameter, this process is irreversible. Creating backups for video files is otherwise cumbersome due to their large size.

   ```cmd
   # Rename subtitle files only
   .\sub-renamer.exe "<video path>" "<subtitle path>"
   
   # Using -c will copy subtitle files to the video folder after renaming
   .\sub-renamer.exe -c "<video path>" "<subtitle path>"
   ```

![](./docs/assets/how_to_use.gif)

## Issue

Feel free to create issues to report bugs or request new features. 

