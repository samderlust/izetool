# Sangtoolbox

This is an open source project to create an toolbox that help you saving time with boilerplate. Whatever your projects is.
The tool box helps run command line to create your project and then create folders and files as follow the structure you provide

At the moment, sangtoolbox only support creating flutter project.

# Install

1. Clone the project

   ```
   	git clone https://github.com/samderlust/sangtoolbox.git
   ```

2. provide the template file in `template` folder. Must be a json file
3. run `make sangtoolbox` in the root folder to install the toolbox into your system
4. run `sangtool flutter_create <name> --template=<template.json>`. if `--template` is not provide, the `example` will be used as default template.
