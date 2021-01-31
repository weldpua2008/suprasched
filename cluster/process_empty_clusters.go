package cluster

// CloseEmptyChannel closes channel for empty
func CloseEmptyChannel() {

	close(emptyClustersChan)
}
