go_image = 'golang:1.15.2-alpine3.12'

def main(ctx):
  before = testing(ctx)

  stages = [
    linux(ctx, 'amd64'),
  ]

  after = manifest(ctx)

  for b in before:
    for s in stages:
      s['depends_on'].append(b['name'])

  for s in stages:
    for a in after:
      a['depends_on'].append(s['name'])

  return before + stages + after

def testing(ctx):
  return [{
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'testing',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'lint',
        'image': 'golangci/golangci-lint:v1.31.0-alpine',
        'commands': [
          'golangci-lint run',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/go',
          },
        ],
      },
      {
        'name': 'test',
        'image': go_image,
        'environment': {
          'CGO_ENABLED': 0,
        },
        'commands': [
          'go test -cover ./...',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/go',
          },
        ],
      },
    ],
    'volumes': [
      {
        'name': 'gopath',
        'temp': {},
      },
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**',
      ],
    },
  }]

def linux(ctx, arch):
  docker = {
    'dockerfile': 'docker/Dockerfile.linux.%s' % (arch),
    'repo': 'grafana/drone-downstream',
    'username': {
      'from_secret': 'docker_username',
    },
    'password': {
      'from_secret': 'docker_password',
    },
  }

  if ctx.build.event == 'pull_request':
    docker.update({
      'dry_run': True,
      'tags': 'linux-%s' % (arch),
    })
  else:
    docker.update({
      'auto_tag': True,
      'auto_tag_suffix': 'linux-%s' % (arch),
    })

  if ctx.build.event == 'tag':
    build = [
      'go build -v -ldflags "-X main.version=%s" -a -o release/linux/%s/drone-downstream ./cmd/drone-downstream' % (ctx.build.ref.replace("refs/tags/v", ""), arch),
    ]
  else:
    build = [
      'go build -v -ldflags "-X main.version=%s" -a -o release/linux/%s/drone-downstream ./cmd/drone-downstream' % (ctx.build.commit[0:8], arch),
    ]

  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'linux-%s' % (arch),
    'platform': {
      'os': 'linux',
      'arch': arch,
    },
    'steps': [
      {
        'name': 'environment',
        'image': go_image,
        'environment': {
          'CGO_ENABLED': '0',
        },
        'commands': [
          'go version',
          'go env',
        ],
      },
      {
        'name': 'build',
        'image': go_image,
        'environment': {
          'CGO_ENABLED': '0',
        },
        'commands': build,
      },
      {
        'name': 'executable',
        'image': go_image,
        'commands': [
          './release/linux/%s/drone-downstream --help' % (arch),
        ],
        'depends_on': [
          'build',
        ],
      },
      {
        'name': 'docker',
        'image': 'plugins/docker',
        'pull': 'always',
        'settings': docker,
      },
    ],
    'depends_on': [],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**',
      ],
    },
  }

def manifest(ctx):
  return [{
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'manifest',
    'steps': [
      {
        'name': 'manifest',
        'image': 'plugins/manifest',
        'pull': 'always',
        'settings': {
          'auto_tag': 'true',
          'username': {
            'from_secret': 'docker_username',
          },
          'password': {
            'from_secret': 'docker_password',
          },
          'spec': 'docker/manifest.tmpl',
          'ignore_missing': 'true',
        },
      },
    ],
    'depends_on': [],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
      ],
    },
  }]
