Use this plugin to trigger builds for a list of downstream repositories. This is useful when updates to a repository have downstream impacts that should also be tested. These are the configuration options:

* `repositories` - list of repositories to trigger
* `server` - the server on which you want to trigger builds

* `DOWNSTREAM_TOKEN` - A secret that specifies the token for the drone server

The following is a sample configuration in your .drone.yml file:

```yaml
notify:
  downstream:
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
    server: https://ci.example.com
    repositories:
      - octocat/Hello-World@develop
      - octocat/Spoon-Knife@master
```

Please be sure to include the `when` section in your `.drone.yml` to avoid triggering builds for pull requests, tags and failed builds.

You can also trigger a new build for a repository using `fork: true` :

```yaml
notify:
  downstream:
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
