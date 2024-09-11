## [Project] Distributed Short Video Engagement Aggregation

### Tech Doc : [here](https://docs.google.com/document/d/1aVsR6AD03AT_RphhOUQOxIasR57CvPzBJ1xGpde5ILU/edit?usp=sharing)

### HLD

![](resources/misc/images/service_design.png)

### How to run?

- Use command `make`. 
- You can find initialisation scripts in `Makefile` and `resources/local-setup/scripts/`
- If you use other config files, make sure to update `APP_ENV` env var.

#### API cURLS
- fetch aggregate views by viewer-id and timestamp

```
curl --location '{host}:8000/api/v1/viewer-count?video_id=2omaodmas1&timestamp_in_min=129139'
```
