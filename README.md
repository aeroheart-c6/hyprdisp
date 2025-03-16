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
