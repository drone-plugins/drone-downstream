Use the downstream trigger plugin to trigger builds for other repositories. This is useful when updates to a repository have downstream impacts that should also be tested. These are the configuration options:

* **repos** - list of repositories to trigger
* **token** - drone API token from your user setttings

The following is a sample configuration in your .drone.yml file:

```yaml
notify:
  downstream:
    repos:
      - octocat/Hello-World
      - octocat/Spoon-Knife@master
    when:
      event: push
      branch: master
      success: true
```

Note that Drone will re-run the lastest build for a repository. It will not create a new build. This behavior may change in future versions of the plugin.