# izetool

This is an opensource project to create an toolbox that help you saving time with boiledplate. Whatever your projects is.
The tool box helps run commandline to create your project and then create folders and files as follow the structure you provide

At the moment, izetool only support creating flutter project.

#### support: Linux, MacOS

# Install

1. Clone the project

   ```
   	git clone https://github.com/samderlust/izetool.git
   ```

2. run `make izetool` in the root folder to install the toolbox into your system
3. run `ize init` to create `ize_templates` folder at home dir. This is where you store your template files

# Available commands

| command                | usage                                                                            | note                                                                                                                              |
| ---------------------- | -------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| flutter create         | `ize flutter_create <name> --template=<template>` or `ize flutter_create <name>` | if `--template` is not provide, the `example.json` will be used as default template.                                              |
| flutter uploadkeystore | `ize flutter uploadkeystore`                                                     | process to create android upload keystore. After that, create `key.properties` file and also modify your `app/build.gradle` file. |
| make                   | `ize make <template> --name=<name>`                                              | must provide a template.`<template>` is the same as file name                                                                     |
