## Todo list

## FEATURES:

## Clusters:
* capacity mesure
* empty clusters should be terminated
* terminate clusters

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
