go_image = 'golang:1.15.2-alpine3.12'

def main(ctx):
  return pipelines('release') + pipelines('master') + pipelines('pr')

def pipelines(ver_mode):
  before = testing(ver_mode)

  stages = [
    linux('amd64', ver_mode),
  ]

  after = manifest(ver_mode)

  for b in before:
    for s in stages:
      s['depends_on'].append(b['name'])

  for s in stages:
    for a in after:
      a['depends_on'].append(s['name'])

  return before + stages + after

def testing(ver_mode):
  return [
    {
      'kind': 'pipeline',
      'type': 'docker',
      'name': '{}-testing'.format(ver_mode),
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
      'trigger': get_triggers(ver_mode),
    },
  ]

def linux(arch, ver_mode):
  docker = {
    'dockerfile': 'docker/Dockerfile.linux.%s' % (arch),
    'repo': 'grafana/drone-downstream',
    'dry_run': True,
    'tags': 'linux-{}'.format(arch),
  }

  if ver_mode in ('master', 'release'):
    docker.update({
      'username': {
        'from_secret': 'docker_username',
      },
      'password': {
        'from_secret': 'docker_password',
      },
      'auto_tag': True,
      'auto_tag_suffix': 'linux-{}'.format(arch),
      'dry_run': False,
    })

  if ver_mode == 'release':
    build = [
      'REF=$(echo ${DRONE_COMMIT_REF} | sed \'s/refs\\/tags\\/v//\')',
      'go build -v -ldflags "-X main.version=$${{REF}}" -a -o release/linux/{}/drone-downstream ./cmd/drone-downstream'.format(
        arch
      ),
    ]
  else:
    build = [
      'COMMIT=${DRONE_COMMIT:0:8}',
      'go build -v -ldflags "-X main.version=$${{COMMIT}}" -a -o release/linux/{}/drone-downstream ./cmd/drone-downstream'.format(
        arch
      ),
    ]

  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': '{}-linux-{}'.format(ver_mode, arch),
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
        'depends_on': [
          'executable',
        ],
      },
    ],
    'depends_on': [],
    'trigger': get_triggers(ver_mode),
  }

def manifest(ver_mode):
  if ver_mode not in ('master', 'release'):
    return []

  return [
    {
      'kind': 'pipeline',
      'type': 'docker',
      'name': '{}-manifest'.format(ver_mode),
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
      'trigger': get_triggers(ver_mode),
    },
  ]

def get_triggers(ver_mode):
  if ver_mode == 'pr':
    return {
      'event': ['pull_request',],
    }

  if ver_mode == 'master':
    return {
      'ref': [
        'refs/heads/master',
      ]
    }

  if ver_mode == 'release':
    return {
      'ref': [
        'refs/tags/**',
      ]
    }
