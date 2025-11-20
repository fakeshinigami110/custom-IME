complete -c custom-ime -f

complete -c custom-ime -n "__fish_use_subcommand" -a create -d "Create a new IME project"
complete -c custom-ime -n "__fish_use_subcommand" -a install -d "Build and install an IME project"
complete -c custom-ime -n "__fish_use_subcommand" -a list -d "List all IME projects"
complete -c custom-ime -n "__fish_use_subcommand" -a delete -d "Delete sources or uninstall an IME"

# create command
complete -c custom-ime -n "__fish_seen_subcommand_from create" -s p -l project -d "Project name (required)" -r
complete -c custom-ime -n "__fish_seen_subcommand_from create" -s n -l name -d "IME name (required)" -r
complete -c custom-ime -n "__fish_seen_subcommand_from create" -s l -l label -d "IME label (default: Custom)" -r
complete -c custom-ime -n "__fish_seen_subcommand_from create" -s i -l icon -d "Icon name (default: fcitx5-keyboard)" -r
complete -c custom-ime -n "__fish_seen_subcommand_from create" -s L -l lang -d "Language code (default: en)" -r
complete -c custom-ime -n "__fish_seen_subcommand_from create" -s D -l desc -d "Description (default: Custom IME)" -r
complete -c custom-ime -n "__fish_seen_subcommand_from create" -s c -l config -d "Custom config file path" -r
complete -c custom-ime -n "__fish_seen_subcommand_from create" -s f -l force -d "Force overwrite existing project"

# install command
complete -c custom-ime -n "__fish_seen_subcommand_from install" -s p -l project -d "Project name to install" -r

# delete command
complete -c custom-ime -n "__fish_seen_subcommand_from delete" -s p -l project -d "Project name to delete" -r
complete -c custom-ime -n "__fish_seen_subcommand_from delete" -s d -l delete -d "Delete source files"
complete -c custom-ime -n "__fish_seen_subcommand_from delete" -s u -l uninstall -d "Uninstall from fcitx5"