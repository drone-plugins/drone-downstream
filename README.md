## Overview

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-downstream/status.svg)](http://beta.drone.io/drone-plugins/drone-downstream)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-downstream/coverage.svg)](https://aircover.co/drone-plugins/drone-downstream)
[![](https://badge.imagelayers.io/plugins/drone-downstream:latest.svg)](https://imagelayers.io/?images=plugins/drone-downstream:latest 'Get your own badge on imagelayers.io')

Drone plugin to trigger downstream repository builds

## Binary

Build the binary using `make`:

```
make deps build
```

### Example

```sh
./drone-anynines <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link_url": "https://beta.drone.io"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
        "repositories": [
            "octocat/Hello-World",
            "octocat/Spoon-Knife"
        ],
        "token": "eyJhbFciHiJISzI1EiIsUnR5cCW6IkpXQCJ9.ezH0ZXh0LjoidGJvZXJnZXIiLCJ0eXBlIjoidXNlciJ9.1m_3QFA6eA7h4wrBby2aIRFAEhQWPrlj4dsO_Gfchtc"
    }
}
EOF
```

## Docker

Build the container using `make`:

```
make deps docker
```

### Example

```sh
docker run -i plugins/drone-anynines <<EOF
{
    "repo": {
        "clone_url": "git://github.com/drone/drone",
        "owner": "drone",
        "name": "drone",
        "full_name": "drone/drone"
    },
    "system": {
        "link_url": "https://beta.drone.io"
    },
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "message": "Update the Readme",
        "author": "johnsmith",
        "author_email": "john.smith@gmail.com"
        "event": "push",
        "branch": "master",
        "commit": "436b7a6e2abaddfd35740527353e78a227ddcb2c",
        "ref": "refs/heads/master"
    },
    "workspace": {
        "root": "/drone/src",
        "path": "/drone/src/github.com/drone/drone"
    },
    "vargs": {
        "repositories": [
            "octocat/Hello-World",
            "octocat/Spoon-Knife"
        ],
        "token": "eyJhbFciHiJISzI1EiIsUnR5cCW6IkpXQCJ9.ezH0ZXh0LjoidGJvZXJnZXIiLCJ0eXBlIjoidXNlciJ9.1m_3QFA6eA7h4wrBby2aIRFAEhQWPrlj4dsO_Gfchtc"
    }
}
EOF
```
