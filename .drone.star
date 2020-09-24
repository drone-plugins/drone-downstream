go_image = 'golang:1.15.2-alpine3.12'

def main(ctx):
  before = testing(ctx)

  stages = [
    linux(ctx, 'amd64'),
    linux(ctx, 'arm64'),
    linux(ctx, 'arm'),
    windows(ctx, '1903'),
    windows(ctx, '1809'),
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
        'name': 'staticcheck',
        'image': go_image,
        'commands': [
          'go run honnef.co/go/tools/cmd/staticcheck ./...',
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/go',
          },
        ],
      },
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
    'repo': 'plugins/downstream',
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
        'image': go_image',
        'commands': [
          './release/linux/%s/drone-downstream --help' % (arch),
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

def windows(ctx, version):
  docker = [
    'echo $env:PASSWORD | docker login --username $env:USERNAME --password-stdin',
  ]

  if ctx.build.event == 'tag':
    build = [
      'go build -v -ldflags "-X main.version=%s" -a -o release/windows/amd64/drone-downstream.exe ./cmd/drone-downstream' % (ctx.build.ref.replace("refs/tags/v", "")),
    ]

    docker.extend([
      'docker build --pull -f docker/Dockerfile.windows.%s -t grafana/docker-downstream:%s-windows-%s-amd64 .' % (
        version, ctx.build.ref.replace("refs/tags/v", ""), version
      ),
      'docker run --rm grafana/docker-downstream:%s-windows-%s-amd64 --help' % (
        ctx.build.ref.replace("refs/tags/v", ""), version
      ),
      'docker push grafana/docker-downstream:%s-windows-%s-amd64' % (
        ctx.build.ref.replace("refs/tags/v", ""), version
      ),
    ])
  else:
    build = [
      'go build -v -ldflags "-X main.version=%s" -a -o release/windows/amd64/drone-downstream.exe ./cmd/drone-downstream' % (
        ctx.build.commit[0:8]
      ),
    ]

    docker.extend([
      'docker build --pull -f docker/Dockerfile.windows.%s -t grafana/docker-downstream:windows-%s-amd64 .' % (
        version, version
      ),
      'docker run --rm grafana/docker-downstream:windows-%s-amd64 --help' % (version),
      'docker push grafana/docker-downstream:windows-%s-amd64' % (version),
    ])

  return {
    'kind': 'pipeline',
    'type': 'ssh',
    'name': 'windows-%s' % (version),
    'platform': {
      'os': 'windows',
    },
    'server': {
      'host': {
        'from_secret': 'windows_server_%s' % (version),
      },
      'user': {
        'from_secret': 'windows_username',
      },
      'password': {
        'from_secret': 'windows_password',
      },
    },
    'steps': [
      {
        'name': 'environment',
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
        'environment': {
          'CGO_ENABLED': '0',
        },
        'commands': build,
      },
      {
        'name': 'executable',
        'commands': [
          './release/windows/amd64/drone-downstream.exe --help',
        ],
      },
      {
        'name': 'docker',
        'environment': {
          'USERNAME': {
            'from_secret': 'docker_username',
          },
          'PASSWORD': {
            'from_secret': 'docker_password',
          },
        },
        'commands': docker,
      },
    ],
    'depends_on': [],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
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
