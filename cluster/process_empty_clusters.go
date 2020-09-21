package cluster

// CloseEmptyChannel closes channel for empty
func CloseEmptyChannel() {

	close(empty_clusters_chan)
}
