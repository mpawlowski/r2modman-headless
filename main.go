package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/mpawlowski/r2modman-headless/r2modman"
	"github.com/mpawlowski/r2modman-headless/zip"

	"go.uber.org/fx"
)

const VERSION = "v0.0.1"

type flags struct {
	installDir string
	profileZip string
	workDir    string

	version      bool
	debugEnabled bool
}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "%s - Apply a profile export from r2modman to a dedicated server.\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example:\n\t%s --install-dir=serverfiles/ --work-dir=work/ --profile-zip=Profile.r2z\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags:\n")

	flag.PrintDefaults()
}

var options flags

func init() {

	flag.Usage = Usage

	options = flags{}
	flag.StringVar(&options.installDir, "install-dir", "", "Installation directory of the server.")
	flag.StringVar(&options.profileZip, "profile-zip", "", "Profile export to apply.")
	flag.StringVar(&options.workDir, "work-dir", "tmp/", "Temporary work directory for downloaded files.")
	flag.BoolVar(&options.version, "version", false, "Display the current version.")
	flag.BoolVar(&options.debugEnabled, "debug", false, "Enable verbose debugging.")
	flag.Parse()

	if options.version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if options.installDir == "" {
		log.Fatal("--install-dir must be defined")
	}
	if options.profileZip == "" {
		log.Fatal("--profile-zip must be defined")
	}

	if _, err := os.Stat(options.installDir); os.IsNotExist(err) {
		log.Fatal("provided --install-dir does not exist: " + options.installDir)
	}

	if _, err := os.Stat(options.profileZip); os.IsNotExist(err) {
		log.Fatal("provided --profile-zip does not exist: " + options.profileZip)
	}

	if _, err := os.Stat(options.workDir); os.IsNotExist(err) {
		log.Fatal("provided --work-dir does not exist: " + options.workDir)
	}

	//normalize directories
	options.installDir = path.Clean(options.installDir)
	options.workDir = path.Clean(options.workDir)

	fmt.Println("options", options)
}

func main() {

	var fxOptions []fx.Option

	if !options.debugEnabled {
		fxOptions = append(fxOptions, fx.NopLogger)
	}

	fxOptions = append(fxOptions,
		r2modman.Module(r2modman.Config{
			InstallDirectory: options.installDir,
			WorkDirectory:    options.workDir,
		}),
		zip.Module(zip.Config{}),
		fx.Invoke(run),
	)

	rand.Seed(time.Now().UnixNano())
	app := fx.New(fxOptions...)
	app.Run()
}

func run(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,

	parser r2modman.ExportParser,
	modutil r2modman.ModUtil,
	extractor zip.Extractor,
) {
	lc.Append(fx.Hook{OnStart: func(c context.Context) error {

		log.Println("r2modman-headless", VERSION)

		log.Println("using profile", options.profileZip)
		metadata, err := parser.Parse(options.profileZip)
		if err != nil {
			return err
		}

		for _, v := range metadata.Mods {

			downloadedZipPath := options.workDir + "/" + v.Filename()

			log.Println("downloading from", v.DownloadUrl())
			err := modutil.Download(v)
			if err != nil {
				fmt.Println("err", err)
				return err
			}

			packagingType, err := r2modman.DeterminePackagingType(downloadedZipPath)
			if err != nil {
				return err
			}
			log.Println("archive has type", packagingType.String())

			installDir := fmt.Sprintf("%s/%s", options.installDir, packagingType.Directory())

			log.Println("extracting to", installDir)

			err = extractor.Extract(downloadedZipPath, installDir)
			if err != nil {
				return err
			}
		}

		bepinDir := options.installDir + "/BepInEx"
		log.Println(fmt.Sprintf("extracting %s to %s", options.profileZip, bepinDir))
		err = extractor.Extract(options.profileZip, bepinDir)
		if err != nil {
			return err
		}

		log.Println("\n\n\tmod install finished successfully, configure your start script at", options.installDir+"/start_server_bepinex.sh\n")

		return shutdowner.Shutdown()

	}})
}
