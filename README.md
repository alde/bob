# Bob

## Purpose
Bob will try to identify what type of project you are executing it from, and based on configured values will run certain tasks in a docker container.

For example, if you're standing in a maven project (identified by the precense of `pom.xml`):
```bash
$ bob checkstyle
DEBU[0000] no local config file found, using default
INFO[0000] project identified as maven due to precense of pom.xml
INFO[0000] Running docker run --rm -w /workdir -e _JAVA_OPTIONS=-Duser.home=/home/alde -e HOME=/home/alde -v /home/alde:/home/alde -v /home/alde/git/mavenproject:/workdir maven mvn checkstyle:check

...

```

## Overriding targets
place a `.bob.yaml` in the project you want to override the defaults of. The structure (and current defaults) looks like the following:

```yaml
version: 1
projects:
- projectType: maven
  identifier: pom.xml
  dockerImage: maven
  environment:
    _JAVA_OPTIONS: -Duser.home=@homeDir
  volumes: {}
  commands:
    checkstyle:
    - mvn
    - checkstyle:check
    test:
    - mvn
    - clean
    - verify
```

`@homeDir` will be replaced by $HOME. Why it doesn't use $HOME directly is because I haven't yet decided how I want to support variable expansion.

