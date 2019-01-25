# escape=`
FROM plugins/base:windows-1809

LABEL maintainer="Drone.IO Community <drone-dev@googlegroups.com>" `
  org.label-schema.name="Drone Downstream" `
  org.label-schema.vendor="Drone.IO Community" `
  org.label-schema.schema-version="1.0"

ADD release/windows/amd64/drone-downstream.exe C:/bin/drone-downstream.exe
ENTRYPOINT [ "C:\\bin\\drone-downstream.exe" ]
