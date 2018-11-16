i3 powertools
=============

User experience enhancements to i3 using the Go i3 IPC package.


## Features

* workspaces: switch to the next or previous workspace, while staying on the same monitor, aka output device. Very handy for multi-monitor workflows.
* more stuff coming over time.. feel free to fork + PR too


## Install/Usage

1. go get -u -d github.com/pkieltyka/i3-powertools
2. edit your i3 config file, and add:

```
bindsym $mod+Next exec i3-powertools -workspace=next
bindsym $mod+Prior exec i3-powertools -workspace=prev
```


## LICENSE

MIT
