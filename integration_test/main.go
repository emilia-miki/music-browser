package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/emilia-miki/music-browser/integration_test/logger"
	"github.com/emilia-miki/music-browser/integration_test/test_case"
)

const BackendBaseUrl string = "http://localhost:3333"

func processRequest(
	tc test_case.TestCase,
) ([]byte, error) {
	req, err := tc.BuildRequest(BackendBaseUrl)
	if err != nil {
		return nil, errors.New("Unable to create request: " + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("Unable to send request: " + err.Error())
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Unable to read response body: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Message string `json:"message"`
		}
		err = json.Unmarshal(bytes, &errorResponse)
		if err != nil {
			return nil, fmt.Errorf(
				"Error %d: Unable to read error message: %s",
				resp.StatusCode, err.Error(),
			)
		}

		return nil, fmt.Errorf("Error %d: %s",
			resp.StatusCode, errorResponse.Message)
	}

	return bytes, nil
}

func startContainer(
	containerTool string,
	spotifyClientId string,
	spotifyClientSecret string,
) (func() error, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(
		ctx,
		containerTool, "run", "-p", "3333:3333",
		"-e", fmt.Sprintf("SPOTIFY_CLIENT_ID=%s", spotifyClientId),
		"-e", fmt.Sprintf("SPOTIFY_CLIENT_SECRET=%s", spotifyClientSecret),
		"music-browser")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("Error opening container stdout pipe: %s", err)
	}

	err = cmd.Start()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("Unable to start container: %s", err)
	}
	cIdB, err := exec.Command("sh", "-c", fmt.Sprintf(
		"%s ps | tail -1 | cut -d ' ' -f 1", containerTool)).Output()
	logger.Info.Printf(`The command sh -c "%s ps | tail -1 | cut -d ' ' -f 1"`+
		` returned %s`, containerTool, string(cIdB))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("Unable to get the container's ID: %s", err)
	}
	containerId := string(cIdB)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Server listening on port 3333") {
			go func() {
				for scanner.Scan() {
				}
				cancel()
			}()

			return func() error {
				logger.Info.Println("Sending SIGINT to the container")
				err = cmd.Process.Signal(os.Interrupt)
				if err != nil {
					return err
				}

				logger.Info.Println("Waiting for I/O to shut down")
				<-ctx.Done()

				logger.Info.Println("Waiting for the container to exit")
				err = cmd.Wait()
				if err != nil {
					return err
				}

				return nil
			}, nil
		}
	}

	cancel()
	return nil, fmt.Errorf("Couldn't start the server. "+
		"Inspect the container's logs by running %s logs %s",
		containerTool, containerId,
	)
}

func findContainerizationTool() (string, error) {
	containerizationTool := ""

	err := exec.Command("which", "docker").Run()
	if err == nil {
		containerizationTool = "docker"
	}

	err = exec.Command("which", "podman").Run()
	if err == nil {
		containerizationTool = "podman"
	}

	if containerizationTool == "" {
		return "", errors.New("Neither podman nor docker found in PATH")
	}

	return containerizationTool, nil
}

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "copies the log output to stdout")
	flag.BoolVar(&verbose, "v", false, "copies the log output to stdout")
	flag.Parse()

	logger.Initialize(verbose)

	containerizationTool, err := findContainerizationTool()
	if err != nil {
		logger.Error.Fatalln(err)
	}

	logger.Info.Println("Building the image with " + containerizationTool)
	cmd := exec.Command(containerizationTool,
		"build", "-t", "music-browser", ".")
	cmd.Dir = "../"

	b, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error.Fatalf(
			"Couldn't build the image: %s\n"+
				"%s build output:\n%s",
			err, containerizationTool, string(b),
		)
	}

	spotifyClientId := os.Getenv("SPOTIFY_CLIENT_ID")
	if spotifyClientId == "" {
		logger.Error.Fatalln(
			"SPOTIFY_CLIENT_ID environment variable is not set")
	}

	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if spotifyClientSecret == "" {
		logger.Error.Fatalln(
			"SPOTIFY_CLIENT_SECRET environment variable is not set")
	}

	spotifyTrackUrl := "https://open.spotify.com/track/2R3enCnuZSLovqEu1LCl6t"
	ytMusicTrackUrl := "https://music.youtube.com/watch?v=ziJUuC4y8s4"
	ytMusicTrackFileName := "suspyre.opus"

	bandcampTrackUrl := "https://macroblank.bandcamp.com/track/--91"
	bandcampTrackFileName := "手遅れです.opus"

	ytMusicTrack, err := os.ReadFile(ytMusicTrackFileName)
	if err != nil {
		logger.Error.Fatalf(
			"Unable to read file %s: %s\n", ytMusicTrackFileName, err)
	}
	spotifyTrack := ytMusicTrack
	bandcampTrack, err := os.ReadFile(bandcampTrackFileName)
	if err != nil {
		logger.Error.Fatalf(
			"Unable to read file %s: %s\n", bandcampTrackFileName, err)
	}

	logger.Info.Println("Starting the container")
	cancel, err := startContainer(
		containerizationTool, spotifyClientId, spotifyClientSecret)
	if err != nil {
		logger.Error.Fatalln("Error starting container: " + err.Error())
	}
	defer func() {
		logger.Info.Println("Shutting down the container")
		cancel()
	}()

	logger.Info.Println("Running tests")

	var testCases []test_case.TestCase

	testCases = append(testCases, test_case.NewSearch(
		test_case.SpotifyBackendName,
		test_case.ReqArtists,
		"Ne Obliviscaris",
	))
	testCases = append(testCases, test_case.NewSearch(
		test_case.YtMusicBackendName,
		test_case.ReqArtists,
		"Ne Obliviscaris",
	))
	testCases = append(testCases, test_case.NewSearch(
		test_case.BandcampBackendName,
		test_case.ReqArtists,
		"Death",
	))

	testCases = append(testCases, test_case.NewGet(
		test_case.GetTypeArtist,
		"https://open.spotify.com/artist/5kbidtcpyRRMdAQUnI1BG4"))
	testCases = append(testCases, test_case.NewGet(
		test_case.GetTypeArtist,
		"https://music.youtube.com/channel/UCevRJ36hDdpdG_4ujQfXfOA"))
	testCases = append(testCases, test_case.NewGet(
		test_case.GetTypeArtist,
		"https://macroblank.bandcamp.com"))

	testCases = append(testCases, test_case.NewGet(
		test_case.GetTypeAlbum,
		"https://open.spotify.com/album/29p0sILcBTXJmhzqJPzcxB"))
	testCases = append(testCases, test_case.NewGet(
		test_case.GetTypeAlbum,
		"https://music.youtube.com/browse/MPREb_AOBS6uUAmis"))
	testCases = append(testCases, test_case.NewGet(
		test_case.GetTypeAlbum,
		"https://macroblank.bandcamp.com/album/--5"))

	testCases = append(testCases, test_case.NewGet(
		test_case.GetTypeTrack,
		spotifyTrackUrl,
	))
	testCases = append(testCases, test_case.NewGet(
		test_case.GetTypeTrack,
		ytMusicTrackUrl,
	))
	testCases = append(testCases, test_case.NewGet(
		test_case.GetTypeTrack,
		bandcampTrackUrl,
	))

	testCases = append(testCases, test_case.NewDownload(spotifyTrackUrl))
	testCases = append(testCases, test_case.NewDownload(ytMusicTrackUrl))
	testCases = append(testCases, test_case.NewDownload(bandcampTrackUrl))

	testCases = append(testCases, test_case.NewGetStream(
		spotifyTrackUrl,
		spotifyTrack,
	))
	testCases = append(testCases, test_case.NewGetStream(
		ytMusicTrackUrl,
		ytMusicTrack,
	))
	testCases = append(testCases, test_case.NewGetStream(
		bandcampTrackUrl,
		bandcampTrack,
	))

	for i, tc := range testCases {
		logger.Info.Printf("Testing %s\n", tc)

		b, err = processRequest(tc)
		if err != nil {
			logger.Error.Printf(
				"Unable to process request for %s: %s\n",
				tc, err,
			)
		}
		testCases[i], err = tc.Check(b)
		if err != nil {
			logger.Error.Println(err)
		}
	}

	var passed int
	for _, result := range testCases {
		if result.DidPass() {
			logger.StdoutGreen.Println(result)
			passed += 1
		} else {
			logger.StdoutRed.Println(result)
		}
	}
	total := len(testCases)
	failed := total - passed

	if failed == 0 {
		logger.StdoutGreen.Printf(
			"\nSummary:\n%d/%d - ALL TESTS PASSED\n",
			passed, total,
		)
	} else {
		logger.StdoutRed.Printf(
			"\nSummary:\n%d/%d - %d TESTS FAILED\n",
			passed, total, failed,
		)
	}
	logger.Inform()
}
