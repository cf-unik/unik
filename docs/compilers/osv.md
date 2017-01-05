# OSv Unikernels

UniK supports [OSv](http://osv.io) unikernel creation.
OSv is somewhat different from other unikernel implementations.
Firstly, it can run application written in any programming language.
And secondly, due to some specifics in OSv kernel implementation, it can make use of
precompiled application packages that are stored in public repository.
In other words, in OSv you can not only run your own application (e.g. your NodeJS server),
but also arbitrary application from the public repository (e.g. MySQL).

UniK provides support for three languages: nodejs, java and native.
Native language refers to C/C++ compiled code that can be either your own
either from the remote repository.

Please note that UniK **does not** compile any of your application. Basically, all that it can
do for you is to upload your files to the unikernel.

---

## Node.js
It's trivial to deploy your Node.js application into unikernel. First you must collect all
external npm libraries with:
```
$ cd $PROJECT_ROOT
$ npm install
```
Then you can consult built-in documentation to see what configuration files are needed by UniK:
```
$ unik describe-compiler
   --base osv
   --language nodejs
   --provider openstack

OVERVIEW
Language "nodejs" allows you to run your NodeJS 4.4.5 application.
Please provide meta/run.yaml file where you describe your project
structure. See below for more details.

HOW TO PREPARE APPLICATION
Install all libraries using `npm install`.

CONFIGURATION FILES
------ /meta/run.yaml ------
config_set:
   conf1:
      main: <relative-path-to-your-entrypoint>
config_set_default: conf1

------ /manifest.yaml ------
image_size: "10GB"  # logical image size
```

As you can see above, you need to create folder `meta` inside your application and create
file `run.yaml` in it. In this file you specify where your entrypoint file is. For example, if your
entrypoint file is `server.js`, then the final content of `meta/run.yaml` would be:
```
config_set:
   conf1:
      main: /server.js
config_set_default: conf1
```
Optionally, you can create file `manifest.yaml` directly inside your project root directory
where you specify some additional parameters. Parameter `image_size` specifys filesystem size
of the created unikenrel.

### Build command
Use this command to build unikernel with your NodeJS application:
```
$ cd $PROJECT_ROOT
$ unik build
   --name myImg
   --path ./
   --base osv
   --language nodejs
   --provider openstack
```

## Java
It's trivial to deploy your Java application into unikernel. First you must compile it with javac:
```
$ cd $PROJECT_ROOT
$ javac ...
```
Then you can consult built-in documentation to see what configuration files are needed by UniK:
```
$ unik describe-compiler
   --base osv
   --language java
   --provider openstack

OVERVIEW
Language "java" allows you to run your Java 1.7 application.
Please provide meta/run.yaml file where you describe your project
structure. See below for more details.

HOW TO PREPARE APPLICATION
Compile your code into .class files using javac.

CONFIGURATION FILES
------ /meta/run.yaml ------
config_set:
   conf1:
      main: <name>
      classpath:
         <list>
      args:
         <list>
      jvmargs:
         <list>
config_set_default: conf1

------ /manifest.yaml ------
image_size: "10GB"  # logical image size
```

As you can see above, you need to create folder `meta` inside your application and create
file `run.yaml` in it. In this file you specify where your main file is. For example, if your
entrypoint file is `server.js`, then the final content of `meta/run.yaml` would be:
```
config_set:
   conf1:
      main: package1.Main
      classpath:
       - /src
      args:
      jvmargs:
       - Dhadoop.log.dir=/hdfs/logs
config_set_default: conf1
```
Optionally, you can create file `manifest.yaml` directly inside your project root directory
where you specify some additional parameters. Parameter `image_size` specifys filesystem size
of the created unikenrel.

### Build command
Use this command to build unikernel with your NodeJS application:
```
$ cd $PROJECT_ROOT
$ unik build
   --name myImg
   --path ./
   --base osv
   --language java
   --provider openstack
```

## Native
Native runtime can be used either for running your C/C++ application or for running
arbitrary precompiled application that exists in public repository (these are called
application packages).

### Running application package
Consult built-in documentation to see what configuration files are needed by UniK:
```
$ unik describe-compiler
   --base osv
   --language native
   --provider openstack

OVERVIEW
Language "native" allows you to either:
a) run precompiled package from the remote repository
b) run your own binary code
c) combination of (a) and (b)
In case you intend to use packages from the remote repository, please provide
meta/run.yaml where you list desired packages.
In any case (i.e. (a), (b), (c)) you need to provide meta/run.yaml. See below for
more details.

HOW TO PREPARE APPLICATION
(this is only needed if you want to run your own C/C++ application)
Compile your application into relocatable shared-object (a file normally
given a ".so" extension) that is PIC (position independent code).

CONFIGURATION FILES
------ /manifest.yaml ------
image_size: "10GB"  # logical image size

------ /meta/run.yaml ------
config_set:
   conf1:
      bootcmd: <boot-command-that-starts-application>
config_set_default: conf1

------ /meta/package.yaml ------
title: <your-unikernel-title>
name: <your-unikernel-name>
author: <your-name>
require:
  - <first-required-package-title>
  - <second-required-package-title>
  # ...

OTHER
Below please find a list of packages in remote repository:

eu.mikelangelo-project.app.hadoop-hdfs
eu.mikelangelo-project.app.mysql-5.6.21
eu.mikelangelo-project.erlang
eu.mikelangelo-project.ompi
eu.mikelangelo-project.openfoam.core
eu.mikelangelo-project.openfoam.pimplefoam
eu.mikelangelo-project.openfoam.pisofoam
eu.mikelangelo-project.openfoam.poroussimplefoam
eu.mikelangelo-project.openfoam.potentialfoam
eu.mikelangelo-project.openfoam.rhoporoussimplefoam
eu.mikelangelo-project.openfoam.rhosimplefoam
eu.mikelangelo-project.openfoam.simplefoam
eu.mikelangelo-project.osv.cli
eu.mikelangelo-project.osv.cloud-init
eu.mikelangelo-project.osv.httpserver
eu.mikelangelo-project.osv.nfs
```
As you can see above, you need to create folder `meta` inside your application and create
file `package.yaml` in it. In this file you specify what application packages do you want to
use. For example, if you want to run MySQL server, then the content of `meta/package.yaml`
would be:
```
name: lemmy.test
title: Lemmy Test
author: Lemmy
require:
 - eu.mikelangelo-project.app.mysql-5.6.21
```
Furthermore, you need to create file `meta/run.yaml` with the following content:
```
config_set:
   conf1:
      bootcmd: /usr/bin/mysqld --datadir=/usr/data --user=root --init-file=/etc/mysql-init.txt
config_set_default: conf1
```
And there you go, you can build your unikernel.

### Build command
Use this command to build unikernel with precompiled application (e.g. MySQL):
```
$ cd $PROJECT_ROOT
$ unik build
   --name myImg
   --path ./
   --base osv
   --language native
   --provider openstack
```

