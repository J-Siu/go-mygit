package lib

const (
	TXT_FLAGS_USE       = "If -r/-g not specified, all remotes in config file will be used. "
	TXT_REPO_CLONE_LONG = "Provide at least one repository name(not url). Need exactly one -r/--remote. Do not accept -g/--group."
	TXT_REPO_CLONE_USE  = "<repository ...>"
	TXT_REPO_DIR_LONG   = "If no repository is specified, current git root will be used as repository name. "
	TXT_REPO_DIR_USE    = "[repository/directory ...]"
)

var Conf TypeConf
var Flag TypeFlag
