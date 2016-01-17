## Overview

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-downstream/status.svg)](http://beta.drone.io/drone-plugins/drone-downstream)
[![](https://badge.imagelayers.io/plugins/drone-downstream:latest.svg)](https://imagelayers.io/?images=plugins/drone-downstream:latest 'Get your own badge on imagelayers.io')

This plugin is responsible for triggering downstream builds:

```
./drone-downstream <<EOF
{
    "system": {
        "link_url": "http://drone.mycompany.com"
    },
    "vargs": {
        "repositories": [
        	"octocat/Hello-World",
        	"octocat/Spoon-Knife"
        ],
        "token": "...."
    }
}
EOF
```
