# r2modman-headless

[![Build - Continuous](https://github.com/mpawlowski/r2modman-headless/actions/workflows/build_continuous.yaml/badge.svg?branch=main)](https://github.com/mpawlowski/r2modman-headless/actions/workflows/build_continuous.yaml)

Command line non-interactive mod loader for [r2modmanPlus](https://github.com/ebkr/r2modmanPlus) zip exports.

## usage

        r2modman-headless - Apply a profile export from r2modman to a dedicated server.
        Example:
                r2modman-headless --install-dir=serverfiles/ --work-dir=work/ --profile-zip=Profile.r2z
        Flags:
        -install-dir string
                Installation directory of the server.
        -profile-zip string
                Profile export to apply.
        -thunderstore-cdn-host string
                Hostname of the thunderstore CDN to use. (default "gcdn.thunderstore.io")
        -thunderstore-cdn-timeout duration
                Timeout while downloading each mod. (default 30s)
        -thunderstore-force-download
                Force re-download of all mods, even if they are already present in the work directory.
        -work-dir string
                Temporary work directory for downloaded files. (default "tmp/")
                