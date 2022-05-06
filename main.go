package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v43/github"
)

func main() {
	// Get download URL
	url, err := listenSource()
	if err != nil {
		panic(err)
	}

	// Get prefix of install path
	install_path := listenInstall()

	// Create tmp directory to download and unfreeze
	dirpath, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dirpath)

	// Download and unfreeze file
	tar_file, err := ioutil.TempFile(dirpath, "")
	if err != nil {
		panic(err)
	}
	download(url, tar_file.Name())
	unfreeze(install_path, tar_file.Name())

	// Remove all tmp files and directory
	err = os.Remove(tar_file.Name())
	if err != nil {
		println("Unable to remove downloaded tar file")
	}
}

func unfreeze(parent_dir string, filepath string) {
	reader, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		panic(err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		filename := parent_dir + "/" + header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(filename, os.FileMode(header.Mode))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case tar.TypeReg:
			writer, err := os.Create(filename)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			writer.Close()
		default:
			fmt.Printf("Unable to untar type: %c in file %s", header.Typeflag, filename)
		}
	}
}

func listenInstall() string {
	i := question("Input install path")
	return i
}

func listenSource() (string, error) {
	choices := []string{"GitHub"}
	i := selective("Select source", choices)
	// GitHub
	if i == 0 {
		owner_repo_name := question("Enter repository name (e.g. flat35hd99/ig)")
		owner, repo, err := split_owner_and_repo(owner_repo_name)
		if err != nil {
			log.Fatal(err)
		}
		// repository service -> ListReleases
		client := github.NewClient(nil)
		releases, _, err := client.Repositories.ListReleases(context.Background(), owner, repo, &github.ListOptions{})
		if err != nil {
			fmt.Println("Trouble happen")
		}

		var release_names []string
		for _, r := range releases {
			release_names = append(release_names, *r.Name)
		}
		release_index := selective("Which release?", release_names)
		release := releases[release_index]

		assets := release.Assets
		var asset_names []string
		for _, a := range assets {
			asset_names = append(asset_names, *a.Name)
		}
		asset_index := selective("Which assets?", asset_names)
		asset := assets[asset_index]

		url := asset.BrowserDownloadURL
		return *url, nil
	}
	return "", nil
}

func download(url string, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func split_owner_and_repo(s string) (string, string, error) {
	slice := strings.Split(s, "/")
	if len(slice) != 2 {
		return "", "", fmt.Errorf("%s is not owner/repo format", s)
	}
	return slice[0], slice[1], nil
}

// Show question and return answer
func question(q string) string {
	fmt.Println("> ", q)

	var result string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		result = scanner.Text()
		break
	}
	return result
}

/*
Show question and choices and return index choiced
*/
func selective(q string, choices []string) int {
	fmt.Println("> ", q)
	for i, choice := range choices {
		fmt.Printf("%d: %s\n", i, choice)
	}

	var result_index int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		// Check string inputs
		is_in, index := contain(choices, input)
		if is_in {
			result_index = index
			break
		}

		// Check index inputs
		int_input, err := strconv.Atoi(input)
		if err == nil {
			result_index = int_input
			break
		} else {
			fmt.Println("cannot recognize")
		}
	}
	return result_index
}

// Return bool and index
// If it does not contain, return -1 as index
func contain(list []string, subject string) (bool, int) {
	for i, v := range list {
		if v == subject {
			return true, i
		}
	}
	return false, -1
}
