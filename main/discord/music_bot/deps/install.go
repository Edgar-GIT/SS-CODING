package deps

// InstallAll is kept for bootstrap compatibility. Voice DAVE support is
// provided by the discordgo fork (pure Go); no native libdave is required.
func InstallAll() error {
	return nil
}
