---
title: "Debug Role"
weight: 10
description: Debug Gravity's performance or logic
---

## Attach to a running gravity instance

Starting with 0.29.0, a gravity debug container image is built. You can use `:latest-debug` or `<commit-sha>-debug` to run a debug image.

In the container, you can run this command to start a debugging adapter and make it available on port 8011.

```shell
gravity cli debug dlv attach 1 --headless --listen=$INSTANCE_IP:8011
```

It is recommended to clone the Gravity Git Repository, open it in VS Code and continue from there.

Modify the included VS Code launch config file in `.vscode/launch.json` to point the IP to the machine you're running Gravity on.

With all of that setup, you can start debugging from VS Code, including setting breakpoints, etc.
