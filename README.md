# krex

Kubernetes Resource Explorer

Krex works by building a directional graph (digraph) in memory of
various Kubernetes resources, and then giving you the graph one layer at
a time to explore using an interactive drop down menu. Explore
Kubernetes right from your terminal.

## Get Involved

Join us in the `#krex` channel in the [Kubernetes Slack community].

## Current state of krex

Handy tool for exploring applications in Kubernetes

## Future of krex

Global tool for exploring all things in Kubernetes

## Building krex

### Mac OSX

We use the C `ncurses` tool internally in Krex to navigate the terminal.

First install `ncurses`

``` bash
brew install ncurses
```

Export the `PKG_CONFIG_PATH` variable

``` bash
export PKG_CONFIG_PATH=/usr/local/opt/ncurses/lib/pkgconfig
```

### Ubuntu

Installing ncurses library in Debian/Ubuntu Linux

It should be already installed starting from `Ubuntu 18.04.2 LTS`, we are going to verify and install eventually.

Verify the libraries are installed

``` bash
/sbin/ldconfig -p | grep -i curses`
``` 
It should return something like:

``` bash
libncursesw.so.5 (libc6,x86-64) => /lib/x86_64-linux-gnu/libncursesw.so.5
libncurses.so.5 (libc6,x86-64) => /lib/x86_64-linux-gnu/libncurses.so.5
```

If not, let's install the following two packages: `libncurses5-dev` and `libncursesw5-dev`

``` bash
# Open a Terminal.
sudo apt-get update
sudo apt-get install libncurses5-dev libncursesw5-dev
```

Export the `PKG_CONFIG_PATH` variable

``` bash
export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig
```

### Continue the installation

Use symbolic links to adjust the library for `pkg-config`. (Note:
more information can be found [here].

``` bash
# Create the symbolic links for the missing files.
# It can be different between Mac and Ubuntu so a conditionnal execution is added here
# in case the link or file already exist.
if [ ! -f $PKG_CONFIG_PATH/form.pc -a ! -L $PKG_CONFIG_PATH/form.pc ]  ; then 
  ln -s /usr/lib/x86_64-linux-gnu/pkgconfig/formw.pc /usr/lib/x86_64-linux-gnu/pkgconfig/form.pc
fi
if [ ! -f $PKG_CONFIG_PATH/menu.pc -a ! -L $PKG_CONFIG_PATH/menu.pc ]  ; then
  ln -s /usr/lib/x86_64-linux-gnu/pkgconfig/menuw.pc /usr/lib/x86_64-linux-gnu/pkgconfig/menu.pc
fi
if [ ! -f $PKG_CONFIG_PATH/panel.pc -a ! -L $PKG_CONFIG_PATH/panel.pc ]  ; then
  ln -s /usr/lib/x86_64-linux-gnu/pkgconfig/panelw.pc /usr/lib/x86_64-linux-gnu/pkgconfig/panel.pc
fi
```

Then build the binary

``` bash
make build
```

and add to your path

``` bash
mv krex /usr/local/bin
```

[Kubernetes Slack community]: http://slack.k8s.io/
[here]: https://gist.github.com/cnruby/960344
