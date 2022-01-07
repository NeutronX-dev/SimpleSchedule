<p align="center">
    <img src="./logos/1500x500-SimpleSchedule.png"
        height="130">
</p>
<p align="center">
    <a href="https://go.dev/" alt="Made In">
        <img src="https://img.shields.io/badge/MADE IN-Go-blue?style=for-the-badge&logo=go&logoColor=white" /></a>
    <a href="https://github.com/NeutronX-dev/SimpleSchedule/releases/latest" alt="Version">
        <img src="https://img.shields.io/badge/VERSION-1.0.0-yellow?style=for-the-badge" /></a>
    <a href="https://github.com/NeutronX-dev/SimpleSchedule/graphs/contributors" alt="Version">
        <img src="https://img.shields.io/github/contributors/NeutronX-dev/SimpleSchedule?style=for-the-badge" /></a>
</p>

# <p align="center">Preview</p>

### No Items
<a src="./logos/screenshot/1.0.0/no-items.png">
    <img align="left" height="100" width=100 src="./logos/screenshot/1.0.0/no-items.png">
</a>

#### This is what the program will look like when it is first installed, or you have no upcoming events
```json
[]
```

---

### Add Items
<a src="./logos/screenshot/1.0.0/add-item.png">
    <img align="right" height="100" width=100 src="./logos/screenshot/1.0.0/add-item.png">
</a>

#### This is what the program will look like when you click on the "+ Add" Button
```json
[]
```

---

### With Items
<a src="./logos/screenshot/1.0.0/items.png">
    <img align="left" height="100" width=100 src="./logos/screenshot/1.0.0/items.png">
</a>

#### This is what the program will look like when you add an event.
```json
[ { "title": "Class", "time": 1641522463165 }, { "title": "Code", "time": 1641522463165 } ]
```

---

### Add Items
<a src="./logos/screenshot/1.0.0/event-triggered.png">
    <img align="right" height="100" width=100 src="./logos/screenshot/1.0.0/event-triggered.png">
</a>

#### This is what the program will look like when the time of an event passed. (as well as a sound effect)
```json
[ { "title": "Class", "time": 1641522463165 } ]
```

# Exit Codes
* Loading Errors (3-5)
* * **3**: Error loading Config (might be unexpected JSON input or missing permissions)
* * **4**: --
* * **5**: --
* Normal Errors (6-??)
* * **6**: Closed and threads Dispatched

# LICENSE
![gnu-logo](logos/gplv3-88x31.png)

This program is free software: you can redistribute it and/or modify
it under the terms of the [GNU General Public License](https://github.com/NeutronX-dev/ws.js/blob/main/LICENSE) as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
