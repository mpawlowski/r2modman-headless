# r2modman-headless

[![Build - Continuous](https://github.com/mpawlowski/r2modman-headless/actions/workflows/build_continuous.yaml/badge.svg?branch=main)](https://github.com/mpawlowski/r2modman-headless/actions/workflows/build_continuous.yaml)

Command line non-interactive mod loader for [r2modmanPlus](https://github.com/ebkr/r2modmanPlus) zip exports.

## usage

    r2modman-headless - Apply a profile export from r2modman to a dedicated server.
    Example:
            ./r2modman-headless --install-dir=serverfiles/ --work-dir=work/ --profile-zip=Profile.r2z
    Flags:
    -debug
            Enable verbose debugging.
    -install-dir string
            Installation directory of the server.
    -profile-zip string
            Profile export to apply.
    -version
            Display the current version.
    -work-dir string
            Temporary work directory for downloaded files. (default "tmp/")