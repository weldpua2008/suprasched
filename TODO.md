## Todo list

## FEATURES:

## Clusters:
* capacity mesure
* empty clusters should be terminated
* terminate clusters
* Adaptive EMR describe -
    if we got no new EMR from API then we should describe old with interval more then  1 minute
## Jobs:
* unassigned jobs must be canceled

## BUGs:
### Jobs:
* No Cluster Information - Potentially Job fetch no new information about CLuster
```golang
(j *Job) updateStatus
```
### Clusters:
* Decomission cluster

### OTHER:
* Tickers aren't update with config update
