# hyprdisp
A utility program for automatically switching display profiles for Hyprland with Hyprpanel config  switching


hyprdisp detect

    - create the configuration files for the monitor configuration and puts in the sane defaults.
    - if the configuration files already exist, then do nothing

hyprdisp apply

    - read the configuration files and apply them to hyprland and hyprpanel

hyprdisp listen

    - listen for changes in the monitor configurations
    - does `hyprdisp detect` on changes
    - does `hyprdisp apply` afterwards

# Configuration

The program will write to the following files:

* `<hyprland-config-dir>/hyprland-monitors.conf`
* `<hyprland-config-dir>/hyprland-workspace.conf`

__Important:__

`hyprdisp` will automatically add the following line to `hyprland-monitors.conf`:

```
monitor = , preferred, auto, auto
```

To be able to automatically listen to any changes in monitor events when it gets added and removed. However, adding this
into another configuration file will make it undetectable to `hyprdisp` and will result into infinite loops.
