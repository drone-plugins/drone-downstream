Use this plugin to trigger builds for a list of downstream repositories. This
is useful when updates to a repository have downstream impacts that should also
be tested.

## Config

The following parameters are used to configure the plugin:

* **token** - secret that specifies the token for the drone server
* **repositories** - list of repositories to trigger
* **server** - the server on which you want to trigger builds

The following secret values can be set to configure the plugin.

* **DOWNSTREAM_TOKEN** - corresponds to **token**

It is highly recommended to put the **DOWNSTREAM_TOKEN** into a secret so it is
not exposed to users. This can be done using the drone-cli.

```bash
drone secret add --image=plugins/downstream \
    octocat/hello-world DOWNSTREAM_TOKEN my-secret-token
```

Then sign the YAML file after all secrets are added.

```bash
drone sign octocat/hello-world
```

See [secrets](http://readme.drone.io/0.5/usage/secrets/) for additional
information on secrets

## Examples

The following is a sample configuration in your .drone.yml file:

```yaml
notify:
  downstream:
    image: plugins/downstream
    server: https://ci.example.com
    repositories:
      - octocat/Hello-World
      - octocat/Spoon-Knife
    when:
      event: push
      branch: master
      success: true
```

In some cases you may want to trigger specific branches:

```
notify:
  downstream:
    image: plugins/downstream
    server: https://ci.example.com
    repositories:
      - octocat/Hello-World@develop
      - octocat/Spoon-Knife@master
```

You can also trigger a new build for a repository using `fork: true`:

```yaml
notify:
  downstream:
    image: plugins/downstream
    server: https://ci.example.com
    repositories:
      - octocat/Hello-World
      - octocat/Spoon-Knife
    fork: true
    when:
      event: push
      branch: master
      success: true
```

Please be sure to include the `when` section in your `.drone.yml` to avoid
triggering builds for pull requests, tags and failed builds. 
