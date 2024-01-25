package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/mpawlowski/r2modman-headless/r2modman"
	"github.com/mpawlowski/r2modman-headless/zip"

	"go.uber.org/fx"
)

type flags struct {
	runTimeout                time.Duration
	installDir                string
	profileZip                string
	workDir                   string
	thunderstoreForceDownload bool
	thunderstoreCDNHost       string
	thunderstoreCdnTimeout    time.Duration
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
	flag.DurationVar(&options.runTimeout, "run-timeout", 5*time.Minute, "Total maximum runtime before giving up.")
	flag.StringVar(&options.installDir, "install-dir", "", "Installation directory of the server.")
	flag.StringVar(&options.profileZip, "profile-zip", "", "Profile export to apply.")
	flag.StringVar(&options.workDir, "work-dir", "tmp/", "Temporary work directory for downloaded files.")
	flag.BoolVar(&options.thunderstoreForceDownload, "thunderstore-force-download", false, "Force re-download of all mods, even if they are already present in the work directory.")
	flag.StringVar(&options.thunderstoreCDNHost, "thunderstore-cdn-host", "gcdn.thunderstore.io", "Hostname of the thunderstore CDN to use.")
	flag.DurationVar(&options.thunderstoreCdnTimeout, "thunderstore-cdn-timeout", 30*time.Second, "Timeout while downloading each mod.")
	flag.Parse()

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
}

func main() {

	var fxOptions []fx.Option

	fxOptions = append(fxOptions,
		r2modman.Module(r2modman.Config{
			InstallDirectory:          options.installDir,
			WorkDirectory:             options.workDir,
			ThunderstoreCDNTimeout:    options.thunderstoreCdnTimeout,
			ThunderstoreCDN:           options.thunderstoreCDNHost,
			ThunderstoreForceDownload: options.thunderstoreForceDownload,
		}),
		zip.Module(zip.Config{}),
		fx.Invoke(registerLifecycleHooks),
	)

	app := fx.New(fxOptions...)

	// start
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		panic(err)
	}

	// wait for app to finish
	signal := <-app.Wait()

	// stop
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	err = app.Stop(stopCtx)
	if err != nil {
		panic(err)
	}

	//exit with signal
	os.Exit(signal.ExitCode)
}

func run(
	ctx context.Context,
	parser r2modman.ExportParser,
	modutil r2modman.ModUtil,
	extractor zip.Extractor,
) error {
	log.Println("Using profile", options.profileZip)
	metadata, err := parser.Parse(options.profileZip)
	if err != nil {
		return err
	}

	packages, err := r2modman.GetPackagesMetadata(ctx)
	if err != nil {
		log.Printf("unable to pull thunderstore api: %s", err)
		return err
	}

	log.Printf("Found packages %d packages from thunderstore\n", len(packages))

	for _, v := range metadata.Mods {

		log.Printf("processing %s\n", v.Name)

		downloadedZipPath := path.Join(options.workDir, v.Filename())

		thunderstoreMeta, ok := packages[v.ThunderstoreKey()]
		if !ok {
			return fmt.Errorf("thunderstore metadata does not exist for: %s", v.ThunderstoreKey())
		}

		err = modutil.Download(v, thunderstoreMeta)
		if err != nil {
			return err
		}

		//extract modes to install directory
		packagingType, prefixToStrip, err := r2modman.DeterminePackagingType(downloadedZipPath)
		if err != nil {
			return err
		}

		installDir := fmt.Sprintf("%s/%s", options.installDir, packagingType.Directory())

		if packagingType == r2modman.ModPackagingTypePlugin {
			installDir += fmt.Sprintf("/plugins/%s", v.Name)
		}

		if packagingType == r2modman.ModPackagingTypeRootDLL {
			installDir += fmt.Sprintf("/%s", v.Name)
		}

		log.Printf("Packaging Type: %s", packagingType)
		log.Printf("Install Dir: %s", installDir)
		log.Printf("Prefix Strip: %s", prefixToStrip)
		err = extractor.Extract(downloadedZipPath, installDir, prefixToStrip)
		if err != nil {
			return err
		}

		log.Printf("")
	}

	// extract profile to bepinex in install dir
	bepinDir := path.Join(options.installDir, "/BepInEx")
	err = extractor.Extract(options.profileZip, bepinDir, "BepInEx")
	if err != nil {
		return err
	}

	log.Printf("Mod install finished successfully, configure your start script at %s/start_server_bepinex.sh\n", options.installDir)

	return nil
}

func registerLifecycleHooks(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,

	parser r2modman.ExportParser,
	modutil r2modman.ModUtil,
	extractor zip.Extractor,
) {

	lc.Append(fx.Hook{OnStart: func(c context.Context) error {

		errorChannel := make(chan error)

		go func(returnChannel chan<- error) {
			myCtx, myCtxCancel := context.WithTimeout(context.Background(), options.runTimeout)
			defer myCtxCancel()
			returnChannel <- run(myCtx, parser, modutil, extractor)

		}(errorChannel)

		go func(errChan chan error, shutdowner fx.Shutdowner) {
			err := <-errChan
			if err != nil {
				shutdowner.Shutdown(fx.ExitCode(1))
				return
			}
			shutdowner.Shutdown(fx.ExitCode(0))
		}(errorChannel, shutdowner)

		return nil

	}})

}
