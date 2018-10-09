# RaiseBuildEvent

Simple go thing to push a build event

Intended to send Cloudwatch Events during a CodePipeline run

curl -LO https://github.com/daraghmartin/RaiseBuildEvent/raw/master/RaiseBuildEvent

```
./RaiseBuildEvent \
                    -BuildName dockerx \
                    -BuildVersion 0.0.0 \
                    -BuildStatus Succeeded \
                    -Commit xyz \
                    -Source builds.docker  \
                    -DetailType build.notification \
                    -Comment "its a build"
```

Optional: -BuildDate=now, -DetailType=build.notification, -v=print things

BuildStatus will convert 0, 1 to Failed, Succeeded