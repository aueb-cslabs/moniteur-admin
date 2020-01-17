[![Build Status](https://travis-ci.com/aueb-cslabs/moniteur-admin.svg?branch=master)](https://travis-ci.com/aueb-cslabs/moniteur-admin) 
[![dependencies Status](https://david-dm.org/aueb-cslabs/moniteur-admin/status.svg)](https://david-dm.org/aueb-cslabs/moniteur-admin)
[![devDependencies Status](https://david-dm.org/aueb-cslabs/moniteur-admin/dev-status.svg)](https://david-dm.org/aueb-cslabs/moniteur-admin?type=dev) 

<a href="https://cslab.aueb.gr"><img src="https://www.aueb.gr/press/logos/2_AUEB-white-HR.jpg" title="AUEB CSLab" alt="AUEB"></a>

# Moniteur Administration

The administration tool for [Moniteur](https://github.com/aueb-cslabs/moniteur).

> Still WIP, but there is a working beta!

![](https://imgur.com/DpLgARW.png)

![](https://imgur.com/kckzf5f.png)

---

## Contents

- [Technologies](#technologies)
- [Building](#building)
- [Download](#download)
- [Platforms](#platforms)
- [Features](#features)
- [Team](#team)
- [License](#license)

---

## Technologies

Moniteur Administration is an application written in [Vue.js](https://vuejs.org/) and [electron](https://electronjs.org/). The application runs on Windows, Linux.

## Building

### Windows

     - Navigate to the downloaded repo
     - npm install
     - npm run build:win

### Linux

     - make .build-admin-linux

## Download

Head over to [releases](https://github.com/aueb-cslabs/moniteur-admin/releases) to get the latest version!

## Platforms
Moniteur Admin is supported on the following platforms:

- Windows (Installer, Portable)
- Linux (AppImage)

##### Warning! Auto updates only works on AppImage and Installer variants

---

## Features

Moniteur Administration provides REST capabilities for the administrators of Moniteur. Using this tool admins can

* Create general announcements
* Create announcements per room
* Create general comments
* Manage users who must have access to Administration
* Manage the academic calendar
* Create unscheduled lessons
* Change settings of the app

Examples:

> General Announcement

![](https://imgur.com/cxKmagS.png)

> Room Announcement

![](https://imgur.com/G7QTCxV.png)

> General Comment

![](https://imgur.com/tnfzS06.png)

> User Management

![](https://imgur.com/FkDfoHJ.png)

> Academic Calendar

![](https://imgur.com/H7vXP4P.png)

> Unscheduled Lessons

![](https://imgur.com/IclfmaG.png)

> Settings

![](https://imgur.com/GMEHlng.png)

---

## Team

Moniteur & Moniteur Administration was created for the needs of AUEB's Computer Science Laboratories (CSLab).

The development of the project was conducted by the [Moniteur Team](https://github.com/orgs/aueb-cslabs/teams/moniteur).

---

## License

GNU General Public License v3.0