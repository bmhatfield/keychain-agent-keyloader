# EXPERIMENTAL

### Keychain ssh-agent keyloader

This tool is designed to decrypt an encrypted SSH private key and load it into a running ssh-agent (similar to `ssh-add` functionality), with the added twist of getting the passphrase from the OSX Keychain.

This functionality (`ssh-add -K`) existed prior to OSX Sierra; now it depends on custom `.ssh/config` settings instead; perhaps you run a non-Apple version of SSH, perhaps you like to share your `.ssh/config` as part of a dotfiles repo - either way, the custom config options won't work for you. This tool is designed to help that.

Currently super experimental (is your environment exactly like mine?), but functional. Needs a lot more work. Loading the key into the Keychain is an exercise left up to the careful reader. The command line tool `security` is your friend.