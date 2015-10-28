## Overview

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
