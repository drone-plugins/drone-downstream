Use this plugin to trigger builds for a list of downstream repositories. This is useful when updates to a repository have downstream impacts that should also be tested. These are the configuration options:

* `repositories` - list of repositories to trigger
* `token` - drone API token from your user setttings

The following is a sample configuration in your .drone.yml file:

```yaml
notify:
  downstream:
    repositories:
      - octocat/Hello-World
      - octocat/Spoon-Knife
    token: e3b0c44298fc1c149afbf4
    when:
      event: push
      branch: master
      success: true
```

In some cases you may want to trigger specific branches:

```
notify:
  downstream:
    repositories:
      - octocat/Hello-World@develop
      - octocat/Spoon-Knife@master
```

Please be sure to include the `when` section in your `.drone.yml` to avoid triggering builds for pull requests, tags and failed builds.

## Caveats

This plugin triggers the lastest build for a repository. It will not create a new build. This behavior may change in future versions of the plugin.
